---
title: "ImmyGo"
linkTitle: "ImmyGo"
---

{{< blocks/cover title="ImmyGo" image_anchor="top" height="full" >}}
<p class="lead mt-5">A high-level Go UI framework with Fluent Design aesthetics</p>
<a class="btn btn-lg btn-primary me-3 mb-4" href="{{< relref "/docs" >}}">
Documentation <i class="fas fa-arrow-alt-circle-right ms-2"></i>
</a>
<a class="btn btn-lg btn-secondary me-3 mb-4" href="https://github.com/amken3d/immygo">
GitHub <i class="fab fa-github ms-2 "></i>
</a>
{{< /blocks/cover >}}

{{% blocks/lead color="primary" %}}

ImmyGo wraps [Gio](https://gioui.org) into an Avalonia-inspired widget toolkit.
Build beautiful desktop apps in Go with a **declarative API**, **Fluent Design** tokens,
and **built-in AI** capabilities — no web stack required.

{{% /blocks/lead %}}

{{< blocks/section color="dark" type="row" >}}

{{% blocks/feature icon="fas fa-code" title="Two API Levels" %}}
Start fast with the declarative `ui` package — or drop down to the lower-level `widget`/`layout` packages for full Gio control.
{{% /blocks/feature %}}

{{% blocks/feature icon="fas fa-palette" title="Fluent Design" url="/docs/theming/" %}}
Built-in light and dark themes with semantic color tokens, typography scale, spacing, corner radii, and elevation — all customizable.
{{% /blocks/feature %}}

{{% blocks/feature icon="fas fa-robot" title="Built-in AI" url="/docs/ai/" %}}
Local LLM inference via Yzma, runtime UI prototyping, MCP server for Claude Code and Cursor, and conversational dev mode.
{{% /blocks/feature %}}

{{% blocks/feature icon="fas fa-th-large" title="25+ Widgets" url="/docs/widgets/" %}}
Button, TextField, Toggle, Card, DataGrid, TreeView, Dialog, Drawer, DatePicker, Navigator, and more — all with smooth animations.
{{% /blocks/feature %}}

{{< /blocks/section >}}

{{< blocks/section >}}

## Quick Install

```bash
go get github.com/amken3d/immygo
```

Or scaffold a project with the CLI:

```bash
go run github.com/amken3d/immygo/cmd/immygo new myapp
```

{{< /blocks/section >}}
