//go:build js

package ai

import (
	"context"
	"errors"
)

// YzmaProvider on the js/wasm target is a stub: the underlying llama.cpp
// bindings require CGO, which the wasm target doesn't support. Methods
// return errors at call time so the rest of the package keeps the same
// API surface and callers don't have to fork their code.
type YzmaProvider struct {
	modelPath string
}

// NewYzmaProvider returns a stub provider. Calls to Complete /
// CompleteStream return ErrYzmaUnsupported.
func NewYzmaProvider(libPath, modelPath string, contextSize uint32, maxTokens int, temperature float32) *YzmaProvider {
	_ = libPath
	_ = contextSize
	_ = maxTokens
	_ = temperature
	return &YzmaProvider{modelPath: modelPath}
}

// ErrYzmaUnsupported is returned by every YzmaProvider method on the
// wasm target.
var ErrYzmaUnsupported = errors.New("yzma: not supported on wasm — use a remote provider (anthropic, ollama) instead")

func (y *YzmaProvider) Name() string { return "yzma (unavailable on wasm)" }

func (y *YzmaProvider) Complete(ctx context.Context, systemPrompt string, messages []Message) (string, error) {
	return "", ErrYzmaUnsupported
}

func (y *YzmaProvider) CompleteStream(ctx context.Context, systemPrompt string, messages []Message) <-chan StreamToken {
	ch := make(chan StreamToken, 1)
	ch <- StreamToken{Error: ErrYzmaUnsupported, Done: true}
	close(ch)
	return ch
}

// YzmaAvailable always reports false on the wasm target — yzma needs CGO.
func YzmaAvailable(libPath, modelPath string) bool {
	_ = libPath
	_ = modelPath
	return false
}
