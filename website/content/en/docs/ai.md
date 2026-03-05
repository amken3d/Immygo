---
title: "AI Integration"
linkTitle: "AI Integration"
description: "AI providers, local models, MCP server, dev tools, and prototyping"
weight: 5
---

ImmyGo includes built-in AI capabilities with a pluggable provider system. Out of the box it supports **Ollama** (local models), **Anthropic Claude** (cloud API), and **MCP** (any MCP-compatible server). Your apps can have AI features — chat, autocomplete, code generation — and the CLI tools (`immygo new --ai`, `immygo dev --ai`) use the same provider system.

ImmyGo also provides an **AI-first developer workflow** with an MCP server, conversational dev mode, runtime UI prototyping, and layout debugging — designed for developers who use Claude Code, Cursor, and other AI tools.

## Provider Setup

ImmyGo auto-detects available providers in this order: **Ollama** (local) > **Anthropic** (if API key set) > **MCP** (if command configured) > **simulation** (fallback). You can also choose explicitly.

### Option 1: Ollama (Local Models)

The recommended option for local, private AI with no API keys or costs.

**Install Ollama:**
```bash
# Linux
curl -fsSL https://ollama.com/install.sh | sh

# macOS
brew install ollama

# Or download from https://ollama.com
```

**Pull a model:**
```bash
ollama pull qwen2.5-coder    # Default model, good for code generation
ollama pull codellama         # Alternative code model
ollama pull llama3.2          # General purpose
```

**Use it:**
```bash
# Auto-detected (Ollama is tried first)
immygo new myapp --ai "a todo list with add and delete"

# Explicit provider and model
immygo new myapp --ai "a calculator" --provider ollama --model codellama
```

ImmyGo connects to Ollama at `http://localhost:11434` by default. Override with:
```bash
export IMMYGO_OLLAMA_HOST=http://myserver:11434
```

### Option 2: Anthropic Claude (Cloud API)

