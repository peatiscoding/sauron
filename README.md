# sauron

**Technical document preview demon.**

A live preview server that watches your editor and renders markdown, MDX, YAML, and JSON in real-time. Edit in Neovim or Claude Code — sauron sees all.

![Go](https://img.shields.io/badge/Go-1.21-00ADD8?logo=go&logoColor=white)
![Svelte](https://img.shields.io/badge/Svelte-5-FF3E00?logo=svelte&logoColor=white)

## Features

- **Real-time preview** — WebSocket-powered live updates as you type
- **Markdown** — GFM, syntax highlighting, relative image paths
- **MDX** — full MDX evaluation with Preact components
- **Mermaid diagrams** — rendered as interactive SVGs with click-to-expand
- **OpenAPI viewer** — interactive API docs from YAML specs via RapiDoc
- **JSON viewer** — syntax-highlighted, color-coded formatting
- **Dark/light themes** — follows system preference, persisted in localStorage
- **PDF export** — print-optimized stylesheets, one click
- **Git status** — shows staged/modified/untracked in the status bar
- **The Eye** — animated connection status indicator that blinks and watches

## Architecture

```
┌─────────────┐     HTTP/WS      ┌──────────────┐     WebSocket     ┌─────────────┐
│   Neovim    │ ──────────────▶  │  Go Server   │ ──────────────▶   │   Browser   │
│  (plugin)   │   :6905/focus    │  (sauron)    │   :6905/ws        │  (Svelte)   │
└─────────────┘                  └──────────────┘                   └─────────────┘
```

**Backend** — Go HTTP server with embedded static assets, WebSocket hub for broadcasting updates.

**Frontend** — Svelte 5 SPA. Renders markdown (marked + highlight.js), MDX (@mdx-js/mdx + Preact), YAML (RapiDoc), JSON, and Mermaid diagrams.

**Editor plugin** — Neovim Lua plugin that detects buffer changes and sends content to the server. Claude Code hook also supported.

## Quick Start

### Build

```bash
./scripts/build.sh
```

Installs npm deps, builds frontend into `server/static/`, compiles Go binary to `./sauron`.

### Run

```bash
./sauron
```

Opens on `http://localhost:6905`.

### Neovim Setup

Add to your Neovim config:

```lua
require('sauron').setup({
  url = "http://localhost:6905",
  debounce_ms = 300,
})
```

Auto-detects filetypes: `markdown`, `mdx`, `yaml`, `json`.

Commands:
- `:SauronFocus` — force push current buffer
- `:SauronAnchor` — pin cursor position as scroll anchor

### Claude Code Integration

Add to `~/.claude/settings.json`:

```json
{
  "hooks": {
    "PostToolUse": [{
      "matcher": "Edit|Write|MultiEdit",
      "hooks": [{
        "type": "command",
        "command": "/path/to/sauron/scripts/sauron-hook.sh"
      }]
    }]
  }
}
```

Files edited by Claude stream to sauron automatically.

## Project Structure

```
├── server/              # Go backend
│   ├── main.go          # HTTP server, routing
│   ├── handlers.go      # /focus, /anchor, /image endpoints
│   ├── hub.go           # WebSocket hub
│   └── static/          # Built frontend (generated)
├── web/                 # Svelte frontend
│   └── src/
│       ├── App.svelte
│       └── lib/
│           ├── MarkdownView.svelte
│           ├── OpenAPIView.svelte
│           ├── JSONView.svelte
│           └── SauronEye.svelte
├── nvim/                # Neovim plugin
│   └── lua/sauron/
│       └── init.lua
├── scripts/
│   ├── build.sh
│   └── sauron-hook.sh   # Claude Code hook
└── sauron               # Compiled binary
```

## WebSocket Protocol

Server broadcasts JSON to connected clients:

```json
{
  "type": "focus",
  "filename": "README.md",
  "filetype": "markdown",
  "filepath": "/absolute/path/README.md",
  "content": "# Hello",
  "gitStatus": "modified"
}
```

```json
{
  "type": "anchor",
  "line": 42,
  "total": 120
}
```

## Tech Stack

| Layer    | Tech                                                                                                                                                                                                                     |
|----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Backend  | [Go](https://go.dev/) 1.21, [gorilla/websocket](https://github.com/gorilla/websocket) 1.5                                                                                                                               |
| Frontend | [Svelte](https://svelte.dev/) 5, [Vite](https://vite.dev/) 5, [marked](https://github.com/markedjs/marked) 13, [highlight.js](https://highlightjs.org/) 11, [Mermaid](https://mermaid.js.org/) 11                       |
| MDX      | [@mdx-js/mdx](https://mdxjs.com/) 3.1, [Preact](https://preactjs.com/) 10, [remark-gfm](https://github.com/remarkjs/remark-gfm) 4, [remark-frontmatter](https://github.com/remarkjs/remark-frontmatter) 5             |
| API Docs | [RapiDoc](https://rapidocweb.com/) 9.3, [js-yaml](https://github.com/nodeca/js-yaml) 4                                                                                                                                 |
| Editor   | [Neovim](https://neovim.io/) (Lua), [Claude Code](https://docs.anthropic.com/en/docs/claude-code) (bash hook)                                                                                                           |
