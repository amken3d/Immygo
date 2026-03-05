package ui

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"gioui.org/layout"

	"github.com/amken3d/immygo/ai"
)

// DebugInfo holds the layout tree for a single widget.
type DebugInfo struct {
	Type       string       `json:"type"`
	MinW, MinH int          `json:"minW,omitempty"`
	MaxW, MaxH int          `json:"maxW,omitempty"`
	ResultW    int          `json:"resultW,omitempty"`
	ResultH    int          `json:"resultH,omitempty"`
	Children   []*DebugInfo `json:"children,omitempty"`
}

// DebugCollector builds a tree of DebugInfo during a frame.
type DebugCollector struct {
	mu    sync.Mutex
	root  *DebugInfo
	stack []*DebugInfo
	frame uint64
}

var (
	debugEnabled   bool
	debugCollector *DebugCollector
	debugOnce      sync.Once
)

// EnableDebug activates layout debugging. Debug info is printed to stderr
// after each frame. Can also be enabled via IMMYGO_DEBUG=1 environment variable.
func EnableDebug() {
	debugOnce.Do(func() {
		debugEnabled = true
		debugCollector = &DebugCollector{}
	})
}

// initDebugFromEnv checks the IMMYGO_DEBUG environment variable.
func initDebugFromEnv() {
	if os.Getenv("IMMYGO_DEBUG") == "1" {
		EnableDebug()
	}
}

// debugEnter marks the start of a widget's layout. No-op when debug is disabled.
func debugEnter(typeName string, gtx layout.Context) {
	if !debugEnabled || debugCollector == nil {
		return
	}
	debugCollector.enter(typeName, gtx)
}

// debugLeave marks the end of a widget's layout with its result dimensions.
func debugLeave(dims layout.Dimensions) {
	if !debugEnabled || debugCollector == nil {
		return
	}
	debugCollector.leave(dims)
}

func (dc *DebugCollector) enter(typeName string, gtx layout.Context) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	info := &DebugInfo{
		Type: typeName,
		MinW: gtx.Constraints.Min.X,
		MinH: gtx.Constraints.Min.Y,
		MaxW: gtx.Constraints.Max.X,
		MaxH: gtx.Constraints.Max.Y,
	}

	if len(dc.stack) > 0 {
		parent := dc.stack[len(dc.stack)-1]
		parent.Children = append(parent.Children, info)
	} else {
		dc.root = info
	}
	dc.stack = append(dc.stack, info)
}

func (dc *DebugCollector) leave(dims layout.Dimensions) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	if len(dc.stack) == 0 {
		return
	}

	current := dc.stack[len(dc.stack)-1]
	current.ResultW = dims.Size.X
	current.ResultH = dims.Size.Y
	dc.stack = dc.stack[:len(dc.stack)-1]
}

// FlushFrame prints the debug tree and optionally gets AI analysis.
// Called after each frame when debug is enabled.
func (dc *DebugCollector) FlushFrame() {
	dc.mu.Lock()
	root := dc.root
	dc.root = nil
	dc.stack = nil
	dc.frame++
	frame := dc.frame
	dc.mu.Unlock()

	if root == nil {
		return
	}

	// Only log every 60th frame to avoid spam.
	if frame%60 != 1 {
		return
	}

	data, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return
	}

	fmt.Fprintf(os.Stderr, "\n[IMMYGO_DEBUG] Frame %d layout tree:\n%s\n", frame, string(data))

	// Optionally get AI explanation.
	if os.Getenv("IMMYGO_DEBUG_AI") == "1" {
		go func() {
			ctx := context.Background()
			prompt := fmt.Sprintf(`Analyze this ImmyGo layout debug tree and explain any potential issues (overflow, zero-size widgets, excessive nesting, wasted space):

%s

Be concise. Only mention actual problems.`, string(data))

			explanation, err := ai.DefaultAssistant().Chat(ctx, prompt)
			if err != nil {
				return
			}
			fmt.Fprintf(os.Stderr, "[IMMYGO_DEBUG] AI Analysis:\n%s\n", explanation)
		}()
	}
}

// debugFlushFrame is called after each frame render.
func debugFlushFrame() {
	if !debugEnabled || debugCollector == nil {
		return
	}
	debugCollector.FlushFrame()
}
