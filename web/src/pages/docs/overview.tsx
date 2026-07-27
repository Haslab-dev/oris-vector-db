import { Link } from "@tanstack/react-router"
import { Card, CardHeader, CardTitle } from "../../components/ui/card"
import { Rocket, Layout, Cpu, Database, Search, Filter, Layers, GitBranch, Code } from "lucide-react"

const sections = [
  { icon: Rocket, title: "Getting Started", desc: "Quick start guide, installation, and your first search", path: "/docs/getting-started" },
  { icon: Layout, title: "Architecture", desc: "High-level system architecture and design philosophy", path: "/docs/architecture" },
  { icon: Cpu, title: "Dense Engine (HNSW)", desc: "HNSW algorithm, graph construction, and search", path: "/docs/dense-engine" },
  { icon: Search, title: "Sparse Engine (BM25)", desc: "BM25 scoring, inverted index, and tokenization", path: "/docs/sparse-engine" },
  { icon: Filter, title: "Metadata Engine", desc: "Filter AST, Roaring Bitmaps, and filter operators", path: "/docs/metadata-engine" },
  { icon: Database, title: "Storage Engine", desc: "Pebble storage, WAL, snapshots, interfaces", path: "/docs/storage-engine" },
  { icon: Layers, title: "Segment Manager", desc: "Segment lifecycle, mutable/immutable, compaction", path: "/docs/segment-manager" },
  { icon: GitBranch, title: "Query Planner", desc: "Query execution, merging, ranking, hybrid search", path: "/docs/query-planner" },
  { icon: Code, title: "API Reference", desc: "Complete API documentation with examples", path: "/docs/api-reference" },
]

export default function DocsOverviewPage() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Documentation</h1>
        <p className="text-sm text-zinc-500 mt-1">Comprehensive guide to the Oris vector retrieval engine</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {sections.map((section) => (
          <Link key={section.path} to={section.path} className="block">
            <Card className="h-full hover:shadow-md transition-shadow cursor-pointer">
              <CardHeader className="flex flex-row items-start gap-4">
                <section.icon className="h-5 w-5 text-zinc-500 mt-0.5 shrink-0" />
                <div>
                  <CardTitle className="text-sm">{section.title}</CardTitle>
                  <p className="text-xs text-zinc-500 mt-1">{section.desc}</p>
                </div>
              </CardHeader>
            </Card>
          </Link>
        ))}
      </div>
    </div>
  )
}