Use Claude for high-quality code generation. Requires an API key from [console.anthropic.com](https://console.anthropic.com).

**Set your API key:**
```bash
export ANTHROPIC_API_KEY=sk-ant-...
```

**Use it:**
```bash
# Auto-detected when ANTHROPIC_API_KEY is set and Ollama isn't running
immygo new myapp --ai "a dashboard with charts"

# Explicit
immygo new myapp --ai "a chat app" --provider anthropic

# Choose a specific model
immygo new myapp --ai "a settings page" --provider anthropic --model claude-sonnet-4-20250514
```

The default model is `claude-sonnet-4-20250514`. Any Anthropic model ID works.

### Option 3: MCP Server (External Tools)

Connect to any MCP-compatible server for code generation. This lets you use custom or third-party AI tools.

```bash
# Set the MCP server command
export IMMYGO_MCP_COMMAND="npx @some/mcp-server"

# Optionally specify which tool to call (default: immygo_generate_code)
export IMMYGO_MCP_TOOL="my_code_gen_tool"

# Use it
immygo new myapp --ai "a form builder"
```

Or with CLI flags:
```bash
immygo new myapp --ai "a form" --provider mcp --mcp-command "my-mcp-server" --mcp-tool "generate"
```

The MCP client spawns the server as a subprocess and communicates via stdio JSON-RPC 2.0.

### Configuration Reference

Provider selection can be configured via CLI flags or environment variables. CLI flags take priority.

| CLI Flag | Environment Variable | Description |
|----------|---------------------|-------------|
| `--provider <name>` | `IMMYGO_PROVIDER` | Provider: `ollama`, `anthropic`, `mcp`, `simulation` |
| `--model <name>` | `IMMYGO_MODEL` | Model name (provider-specific) |
| | `IMMYGO_OLLAMA_HOST` | Ollama API URL (default: `http://localhost:11434`) |
| `--mcp-command <cmd>` | `IMMYGO_MCP_COMMAND` | MCP server command to spawn |
| `--mcp-tool <name>` | `IMMYGO_MCP_TOOL` | MCP tool name (default: `immygo_generate_code`) |
| | `ANTHROPIC_API_KEY` | Anthropic API key |

### Programmatic Provider Configuration

When embedding AI in your own app, you can configure the provider explicitly:

```go
import "github.com/amken3d/immygo/ai"

// Use Ollama
provider := ai.NewOllamaProvider("http://localhost:11434", "qwen2.5-coder")

// Use Anthropic
provider := ai.NewAnthropicProvider(os.Getenv("ANTHROPIC_API_KEY"), "claude-sonnet-4-20250514")

// Use MCP
provider := ai.NewMCPClientProvider("my-mcp-server", "generate_code")

// Create an engine with a specific provider
cfg := ai.DefaultConfig()
cfg.SystemPrompt = "You are a helpful coding assistant."
engine := ai.NewEngineWithProvider(cfg, provider)
assistant := ai.NewAssistant("My Assistant", engine)

// Or let ImmyGo auto-detect
ai.SetDefaultProviderConfig(ai.ProviderConfig{
    Provider: "ollama",
    Model:    "codellama",
})
assistant := ai.DefaultAssistant()
```

## Architecture

```
ai.Provider            Interface: Complete(), CompleteStream(), Name()
    |
    +-- OllamaProvider      Local models via Ollama HTTP API
    +-- AnthropicProvider    Claude API (cloud)
    +-- MCPClientProvider    Any MCP server (subprocess + stdio JSON-RPC)
    +-- simulationProvider   Fallback for development/testing

ai.Engine          Manages conversation history, routes to Provider
    |
ai.Assistant       High-level: Chat(), Suggest(), Summarize() with history
    |
ai.ChatPanel       UI widget: complete chat interface with bubbles + input

ai.WidgetCatalog()      Shared API reference for AI prompts
ai.DefaultAssistant()   Singleton assistant (auto-detects provider)
ai.GenerateCode()       One-call code generation

MCP Server (immygo mcp)      External AI tool integration
immygo dev --ai               Conversational dev mode
immygo new --ai "desc"        AI-guided scaffolding
ui.Prototype("desc")          Runtime UI generation
ui.EnableDebug() / IMMYGO_DEBUG=1   Layout debugging
```

## Quick Start

### 1. Add a Chat Panel

The fastest way to add AI to your app — a complete chat UI in 4 lines:

```go
import "github.com/amken3d/immygo/ai"

// Create engine and assistant
engine := ai.NewEngine(ai.DefaultConfig())
assistant := ai.NewAssistant("My Assistant", engine)
chatPanel := ai.NewChatPanel(assistant)

// In your layout:
func myLayout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
    return chatPanel.Layout(gtx, th)
}
```

The `ChatPanel` provides:
- Scrollable message history with styled bubbles
- Text input with submit-on-Enter
- Send button
- Async response handling (UI stays responsive)

### 2. Use a Provider

To use real AI inference, configure a provider:

```go
// Ollama (local)
provider := ai.NewOllamaProvider("", "qwen2.5-coder") // "" = default localhost

engine := ai.NewEngineWithProvider(ai.Config{
    ContextSize:  4096,
    MaxTokens:    1024,
    Temperature:  0.7,
    SystemPrompt: "You are a helpful assistant.",
}, provider)
```

Or let auto-detection handle it (checks Ollama, then Anthropic key, then MCP):

```go
engine := ai.NewEngine(ai.DefaultConfig())
engine.Load() // resolves provider automatically
```

### 3. Async Model Loading

For large models, load asynchronously so the UI doesn't freeze:

```go
assistant := ai.NewAssistant("Helper", engine)

// Load in background, get notified when ready
assistant.LoadAsync(func(err error) {
    fmt.Println("Model loaded and ready!")
})

// Check status:
if assistant.IsLoading() { /* show spinner */ }
if assistant.IsReady()   { /* enable chat */ }
```

## Engine API

The engine is the low-level interface for direct model interaction.

### Configuration

```go
// Default configuration
cfg := ai.DefaultConfig()
// Returns: Config{ContextSize: 2048, MaxTokens: 256, Temperature: 0.7}

// Custom configuration with provider
cfg := ai.Config{
    ContextSize:  4096,
    MaxTokens:    1024,
    Temperature:  0.5,
    SystemPrompt: "You are a coding assistant.",
    ProviderConfig: ai.ProviderConfig{
        Provider: "ollama",
        Model:    "codellama",
    },
}
```

### Completion

```go
engine := ai.NewEngine(cfg)
engine.Load()

// Single completion
response, err := engine.Complete(ctx, "What is the capital of France?")
// response = "The capital of France is Paris."

// Streaming completion (token by token)
stream := engine.CompleteStream(ctx, "Explain quantum computing")
for token := range stream {
    if token.Error != nil {
        log.Fatal(token.Error)
    }
    if token.Done {
        break
    }
    fmt.Print(token.Text) // Print each token as it arrives
}

// Conversation history
history := engine.History()       // []Message
engine.ClearHistory()             // Reset conversation

// Check which provider is active
fmt.Println(engine.ProviderName()) // e.g. "ollama (qwen2.5-coder)"
```

## Assistant API

The assistant provides high-level AI operations with conversation management.

```go
assistant := ai.NewAssistant("CodeHelper", engine)
```

### Chat

Send a message and get a response, maintaining conversation history:

```go
response, err := assistant.Chat(ctx, "How do I read a file in Go?")
// response = "To read a file in Go, you can use os.ReadFile()..."
```

### Chat with Streaming

Get responses token-by-token for real-time display:

```go
stream := assistant.ChatStream(ctx, "Explain goroutines")
for token := range stream {
    if token.Done {
        break
    }
    fmt.Print(token.Text)
}
```

### Suggest (Code Completion)

Get completions for partial input — useful for editor integration:

```go
suggestion, err := assistant.Suggest(ctx, "func calculateTax(amount float64) ")
// suggestion = "float64 {\n\treturn amount * 0.2\n}"
```

### Summarize

Condense long text:

```go
summary, err := assistant.Summarize(ctx, longArticleText)
// summary = "The article discusses three main approaches to..."
```

### History Management

```go
assistant.ClearHistory()  // Reset conversation
```

## ChatPanel Widget

The `ChatPanel` is a ready-to-use chat UI widget that combines message display with text input.

```go
chatPanel := ai.NewChatPanel(assistant)

// Layout it like any other widget:
chatPanel.Layout(gtx, th)
```

### Features

- **Message bubbles**: User messages in accent color (right-aligned feel), assistant messages in surface color
- **Scrollable history**: Built-in scroll for long conversations
- **Text input**: Single-line editor with Enter-to-submit
- **Send button**: Circular button with arrow icon
- **Async responses**: Chat stays responsive while the model generates

### Accessing Messages

```go
// Read the message history
for _, msg := range chatPanel.Messages {
    fmt.Printf("[%s] %s\n", msg.Role, msg.Content)
}

// Programmatically add messages
chatPanel.Messages = append(chatPanel.Messages, ai.Message{
    Role:    ai.RoleAssistant,
    Content: "Welcome! How can I help you?",
})
```

## Message Types

```go
type Message struct {
    Role    Role    // RoleSystem, RoleUser, or RoleAssistant
    Content string
}

type Role string

const (
    RoleSystem    Role = "system"
    RoleUser      Role = "user"
    RoleAssistant Role = "assistant"
)
```

## Integration Patterns

### AI-Powered Search

```go
var searchField = widget.NewSearchField()

func onSearch() {
    query := searchField.Text()
    response, _ := assistant.Chat(ctx, "Search for: "+query)
    // Display response...
}
```

### Smart Form Validation

```go
func validateForm() {
    email := emailField.Text()
    prompt := fmt.Sprintf("Is '%s' a valid email address? Reply YES or NO only.", email)
    response, _ := assistant.Chat(ctx, prompt)
    // Use response for validation feedback
}
```

### Contextual Help

```go
func showHelp(widgetName string) {
    prompt := fmt.Sprintf("Briefly explain the '%s' widget in ImmyGo.", widgetName)
    response, _ := assistant.Chat(ctx, prompt)
    // Display in a tooltip or help panel
}
```

## AI Developer Workflow

ImmyGo provides a suite of AI-powered tools for developers who build UIs with AI assistants. All CLI AI features use the same provider system described in [Provider Setup](#provider-setup).

### MCP Server

The MCP (Model Context Protocol) server lets external AI tools like Claude Code and Cursor interact with ImmyGo's API reference and code generation capabilities.

**Start the server:**

```bash
immygo mcp
```

The server communicates via stdio using newline-delimited JSON-RPC 2.0. It exposes three tools:

| Tool | Input | Description |
|------|-------|-------------|
| `immygo_widget_catalog` | `widget?: string` | Returns the full `ui` package API reference, optionally filtered by widget name |
| `immygo_generate_code` | `description: string` | Generates compilable ImmyGo code from a natural language description |
| `immygo_search_docs` | `query: string` | Searches `docs/` markdown files for relevant sections |

**Claude Code configuration** (add to `.mcp.json`):

```json
{
  "mcpServers": {
    "immygo": {
      "command": "go",
      "args": ["run", "./cmd/immygo", "mcp"]
    }
  }
}
```

**Test it:**

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"capabilities":{}}}' | immygo mcp
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | immygo mcp
```

### AI-Guided Scaffolding (`immygo new --ai`)

Generate a new project with AI-written code instead of the default template:

```bash
# Default template (no AI needed)
immygo new myapp

