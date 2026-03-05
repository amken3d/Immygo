// Command immygo provides developer tooling for ImmyGo applications.
//
// Usage:
//
//	immygo dev [--ai] [file|dir]     Live-reload development server
//	immygo new <name> [--ai "desc"]  Scaffold a new ImmyGo project
//	immygo mcp                       Start MCP server (stdio)
//	immygo version                   Print version
package main

import (
	"fmt"
	"os"

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
		target, aiEnabled := parseDevArgs(os.Args[2:])
		if err := dev.RunWithConfig(target, aiEnabled); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "new":
		name, aiDesc := parseNewArgs(os.Args[2:])
		if name == "" {
			fmt.Fprintln(os.Stderr, "usage: immygo new <project-name> [--ai \"description\"]")
			os.Exit(1)
		}
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

// parseDevArgs extracts target and --ai flag from dev arguments.
func parseDevArgs(args []string) (target string, aiEnabled bool) {
	target = "."
	for i := 0; i < len(args); i++ {
		if args[i] == "--ai" {
			aiEnabled = true
		} else {
			target = args[i]
		}
	}
	return
}

// parseNewArgs extracts project name and --ai description from new arguments.
func parseNewArgs(args []string) (name, aiDesc string) {
	for i := 0; i < len(args); i++ {
		if args[i] == "--ai" && i+1 < len(args) {
			i++
			aiDesc = args[i]
		} else if name == "" {
			name = args[i]
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

Examples:
  immygo dev                            # Watch current directory
  immygo dev ./examples/hello           # Watch specific directory
  immygo dev --ai ./examples/hello      # Watch with AI assistant
  immygo new myapp                      # Create new project "myapp"
  immygo new myapp --ai "a calculator"  # AI-generated calculator app
  immygo mcp                            # Start MCP server for Claude Code`)
}
