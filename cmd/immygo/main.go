// Command immygo provides developer tooling for ImmyGo applications.
//
// Usage:
//
//	immygo dev [--ai] [--provider X] [--model Y] [file|dir]
//	immygo new <name> [--ai "desc"] [--provider X] [--model Y]
//	immygo mcp                       Start MCP server (stdio)
//	immygo version                   Print version
package main

import (
	"fmt"
	"os"

	"github.com/amken3d/immygo/ai"
	"github.com/amken3d/immygo/cmd/immygo/dev"
	"github.com/amken3d/immygo/cmd/immygo/mcp"
	"github.com/amken3d/immygo/cmd/immygo/scaffold"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "dev":
		target, aiEnabled, pcfg := parseDevArgs(os.Args[2:])
		ai.SetDefaultProviderConfig(pcfg)
		if err := dev.RunWithConfig(target, aiEnabled); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "new":
		name, aiDesc, pcfg := parseNewArgs(os.Args[2:])
		if name == "" {
			fmt.Fprintln(os.Stderr, "usage: immygo new <project-name> [--ai \"description\"] [--provider X] [--model Y]")
			os.Exit(1)
		}
		ai.SetDefaultProviderConfig(pcfg)
		if err := scaffold.Run(name, aiDesc); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "mcp":
		if err := mcp.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "version":
		fmt.Printf("immygo %s\n", version)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

// parseDevArgs extracts target, --ai flag, and provider options from dev arguments.
func parseDevArgs(args []string) (target string, aiEnabled bool, pcfg ai.ProviderConfig) {
	target = "."
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--ai":
			aiEnabled = true
		case "--provider":
			if i+1 < len(args) {
				i++
				pcfg.Provider = args[i]
			}
		case "--model":
			if i+1 < len(args) {
				i++
				pcfg.Model = args[i]
			}
		case "--mcp-command":
			if i+1 < len(args) {
				i++
				pcfg.MCPCommand = args[i]
			}
		case "--mcp-tool":
			if i+1 < len(args) {
				i++
				pcfg.MCPTool = args[i]
			}
		default:
			target = args[i]
		}
	}
	return
}

// parseNewArgs extracts project name, --ai description, and provider options from new arguments.
func parseNewArgs(args []string) (name, aiDesc string, pcfg ai.ProviderConfig) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--ai":
			if i+1 < len(args) {
				i++
				aiDesc = args[i]
			}
		case "--provider":
			if i+1 < len(args) {
				i++
				pcfg.Provider = args[i]
			}
		case "--model":
			if i+1 < len(args) {
				i++
				pcfg.Model = args[i]
			}
		case "--mcp-command":
			if i+1 < len(args) {
				i++
				pcfg.MCPCommand = args[i]
			}
		case "--mcp-tool":
			if i+1 < len(args) {
				i++
				pcfg.MCPTool = args[i]
			}
		default:
			if name == "" {
				name = args[i]
			}
		}
	}
	return
}

func printUsage() {
	fmt.Println(`ImmyGo - Beautiful Go UIs made easy

Usage:
  immygo <command> [arguments]

Commands:
  dev [path]                Start live-reload dev server (watches for changes, rebuilds & restarts)
  dev --ai [path]           Start dev server with AI assistant (type at ai> prompt)
  new <name>                Scaffold a new ImmyGo project with starter template
  new <name> --ai "desc"    Scaffold with AI-generated code from description
  mcp                       Start MCP server for AI tool integration (stdio JSON-RPC)
  version                   Print version info
  help                      Show this help

AI Provider Options (for dev --ai and new --ai):
  --provider <name>         Provider: ollama, anthropic, mcp, simulation (default: auto-detect)
  --model <name>            Model name (default: qwen2.5-coder for Ollama)
  --mcp-command <cmd>       MCP server command (e.g. "npx @some/mcp-server")
  --mcp-tool <name>         MCP tool name to call (default: immygo_generate_code)

Environment Variables:
  IMMYGO_PROVIDER           Same as --provider
  IMMYGO_MODEL              Same as --model
  IMMYGO_OLLAMA_HOST        Ollama API URL (default: http://localhost:11434)
  IMMYGO_MCP_COMMAND        Same as --mcp-command
  IMMYGO_MCP_TOOL           Same as --mcp-tool
  ANTHROPIC_API_KEY         Anthropic API key (enables Claude provider)

Examples:
  immygo dev                                         # Watch current directory
  immygo dev ./examples/hello                        # Watch specific directory
  immygo dev --ai ./examples/hello                   # Watch with AI (auto-detect provider)
  immygo dev --ai --provider ollama --model codellama # Use specific Ollama model
  immygo new myapp                                   # Create new project
  immygo new myapp --ai "a calculator"               # AI-generated calculator app
  immygo new myapp --ai "a todo app" --provider mcp --mcp-command "my-mcp-server"
  immygo mcp                                         # Start MCP server for Claude Code`)
}