# AI-generated from description (auto-detects provider)
immygo new myapp --ai "a calculator with basic arithmetic operations"
immygo new myapp --ai "a todo list with add, delete, and mark complete"
immygo new myapp --ai "a dashboard with sidebar navigation and charts"

# With a specific provider
immygo new myapp --ai "a file browser" --provider ollama --model codellama
immygo new myapp --ai "a chat app" --provider anthropic
```

The AI generates a complete `main.go` using the `ui` package with proper state management, layout, and event handling.

### Conversational Dev Mode (`immygo dev --ai`)

Start the dev server with an interactive AI prompt alongside file watching:

```bash
immygo dev --ai ./myapp/

# With a specific provider
immygo dev --ai --provider ollama --model qwen2.5-coder ./myapp/
```

This starts the normal live-reload dev server plus an `ai>` prompt where you can describe UI changes in natural language:

```
ai> add a dark mode toggle in the top right corner
  Generating...
  Code updated — rebuilding...

ai> replace the list with a data grid that has Name, Email, Status columns
  Generating...
  Code updated — rebuilding...
```

The AI reads your current `main.go` for context, generates updated code, and writes it. The file watcher detects the change and auto-rebuilds.

### Runtime UI Prototyping (`ui.Prototype()`)

Generate UI at runtime from a description — useful for rapid exploration:

```go
ui.Run("Prototype", func() ui.View {
    return ui.Prototype("a login form with email, password, and remember me checkbox")
})
```

- Shows "Generating..." while the AI works
- Renders the generated view when ready
- Call `.Eject()` to print the generated Go source code for permanent use:

```go
proto := ui.Prototype("a settings panel with toggles")

