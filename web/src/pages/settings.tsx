import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card"

export default function SettingsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Settings</h1>
        <p className="text-sm text-zinc-500 mt-1">Instance configuration</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Card>
          <CardHeader><CardTitle className="text-sm">Storage</CardTitle></CardHeader>
          <CardContent className="space-y-3">
            <div className="flex justify-between items-center">
              <span className="text-sm text-zinc-600 dark:text-zinc-400">Collections Path</span>
              <code className="text-xs bg-zinc-100 dark:bg-zinc-800 px-2 py-1 rounded">/data/oris/collections</code>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader><CardTitle className="text-sm">Engine</CardTitle></CardHeader>
          <CardContent className="space-y-3">
            <div className="flex justify-between items-center">
              <span className="text-sm text-zinc-600 dark:text-zinc-400">SIMD</span>
              <span className="text-sm font-medium text-emerald-600">Enabled (NEON)</span>
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
              <span className="text-sm text-zinc-600 dark:text-zinc-400">Interval</span>
              <span className="text-sm font-medium text-zinc-900 dark:text-zinc-50">Every 1 hour</span>
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
              <span className="text-sm text-zinc-600 dark:text-zinc-400">Go Runtime</span>
              <span className="text-sm font-medium text-zinc-900 dark:text-zinc-50">1.26</span>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
