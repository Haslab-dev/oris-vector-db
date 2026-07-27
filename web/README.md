# Oris Web UI

A lightweight, embedded-first web interface for the Oris vector retrieval engine.

Built with **React 19** + **Vite 8** + **Tailwind CSS v4** + **TanStack Router**.

## Quick Start

```bash
cd web
bun install
bun run dev
```

Open http://localhost:5173 in your browser.

## Pages

| Page | Route | Description |
|------|-------|-------------|
| Dashboard | `/` | Stats overview (collections, vectors, storage, memory) |
| Collections | `/collections` | List all collections with metadata |
| Collection Detail | `/collections/$name` | Per-collection stats, segments, and tabs |
| Search | `/search` | Query playground with result detail panel |
| Performance | `/performance` | Latency, insert, memory, and storage charts |
| Settings | `/settings` | Engine and instance configuration |

## Tech Stack

- **React 19** — UI framework
- **TanStack Router** — type-safe routing with hash history
- **Tailwind CSS v4** — utility-first styling
- **Radix UI** — accessible primitives (Tabs, Slot)
- **Recharts** — charts and visualizations
- **Lucide React** — icons
- **shadcn/ui style** — component patterns (Button, Card, Badge)

## Architecture

```
web/
├── src/
│   ├── components/
│   │   ├── layout.tsx          # App shell with sidebar navigation
│   │   └── ui/                 # Primitive components (button, card, badge)
│   ├── lib/
│   │   └── utils.ts            # cn() utility for className merging
│   ├── pages/
│   │   ├── dashboard.tsx
│   │   ├── collections.tsx
│   │   ├── collection-detail.tsx
│   │   ├── search.tsx
│   │   ├── performance.tsx
│   │   └── settings.tsx
│   ├── router.tsx              # Route tree definition
│   ├── App.tsx                 # Root component
│   ├── main.tsx                # Entry point
│   └── index.css               # Tailwind import
├── index.html
├── vite.config.ts
├── tsconfig.json
└── package.json
```

## Connecting to Oris Backend

The web UI communicates with the Oris Go HTTP API. Start the API server:

```go
col, _ := api.Open("/path/to/collection", cfg)
http.ListenAndServe(":8080", api.NewHTTPServer(col))
```

Update the API base URL in the search page's `handleSearch` function to point to your server.

## Build for Production

```bash
bun run build
```

Output goes to `web/dist/`.
