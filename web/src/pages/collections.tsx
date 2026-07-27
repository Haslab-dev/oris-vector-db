import { Card, CardContent } from "../components/ui/card"
import { Badge } from "../components/ui/badge"
import { Link } from "@tanstack/react-router"

const collections = [
  { name: "Products", vectors: "1,245,221", dimension: 768, distance: "Cosine", segments: 12, size: "420 MB" },
  { name: "Docs", vectors: "250,000", dimension: 1536, distance: "Cosine", segments: 4, size: "380 MB" },
  { name: "Users", vectors: "85,432", dimension: 128, distance: "Dot", segments: 2, size: "42 MB" },
  { name: "Images", vectors: "45,100", dimension: 512, distance: "Cosine", segments: 1, size: "92 MB" },
]

export default function CollectionsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Collections</h1>
        <p className="text-sm text-zinc-500 mt-1">{collections.length} collections total</p>
      </div>

      <Card>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-zinc-200 dark:border-zinc-700 text-left">
                  <th className="p-4 font-medium text-zinc-500">Name</th>
                  <th className="p-4 font-medium text-zinc-500 text-right">Vectors</th>
                  <th className="p-4 font-medium text-zinc-500 text-right">Dimension</th>
                  <th className="p-4 font-medium text-zinc-500">Distance</th>
                  <th className="p-4 font-medium text-zinc-500 text-right">Segments</th>
                  <th className="p-4 font-medium text-zinc-500 text-right">Size</th>
                </tr>
              </thead>
              <tbody>
                {collections.map((col) => (
                  <tr key={col.name} className="border-b border-zinc-100 dark:border-zinc-800 hover:bg-zinc-50 dark:hover:bg-zinc-800/50">
                    <td className="p-4">
                      <Link to="/collections/$name" params={{ name: col.name.toLowerCase() }} className="font-medium text-zinc-900 dark:text-zinc-50 hover:text-blue-600 dark:hover:text-blue-400">
                        {col.name}
                      </Link>
                    </td>
                    <td className="p-4 text-right text-zinc-600 dark:text-zinc-400">{col.vectors}</td>
                    <td className="p-4 text-right text-zinc-600 dark:text-zinc-400">{col.dimension}</td>
                    <td className="p-4"><Badge variant="secondary">{col.distance}</Badge></td>
                    <td className="p-4 text-right text-zinc-600 dark:text-zinc-400">{col.segments}</td>
                    <td className="p-4 text-right text-zinc-600 dark:text-zinc-400">{col.size}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
