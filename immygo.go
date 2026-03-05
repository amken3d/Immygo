// Package immygo provides a high-level Go UI framework built on Gio.
//
// ImmyGo makes building beautiful, modern desktop applications in Go as easy
// as building them with Avalonia in .NET. It wraps Gio's powerful but low-level
// rendering engine into an ergonomic, Fluent Design-inspired widget toolkit
// with built-in AI capabilities via Yzma.
//
// # Quick Start
//
//	package main
//
//	import (
//	    "github.com/amken3d/immygo/app"
//	    "github.com/amken3d/immygo/widget"
//	    "github.com/amken3d/immygo/layout"
//	    "github.com/amken3d/immygo/theme"
//	    giolayout "gioui.org/layout"
//	)
//
//	func main() {
//	    app.New("My App").
//	        WithLayout(func(gtx giolayout.Context, th *theme.Theme) giolayout.Dimensions {
//	            return layout.NewVStack().
//	                Child(func(gtx giolayout.Context) giolayout.Dimensions {
//	                    return widget.H1("Hello, ImmyGo!").Layout(gtx, th)
//	                }).
//	                Child(func(gtx giolayout.Context) giolayout.Dimensions {
//	                    return widget.NewButton("Click Me").Layout(gtx, th)
//	                }).
//	                Layout(gtx)
//	        }).
//	        Run()
//	}
//
// # Architecture
//
// ImmyGo is organized into focused packages:
//
//   - app:    Application scaffold, window management, event loop
//   - theme:  Fluent Design-inspired colors, typography, spacing
//   - widget: Ready-to-use controls (Button, TextField, Card, ListView, etc.)
//   - layout: Avalonia-style layout panels (VStack, HStack, Grid, Dock, Wrap)
//   - style:  CSS-like styling with pseudo-class state management
//   - ai:     AI capabilities via Yzma (local LLM chat, autocomplete, summarization)
//
// # Design Principles
//
//   - Beautiful by default: Fluent Design theme with light/dark modes
//   - Easy to learn: Chainable builder APIs, no raw Gio ops needed
//   - AI-native: Built-in AI assistant and smart widgets via Yzma
//   - Composable: Mix and match layouts and widgets freely
//   - Go-idiomatic: No code generation, no macros, just Go
package immygo

const Version = "0.1.0"
