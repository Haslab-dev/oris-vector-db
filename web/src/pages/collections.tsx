import { useEffect, useState } from "react"
import { Card, CardContent } from "../components/ui/card"
import { Badge } from "../components/ui/badge"
import { Link } from "@tanstack/react-router"
import { fetchCollections, type CollectionInfo } from "../lib/api"

export default function CollectionsPage() {
  const [collections, setCollections] = useState<CollectionInfo[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchCollections()
      .then(setCollections)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <div className="text-sm text-zinc-500">Loading...</div>
  if (error) return <div className="text-sm text-red-500">Error: {error}</div>

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Collections</h1>
        <p className="text-sm text-zinc-500 mt-1">{collections.length} collection{collections.length !== 1 ? "s" : ""} total</p>
      </div>

      {collections.length === 0 ? (
        <Card>
          <CardContent className="p-8 text-center">
            <p className="text-sm text-zinc-500">No collections found. Start by seeding data or creating one via the API.</p>
          </CardContent>
        </Card>
      ) : (
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
                  </tr>
                </thead>
                <tbody>
                  {collections.map((col) => (
                    <tr key={col.name} className="border-b border-zinc-100 dark:border-zinc-800 hover:bg-zinc-50 dark:hover:bg-zinc-800/50">
                      <td className="p-4">
                        <Link to="/collections/$name" params={{ name: col.name }} className="font-medium text-zinc-900 dark:text-zinc-50 hover:text-blue-600 dark:hover:text-blue-400">
                          {col.name}
                        </Link>
                      </td>
                      <td className="p-4 text-right text-zinc-600 dark:text-zinc-400">{col.vectors.toLocaleString()}</td>
                      <td className="p-4 text-right text-zinc-600 dark:text-zinc-400">{col.dimension}</td>
                      <td className="p-4"><Badge variant="secondary">{col.distance}</Badge></td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
