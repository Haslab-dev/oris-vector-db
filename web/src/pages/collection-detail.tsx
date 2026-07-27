import { useEffect, useState } from "react"
import { useParams } from "@tanstack/react-router"
import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card"
import { Badge } from "../components/ui/badge"
import { fetchCollection, fetchCount, type CollectionInfo } from "../lib/api"

export default function CollectionDetailPage() {
  const { name } = useParams({ from: "/collections/$name" })
  const [info, setInfo] = useState<CollectionInfo | null>(null)
  const [count, setCount] = useState(0)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    Promise.all([
      fetchCollection(name),
      fetchCount(name),
    ])
      .then(([i, c]) => {
        setInfo(i)
        setCount(c)
      })
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [name])

  if (loading) return <div className="text-sm text-zinc-500">Loading...</div>
  if (!info) return <div className="text-sm text-red-500">Collection not found</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">{name}</h1>
        <Badge>Active</Badge>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {[
          { label: "Vectors", value: count.toLocaleString() },
          { label: "Dimension", value: String(info.dimension) },
          { label: "Distance", value: info.distance },
          { label: "Segments", value: String(info.segments ?? "—") },
        ].map((s) => (
          <Card key={s.label}>
            <CardHeader className="pb-1">
              <CardTitle className="text-xs text-zinc-500 font-medium">{s.label}</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-xl font-bold text-zinc-900 dark:text-zinc-50">{s.value}</div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
