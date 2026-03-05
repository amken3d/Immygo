package ai

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/hybridgroup/yzma/pkg/llama"
)

// YzmaProvider runs local LLM inference via Yzma (llama.cpp Go bindings).
// It loads a GGUF model file and generates completions entirely on-device.
type YzmaProvider struct {
	modelPath   string
	libPath     string
	contextSize uint32
	maxTokens   int
	temperature float32

	mu    sync.Mutex
	model llama.Model
	vocab llama.Vocab
	ready bool
}

// NewYzmaProvider creates a provider that loads a GGUF model via Yzma.
// libPath is the path to the llama.cpp shared library (or YZMA_LIB env var).
// modelPath is the path to the .gguf model file.
func NewYzmaProvider(libPath, modelPath string, contextSize uint32, maxTokens int, temperature float32) *YzmaProvider {
	if libPath == "" {
		libPath = os.Getenv("YZMA_LIB")
	}
	if contextSize == 0 {
		contextSize = 4096
	}
	if maxTokens == 0 {
		maxTokens = 2048
	}
	if temperature == 0 {
		temperature = 0.7
	}
	return &YzmaProvider{
		modelPath:   modelPath,
		libPath:     libPath,
		contextSize: contextSize,
		maxTokens:   maxTokens,
		temperature: temperature,
	}
}

func (y *YzmaProvider) Name() string {
	base := filepath.Base(y.modelPath)
	return fmt.Sprintf("yzma (%s)", base)
}

// init loads the library and model if not already loaded.
func (y *YzmaProvider) init() error {
	y.mu.Lock()
	defer y.mu.Unlock()

	if y.ready {
		return nil
	}

	if y.libPath == "" {
		return fmt.Errorf("yzma: YZMA_LIB not set — path to llama.cpp shared library required")
	}
	if y.modelPath == "" {
		return fmt.Errorf("yzma: no model path specified")
	}

	llama.Load(y.libPath)
	llama.LogSet(llama.LogSilent())
	llama.Init()

	model, err := llama.ModelLoadFromFile(y.modelPath, llama.ModelDefaultParams())
	if err != nil {
		return fmt.Errorf("yzma: failed to load model %q: %w", y.modelPath, err)
	}
	if model == 0 {
		return fmt.Errorf("yzma: model loaded as null from %q", y.modelPath)
	}

	y.model = model
	y.vocab = llama.ModelGetVocab(model)
	y.ready = true
	return nil
}

// formatPrompt applies the model's chat template to format messages.
func (y *YzmaProvider) formatPrompt(systemPrompt string, messages []Message) string {
	var chatMsgs []llama.ChatMessage

	// Add system prompt
	if systemPrompt != "" {
		chatMsgs = append(chatMsgs, llama.NewChatMessage("system", systemPrompt))
	}

	// Add conversation messages (skip system role, already handled)
	for _, m := range messages {
		if m.Role == RoleSystem {
			continue
		}
		chatMsgs = append(chatMsgs, llama.NewChatMessage(string(m.Role), m.Content))
	}

	// Get model's chat template
	template := llama.ModelChatTemplate(y.model, "")
	if template == "" {
		template = "chatml"
	}

	// Apply template
	buf := make([]byte, 32768)
	length := llama.ChatApplyTemplate(template, chatMsgs, true, buf)
	if length <= 0 {
		// Fallback: concatenate messages manually
		var result string
		for _, m := range chatMsgs {
			result += fmt.Sprintf("<%s>\n%s\n</%s>\n", "msg", m.Content, "msg")
		}
		return result
	}
	return string(buf[:length])
}

// generate runs inference and returns the full response.
func (y *YzmaProvider) generate(prompt string, maxTokens int, stream func(string)) (string, error) {
	// Create a fresh context for each generation
	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = y.contextSize

	ctx, err := llama.InitFromModel(y.model, ctxParams)
	if err != nil {
		return "", fmt.Errorf("yzma: failed to create context: %w", err)
	}
	if ctx == 0 {
		return "", fmt.Errorf("yzma: context is null")
	}
	defer llama.Free(ctx)

	// Tokenize the prompt
	tokens := llama.Tokenize(y.vocab, prompt, true, true)
	if len(tokens) == 0 {
		return "", fmt.Errorf("yzma: tokenization produced no tokens")
	}

	batch := llama.BatchGetOne(tokens)

	// Configure sampler with temperature
	samplerParams := llama.DefaultSamplerParams()
	samplerParams.Temp = y.temperature
	samplerParams.TopK = 40
	samplerParams.TopP = 0.95
	samplerParams.MinP = 0.05
	samplerParams.PenaltyRepeat = 1.1
	samplerParams.PenaltyLastN = 64

	samplerTypes := []llama.SamplerType{
		llama.SamplerTypePenalties,
		llama.SamplerTypeTopK,
		llama.SamplerTypeTopP,
		llama.SamplerTypeMinP,
		llama.SamplerTypeTemperature,
	}
	sampler := llama.NewSampler(y.model, samplerTypes, samplerParams)
	defer llama.SamplerFree(sampler)

	// Generate tokens
	var response string
	buf := make([]byte, 256)

	for i := 0; i < maxTokens; i++ {
		llama.Decode(ctx, batch)
		token := llama.SamplerSample(sampler, ctx, -1)

		if llama.VocabIsEOG(y.vocab, token) {
			break
		}

		length := llama.TokenToPiece(y.vocab, token, buf, 0, true)
		piece := string(buf[:length])
		response += piece

		if stream != nil {
			stream(piece)
		}

		llama.SamplerAccept(sampler, token)
		batch = llama.BatchGetOne([]llama.Token{token})
	}

	return response, nil
}

func (y *YzmaProvider) Complete(_ context.Context, systemPrompt string, messages []Message) (string, error) {
	if err := y.init(); err != nil {
		return "", err
	}

	prompt := y.formatPrompt(systemPrompt, messages)
	return y.generate(prompt, y.maxTokens, nil)
}

func (y *YzmaProvider) CompleteStream(goCtx context.Context, systemPrompt string, messages []Message) <-chan StreamToken {
	ch := make(chan StreamToken, 32)

	go func() {
		defer close(ch)

		if err := y.init(); err != nil {
			ch <- StreamToken{Error: err}
			return
		}

		prompt := y.formatPrompt(systemPrompt, messages)
		_, err := y.generate(prompt, y.maxTokens, func(piece string) {
			select {
			case ch <- StreamToken{Text: piece}:
			case <-goCtx.Done():
				return
			}
		})
		if err != nil {
			ch <- StreamToken{Error: err}
			return
		}
		ch <- StreamToken{Done: true}
	}()

	return ch
}

// YzmaAvailable checks whether Yzma can be used (library and model exist).
func YzmaAvailable(libPath, modelPath string) bool {
	if libPath == "" {
		libPath = os.Getenv("YZMA_LIB")
	}
	if libPath == "" || modelPath == "" {
		return false
	}
	if _, err := os.Stat(libPath); err != nil {
		return false
	}
	if _, err := os.Stat(modelPath); err != nil {
		return false
	}
	return true
}

// YzmaDefaultLibPath returns a sensible default path for the llama.cpp library.
func YzmaDefaultLibPath() string {
	if p := os.Getenv("YZMA_LIB"); p != "" {
		return p
	}
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "linux":
		return filepath.Join(home, ".local", "lib", "libllama.so")
	case "darwin":
		return filepath.Join(home, ".local", "lib", "libllama.dylib")
	case "windows":
		return filepath.Join(home, "AppData", "Local", "yzma", "llama.dll")
	}
	return ""
}
