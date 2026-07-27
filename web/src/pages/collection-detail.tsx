import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card"
import { Badge } from "../components/ui/badge"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@radix-ui/react-tabs"
import { cn } from "../lib/utils"

const segments = [
  { id: "#1", status: "Immutable" as const, vectors: "250K", size: "120 MB" },
  { id: "#2", status: "Immutable" as const, vectors: "310K", size: "140 MB" },
  { id: "#3", status: "Mutable" as const, vectors: "25K", size: "15 MB" },
]

export default function CollectionDetailPage() {
  return (
    <div className="space-y-6">
      <div>
        <div className="flex items-center gap-3">
          <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Products</h1>
          <Badge>Active</Badge>
        </div>
        <p className="text-sm text-zinc-500 mt-1">Collection details and data</p>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {[
          { label: "Vectors", value: "1,245,221" },
          { label: "Dimension", value: "768" },
          { label: "Distance", value: "Cosine" },
          { label: "Segments", value: "12" },
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

      <Tabs defaultValue="overview" className="space-y-4">
        <TabsList className="flex gap-1 border-b border-zinc-200 dark:border-zinc-800">
          {["overview", "search", "data", "settings"].map((tab) => (
            <TabsTrigger
              key={tab}
              value={tab}
              className={cn(
                "px-4 py-2 text-sm font-medium text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-50 data-[state=active]:text-zinc-900 dark:data-[state=active]:text-zinc-50 data-[state=active]:border-b-2 data-[state=active]:border-zinc-900 dark:data-[state=active]:border-zinc-50 capitalize cursor-pointer",
              )}
            >
              {tab}
            </TabsTrigger>
          ))}
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          <Card>
            <CardHeader><CardTitle>Segments</CardTitle></CardHeader>
            <CardContent className="p-0">
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-zinc-200 dark:border-zinc-700 text-left">
                      <th className="p-4 font-medium text-zinc-500">Segment</th>
                      <th className="p-4 font-medium text-zinc-500">Status</th>
                      <th className="p-4 font-medium text-zinc-500 text-right">Vectors</th>
                      <th className="p-4 font-medium text-zinc-500 text-right">Size</th>
                    </tr>
                  </thead>
                  <tbody>
                    {segments.map((seg) => (
                      <tr key={seg.id} className="border-b border-zinc-100 dark:border-zinc-800">
                        <td className="p-4 font-medium text-zinc-900 dark:text-zinc-50">{seg.id}</td>
                        <td className="p-4">
                          <Badge variant={seg.status === "Immutable" ? "green" : "amber"}>{seg.status}</Badge>
                        </td>
                        <td className="p-4 text-right text-zinc-600 dark:text-zinc-400">{seg.vectors}</td>
                        <td className="p-4 text-right text-zinc-600 dark:text-zinc-400">{seg.size}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="search">
          <p className="text-sm text-zinc-500">Search playground — navigate to the Search page.</p>
        </TabsContent>

        <TabsContent value="data">
          <p className="text-sm text-zinc-500">Data viewer — coming soon.</p>
        </TabsContent>

        <TabsContent value="settings">
          <p className="text-sm text-zinc-500">Collection settings — coming soon.</p>
        </TabsContent>
      </Tabs>
    </div>
  )
}