// Later, when you want to keep the generated code:
proto.Eject() // prints Go source to stdout
```

The AI generates a JSON widget tree which is mapped to real `ui` views. Supported widget types: VStack, HStack, Text, Button, Input, Card, Spacer, Divider, Toggle, Checkbox, Progress.

### Layout Debugger (`IMMYGO_DEBUG`)

Inspect layout constraints and result sizes at runtime to diagnose layout issues:

```bash
# Enable via environment variable
IMMYGO_DEBUG=1 go run ./myapp/

# Also enable AI analysis of the layout tree
IMMYGO_DEBUG=1 IMMYGO_DEBUG_AI=1 go run ./myapp/
```

Or enable programmatically:

```go
ui.EnableDebug()
```

When enabled, the debugger:
1. Records every widget's layout constraints (min/max) and result size
2. Builds a tree of `DebugInfo` nodes matching the view hierarchy
3. Prints the tree as JSON to stderr every 60th frame (to avoid spam)
4. Optionally sends the tree to an AI assistant for analysis (`IMMYGO_DEBUG_AI=1`)

**Example output:**

```json
{
  "type": "VStack",
  "maxW": 1024, "maxH": 768,
  "resultW": 1024, "resultH": 200,
  "children": [
    {"type": "Text", "maxW": 1024, "resultW": 120, "resultH": 24},
    {"type": "Button", "maxW": 1024, "resultW": 100, "resultH": 40}
  ]
}
```

Instrumented widgets: VStack, HStack, Text, Button, Input, Toggle, Card, Checkbox, Dropdown, Progress.

## Code Generation Helpers (`ai/prompts.go`)

The `ai` package includes shared prompt templates used by all AI features:

```go
// Get the full ui package API reference
catalog := ai.WidgetCatalog()

// Get the system prompt for code generation (includes catalog + rules)
system := ai.ImmyGoSystemPrompt()

// Wrap a description for code generation
prompt := ai.CodeGenPrompt("a todo list app")

// Wrap for full app scaffolding
prompt := ai.ScaffoldPrompt("myapp", "a calculator")

// Get a pre-configured singleton assistant
assistant := ai.DefaultAssistant()

// One-call convenience function
code, err := ai.GenerateCode(ctx, "a counter app with reset button")
```

These are used internally by the MCP server, `immygo new --ai`, `immygo dev --ai`, and `ui.Prototype()`.

## Model Recommendations

### Ollama Models

| Use Case | Model | Pull Command | Notes |
|----------|-------|-------------|-------|
| Code generation (default) | Qwen 2.5 Coder | `ollama pull qwen2.5-coder` | Best code quality at 7B, ImmyGo default |
| Code generation (alt) | CodeLlama 7B | `ollama pull codellama` | Good alternative for code |
| General purpose | Llama 3.2 | `ollama pull llama3.2` | Good for chat and general tasks |
| Fast responses | Phi-3 Mini | `ollama pull phi3:mini` | Smaller, faster, still capable |
| High quality | Mistral 7B | `ollama pull mistral` | Strong general quality |

### Anthropic Models

| Model | ID | Notes |
|-------|----|-------|
| Claude Sonnet 4 (default) | `claude-sonnet-4-20250514` | Best balance of speed and quality |
| Claude Haiku 4.5 | `claude-haiku-4-5-20251001` | Fastest, good for simple generation |
| Claude Opus 4.6 | `claude-opus-4-6` | Highest quality, slower |

Override the model with `--model <id>` or `IMMYGO_MODEL=<id>`.
