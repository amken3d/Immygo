// Package dev implements the live-reload development server for ImmyGo.
//
// It watches .go files for changes, rebuilds the binary, and restarts the
// application automatically. This dramatically speeds up the UI development
// cycle — save a file and see the result instantly.
package dev

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Config holds the dev server configuration.
type Config struct {
	// Target is the file or directory to build.
	Target string
	// BuildDir is the directory containing the Go files.
	BuildDir string
	// BinPath is the path to the compiled binary.
	BinPath string
	// PollInterval is how often to check for changes.
	PollInterval time.Duration
	// BuildTags are extra build tags to pass to go build.
	BuildTags string
	// AIEnabled starts the conversational AI mode alongside the dev server.
	AIEnabled bool
}

// Run starts the live-reload dev server. It watches for changes in .go files,
// rebuilds the project, and restarts the application.
func Run(target string) error {
	return RunWithConfig(target, false)
}

// RunWithConfig starts the dev server with optional AI mode.
func RunWithConfig(target string, aiEnabled bool) error {
	cfg, err := resolveConfig(target)
	if err != nil {
		return err
	}
	cfg.AIEnabled = aiEnabled
	defer os.Remove(cfg.BinPath)

	printBanner(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C gracefully
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\n\033[33m⟳ Shutting down...\033[0m")
		cancel()
	}()

	server := &Server{
		cfg:    cfg,
		hashes: make(map[string][32]byte),
	}

	if cfg.AIEnabled {
		aiMode := NewAIMode(cfg.BuildDir)
		go aiMode.Run(ctx)
	}

	return server.run(ctx)
}

// Server manages the build-watch-restart cycle.
type Server struct {
	cfg    *Config
	proc   *exec.Cmd
	mu     sync.Mutex
	hashes map[string][32]byte
}

func (s *Server) run(ctx context.Context) error {
	// Initial build and run
	if err := s.buildAndRun(ctx); err != nil {
		fmt.Printf("\033[31m✗ Build failed:\033[0m %v\n", err)
		fmt.Println("\033[33m⟳ Watching for changes...\033[0m")
	}

	// Snapshot initial hashes
	s.snapshotHashes()

	ticker := time.NewTicker(s.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.stop()
			return nil
		case <-ticker.C:
			if s.hasChanges() {
				fmt.Printf("\033[36m⟳ Change detected, rebuilding...\033[0m\n")
				s.stop()
				if err := s.buildAndRun(ctx); err != nil {
					fmt.Printf("\033[31m✗ Build failed:\033[0m %v\n", err)
					fmt.Println("\033[33m⟳ Watching for changes...\033[0m")
				}
			}
		}
	}
}

func (s *Server) buildAndRun(ctx context.Context) error {
	start := time.Now()

	// Build
	args := []string{"build", "-o", s.cfg.BinPath}
	if s.cfg.BuildTags != "" {
		args = append(args, "-tags", s.cfg.BuildTags)
	}
	args = append(args, s.cfg.BuildDir)

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build: %w", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("\033[32m✓ Built in %dms\033[0m\n", elapsed.Milliseconds())

	// Run the binary
	s.mu.Lock()
	s.proc = exec.CommandContext(ctx, s.cfg.BinPath)
	s.proc.Stdout = os.Stdout
	s.proc.Stderr = os.Stderr
	s.proc.Env = os.Environ()
	// Start in its own process group so we can kill it cleanly
	s.proc.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	s.mu.Unlock()

	if err := s.proc.Start(); err != nil {
		return fmt.Errorf("start app: %w", err)
	}

	fmt.Printf("\033[32m▶ Running (PID %d)\033[0m\n", s.proc.Process.Pid)

	// Wait for exit in background (non-blocking)
	go func() {
		_ = s.proc.Wait()
	}()

	return nil
}

func (s *Server) stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.proc == nil || s.proc.Process == nil {
		return
	}

	// Kill the process group
	pgid, err := syscall.Getpgid(s.proc.Process.Pid)
	if err == nil {
		_ = syscall.Kill(-pgid, syscall.SIGTERM)
	}

	// Give it a moment to exit gracefully
	done := make(chan struct{})
	go func() {
		_ = s.proc.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		if pgid > 0 {
			_ = syscall.Kill(-pgid, syscall.SIGKILL)
		}
	}

	s.proc = nil
}

func (s *Server) snapshotHashes() {
	s.hashes = make(map[string][32]byte)
	_ = filepath.WalkDir(s.cfg.BuildDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if name == "vendor" || name == ".git" || strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if isWatchedFile(path) {
			if h, err := hashFile(path); err == nil {
				s.hashes[path] = h
			}
		}
		return nil
	})
}

func (s *Server) hasChanges() bool {
	changed := false
	newHashes := make(map[string][32]byte)

	_ = filepath.WalkDir(s.cfg.BuildDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if name == "vendor" || name == ".git" || strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if isWatchedFile(path) {
			h, err := hashFile(path)
			if err != nil {
				return nil
			}
			newHashes[path] = h
			if old, ok := s.hashes[path]; !ok || old != h {
				changed = true
			}
		}
		return nil
	})

	// Check for deleted files
	for path := range s.hashes {
		if _, ok := newHashes[path]; !ok {
			changed = true
		}
	}

	if changed {
		s.hashes = newHashes
	}
	return changed
}

func resolveConfig(target string) (*Config, error) {
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	info, err := os.Stat(absTarget)
	if err != nil {
		return nil, fmt.Errorf("stat target: %w", err)
	}

	buildDir := absTarget
	if !info.IsDir() {
		buildDir = filepath.Dir(absTarget)
	}

	// Create temp file for the binary
	tmpDir := os.TempDir()
	binPath := filepath.Join(tmpDir, fmt.Sprintf("immygo-dev-%d", os.Getpid()))

	return &Config{
		Target:       absTarget,
		BuildDir:     buildDir,
		BinPath:      binPath,
		PollInterval: 500 * time.Millisecond,
	}, nil
}

func isWatchedFile(path string) bool {
	ext := filepath.Ext(path)
	switch ext {
	case ".go", ".mod", ".sum":
		return true
	}
	return false
}

func hashFile(path string) ([32]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return [32]byte{}, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return [32]byte{}, err
	}

	var sum [32]byte
	copy(sum[:], h.Sum(nil))
	return sum, nil
}

func printBanner(cfg *Config) {
	fmt.Println("\033[1;36m")
	fmt.Println("  ╔══════════════════════════════════════╗")
	fmt.Println("  ║       ImmyGo Dev Server v0.1.0       ║")
	fmt.Println("  ╚══════════════════════════════════════╝")
	fmt.Println("\033[0m")
	fmt.Printf("  \033[90mWatching:\033[0m  %s\n", cfg.BuildDir)
	fmt.Printf("  \033[90mPolling:\033[0m   every %s\n", cfg.PollInterval)
	fmt.Printf("  \033[90mFiles:\033[0m     *.go, go.mod, go.sum\n")
	fmt.Println()
	fmt.Println("  \033[33mPress Ctrl+C to stop\033[0m")
	fmt.Println()
}
