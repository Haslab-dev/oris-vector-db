import { useEffect, useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card"
import { fetchStats, type StatsResult } from "../lib/api"

export default function SettingsPage() {
  const [stats, setStats] = useState<StatsResult | null>(null)

  useEffect(() => {
    fetchStats().then(setStats).catch(() => {})
  }, [])

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Settings</h1>
        <p className="text-sm text-zinc-500 mt-1">Instance configuration</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Card>
          <CardHeader><CardTitle className="text-sm">Instance</CardTitle></CardHeader>
          <CardContent className="space-y-3">
            <div className="flex justify-between items-center">
              <span className="text-sm text-zinc-600 dark:text-zinc-400">Collections</span>
              <span className="text-sm font-medium text-zinc-900 dark:text-zinc-50">{stats?.count ?? "—"}</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-zinc-600 dark:text-zinc-400">Total Vectors</span>
              <span className="text-sm font-medium text-zinc-900 dark:text-zinc-50">{stats?.totalVectors.toLocaleString() ?? "—"}</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-zinc-600 dark:text-zinc-400">Storage</span>
              <span className="text-sm font-medium text-zinc-900 dark:text-zinc-50">{stats ? `${stats.storageMB.toFixed(1)} MB` : "—"}</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader><CardTitle className="text-sm">Engine</CardTitle></CardHeader>
          <CardContent className="space-y-3">
            <div className="flex justify-between items-center">
              <span className="text-sm text-zinc-600 dark:text-zinc-400">SIMD</span>
              <span className="text-sm font-medium text-emerald-600">Enabled (auto-detect)</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-zinc-600 dark:text-zinc-400">Auto Flush</span>
              <span className="text-sm font-medium text-zinc-900 dark:text-zinc-50">Enabled</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader><CardTitle className="text-sm">Snapshots</CardTitle></CardHeader>
          <CardContent className="space-y-3">
            <div className="flex justify-between items-center">
              <span className="text-sm text-zinc-600 dark:text-zinc-400">Snapshot</span>
              <span className="text-sm font-medium text-zinc-900 dark:text-zinc-50">POST /api/snapshot</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader><CardTitle className="text-sm">System</CardTitle></CardHeader>
          <CardContent className="space-y-3">
            <div className="flex justify-between items-center">
              <span className="text-sm text-zinc-600 dark:text-zinc-400">Version</span>
              <span className="text-sm font-medium text-zinc-900 dark:text-zinc-50">0.1.0</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-zinc-600 dark:text-zinc-400">Runtime</span>
              <span className="text-sm font-medium text-zinc-900 dark:text-zinc-50">Go</span>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
