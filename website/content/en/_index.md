---
title: "ImmyGo"
linkTitle: "ImmyGo"
---

{{< blocks/cover title="" image_anchor="top" height="full" color="white" >}}
<div class="immygo-hero">
  <img src="/images/immygo-logo-sq.svg" alt="ImmyGo Logo" class="immygo-logo" />
  <h1 class="immygo-title">immy<span class="immygo-go">go</span></h1>
  <p class="lead mt-3">Pure Go. Native performance. Modern design. AI-powered.</p>
  <div class="mt-5">
    <a class="btn btn-lg btn-primary me-3 mb-4" href="{{< relref "/docs" >}}">
      Get Started <i class="fas fa-arrow-alt-circle-right ms-2"></i>
    </a>
    <a class="btn btn-lg btn-outline-dark me-3 mb-4" href="https://github.com/amken3d/immygo">
      <i class="fab fa-github me-2"></i> GitHub
    </a>
  </div>
</div>
{{< /blocks/cover >}}

{{% blocks/lead color="primary" %}}

ImmyGo wraps [Gio](https://gioui.org) into a modern declarative widget toolkit.
Build beautiful desktop apps in Go with a **declarative API**, **Fluent Design** tokens,
and **pluggable AI providers** (Yzma, Ollama, Claude, MCP) — no web stack required.

{{% /blocks/lead %}}

{{< blocks/section color="white" >}}

<div class="text-center mb-4">
  <h2>See It in Action</h2>
  <p class="text-muted">Watch ImmyGo build a native desktop UI — no web stack, just Go</p>
</div>

<div class="immygo-demo-video">
  <video autoplay loop muted playsinline class="immygo-screenshot">
    <source src="/images/demo.webm" type="video/webm" />
    Your browser does not support the video tag.
  </video>
</div>

<div class="text-center mt-5 mb-4">
  <h2>Showcase</h2>
  <p class="text-muted">Light and dark themes with 25+ built-in widgets</p>
</div>

<div class="immygo-showcase-pair">
  <img src="/images/showcase_light.png" alt="ImmyGo Showcase - Light Theme" class="immygo-screenshot" />
  <img src="/images/showcase_dark.png" alt="ImmyGo Showcase - Dark Theme" class="immygo-screenshot" />
</div>

{{< /blocks/section >}}

{{< blocks/section color="dark" type="row" >}}

{{% blocks/feature icon="fas fa-code" title="Two API Levels" %}}
Start fast with the declarative `ui` package — or drop down to the lower-level `widget`/`layout` packages for full Gio control.
{{% /blocks/feature %}}

{{% blocks/feature icon="fas fa-palette" title="Fluent Design" url="/docs/theming/" %}}
Built-in light and dark themes with semantic color tokens, typography scale, spacing, corner radii, and elevation — all customizable.
{{% /blocks/feature %}}

{{% blocks/feature icon="fas fa-robot" title="AI-Powered" url="/docs/ai/" %}}
Pluggable AI providers — Yzma for in-process local inference, Ollama, Anthropic Claude, or any MCP server. AI scaffolding, conversational dev mode, and runtime prototyping.
{{% /blocks/feature %}}

{{% blocks/feature icon="fas fa-th-large" title="25+ Widgets" url="/docs/widgets/" %}}
Button, TextField, Toggle, Card, DataGrid, TreeView, Dialog, Drawer, DatePicker, Navigator, and more — all with smooth animations.
{{% /blocks/feature %}}

{{< /blocks/section >}}

{{< blocks/section color="white" >}}

<div class="text-center mb-4">
  <h2>Widget Gallery</h2>
  <p class="text-muted">Forms, lists, data grids, overlays, and more — all from a single Go codebase</p>
</div>

<div class="immygo-gallery">
  <figure class="immygo-gallery-item">
    <img src="/images/showcase_forms.png" alt="Form Inputs" class="immygo-screenshot" />
    <figcaption>Form Inputs — text fields, dropdowns, password fields</figcaption>
  </figure>
  <figure class="immygo-gallery-item">
    <img src="/images/showcase_Lists.png" alt="Lists and Dropdowns" class="immygo-screenshot" />
    <figcaption>Lists & Dropdowns — scrollable, selectable lists</figcaption>
  </figure>
  <figure class="immygo-gallery-item">
    <img src="/images/showcase_data-nav.png" alt="Data Grid and Navigation" class="immygo-screenshot" />
    <figcaption>Data & Navigation — sortable data grids, tree views</figcaption>
  </figure>
  <figure class="immygo-gallery-item">
    <img src="/images/showcase_overlays-state.png" alt="Overlays and State" class="immygo-screenshot" />
    <figcaption>Overlays & State — snackbars, dialogs, computed state</figcaption>
  </figure>
</div>

{{< /blocks/section >}}

{{% blocks/section %}}

## Quick Install

```bash
go get github.com/amken3d/immygo
```

Or scaffold a project with the CLI:

```bash
# Default template
immygo new myapp

# AI-generated from a description
immygo new myapp --ai "a todo list with add and delete"

# With a specific AI provider
immygo new myapp --ai "a dashboard" --provider ollama --model qwen2.5-coder
```

{{% /blocks/section %}}

{{% blocks/section color="dark" %}}

## Hello, ImmyGo

```go
package main

import (
    "fmt"
    "github.com/amken3d/immygo/ui"
)

func main() {
    count := ui.NewState(0)

    ui.Run("My App", func() ui.View {
        return ui.Centered(
            ui.VStack(
                ui.Text("Hello, ImmyGo!").Headline(),
                ui.Button("+1").OnClick(func() {
                    count.Update(func(n int) int { return n + 1 })
                }),
                ui.Text(fmt.Sprintf("Count: %d", count.Get())).Bold(),
            ).Spacing(12),
        )
    })
}
```

{{% /blocks/section %}}
