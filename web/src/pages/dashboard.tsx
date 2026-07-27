import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card"
import { Database, HardDrive, MemoryStick } from "lucide-react"

const stats = [
  { icon: Database, label: "Collections", value: "8", change: "+2 this week" },
  { icon: Database, label: "Vectors", value: "1,245,322", change: "+50K today" },
  { icon: HardDrive, label: "Storage", value: "512 MB", change: "12 segments" },
  { icon: MemoryStick, label: "Memory", value: "280 MB", change: "Peak 320 MB" },
]

const recentCollections = [
  { name: "Products", vectors: "1.2M", dimension: 768, distance: "Cosine" },
  { name: "Articles", vectors: "250K", dimension: 1536, distance: "Cosine" },
  { name: "Users", vectors: "85K", dimension: 128, distance: "Dot" },
  { name: "Images", vectors: "45K", dimension: 512, distance: "Cosine" },
]

export default function DashboardPage() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Dashboard</h1>
        <p className="text-sm text-zinc-500 mt-1">Overview of your Oris instance</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat) => (
          <Card key={stat.label}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-zinc-500">{stat.label}</CardTitle>
              <stat.icon className="h-4 w-4 text-zinc-400" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">{stat.value}</div>
              <p className="text-xs text-zinc-500 mt-1">{stat.change}</p>
            </CardContent>
          </Card>
        ))}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Recent Collections</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-zinc-200 dark:border-zinc-700 text-left">
                  <th className="pb-3 font-medium text-zinc-500">Name</th>
                  <th className="pb-3 font-medium text-zinc-500 text-right">Vectors</th>
                  <th className="pb-3 font-medium text-zinc-500 text-right">Dimension</th>
                  <th className="pb-3 font-medium text-zinc-500">Distance</th>
                </tr>
              </thead>
              <tbody>
                {recentCollections.map((col) => (
                  <tr key={col.name} className="border-b border-zinc-100 dark:border-zinc-800">
                    <td className="py-3 font-medium text-zinc-900 dark:text-zinc-50">{col.name}</td>
                    <td className="py-3 text-right text-zinc-600 dark:text-zinc-400">{col.vectors}</td>
                    <td className="py-3 text-right text-zinc-600 dark:text-zinc-400">{col.dimension}</td>
                    <td className="py-3 text-zinc-600 dark:text-zinc-400">{col.distance}</td>
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
