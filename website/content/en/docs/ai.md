---
title: "AI Integration"
linkTitle: "AI Integration"
description: "Local AI, MCP server, dev tools, and prototyping"
weight: 5
---

ImmyGo includes built-in AI capabilities powered by [Yzma](https://github.com/nicholasgasior/yzma) for local LLM inference. This means your apps can have AI features — chat, autocomplete, summarization — without external API calls or API keys. Everything runs locally on the user's machine.

ImmyGo also provides an **AI-first developer workflow** with an MCP server, conversational dev mode, runtime UI prototyping, and layout debugging — designed for developers who use Claude Code, Cursor, and other AI tools.

## Architecture

```
ai.Engine          Low-level: loads model, runs inference, manages context
    |
ai.Assistant       High-level: Chat(), Suggest(), Summarize() with history
    |
ai.ChatPanel       UI widget: complete chat interface with bubbles + input

ai.WidgetCatalog()      Shared API reference for AI prompts
ai.DefaultAssistant()   Singleton assistant for CLI tools
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

### 2. Load a Real Model

To use actual AI inference, point the engine at a GGUF model file:

```go
engine := ai.NewEngine(ai.Config{
    ModelPath:  "/path/to/model.gguf",  // Any GGUF model (e.g., Llama, Mistral, Phi)
    ContextLen: 2048,                    // Context window size
    MaxTokens:  512,                     // Max tokens per response
    Temperature: 0.7,                    // Creativity (0.0 = deterministic, 1.0 = creative)
    Threads:    4,                       // CPU threads for inference
})

// Load the model (call once at startup)
err := engine.Load()
if err != nil {
    log.Fatal(err)
}
```

### 3. Async Model Loading

For large models, load asynchronously so the UI doesn't freeze:

```go
assistant := ai.NewAssistant("Helper", engine)

// Load in background, get notified when ready
assistant.LoadAsync(func() {
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
// Default configuration (no model loaded)
cfg := ai.DefaultConfig()
// Returns: Config{ContextLen: 2048, MaxTokens: 512, Temperature: 0.7, Threads: 4}

// Custom configuration
cfg := ai.Config{
    ModelPath:  "models/phi-3-mini.gguf",
    LibPath:    "",        // Optional: custom path to Yzma/llama.cpp library
    ContextLen: 4096,
    MaxTokens:  1024,
    Temperature: 0.5,
    Threads:    8,
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

ImmyGo provides a suite of AI-powered tools for developers who build UIs with AI assistants.

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
# Default template
immygo new myapp

# AI-generated from description
immygo new myapp --ai "a calculator with basic arithmetic operations"
immygo new myapp --ai "a todo list with add, delete, and mark complete"
immygo new myapp --ai "a dashboard with sidebar navigation and charts"
```

The AI generates a complete `main.go` using the `ui` package with proper state management, layout, and event handling.

### Conversational Dev Mode (`immygo dev --ai`)

Start the dev server with an interactive AI prompt alongside file watching:

```bash
immygo dev --ai ./myapp/
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

ImmyGo works with any GGUF-format model. Recommendations by use case:

| Use Case | Model | Size | Notes |
|----------|-------|------|-------|
| General chat | Phi-3 Mini | ~2.3GB | Good balance of quality and speed |
| Code assistance | CodeLlama 7B | ~4GB | Strong code generation |
| Fast responses | TinyLlama 1.1B | ~0.6GB | Very fast, lower quality |
| High quality | Mistral 7B | ~4GB | Best quality at 7B scale |

Download GGUF models from [HuggingFace](https://huggingface.co/models?sort=trending&search=gguf).
