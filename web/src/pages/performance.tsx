import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card"
import { BarChart, Bar, XAxis, YAxis, ResponsiveContainer, AreaChart, Area } from "recharts"

const latencyData = [
  { name: "00:00", value: 4.2 }, { name: "04:00", value: 3.8 }, { name: "08:00", value: 7.1 },
  { name: "12:00", value: 8.5 }, { name: "16:00", value: 6.9 }, { name: "20:00", value: 5.3 },
]

const insertData = [
  { name: "Mon", value: 1200 }, { name: "Tue", value: 980 }, { name: "Wed", value: 1500 },
  { name: "Thu", value: 1100 }, { name: "Fri", value: 1350 }, { name: "Sat", value: 700 }, { name: "Sun", value: 450 },
]

export default function PerformancePage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Performance</h1>
        <p className="text-sm text-zinc-500 mt-1">Engine metrics and throughput</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Card>
          <CardHeader><CardTitle className="text-sm">Search Latency (ms)</CardTitle></CardHeader>
          <CardContent>
            <div className="h-48">
              <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={latencyData}>
                  <Area type="monotone" dataKey="value" stroke="#18181b" fill="#18181b" fillOpacity={0.1} strokeWidth={2} />
                  <XAxis dataKey="name" tick={{ fontSize: 11 }} axisLine={false} tickLine={false} />
                  <YAxis tick={{ fontSize: 11 }} axisLine={false} tickLine={false} />
                </AreaChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader><CardTitle className="text-sm">Inserts / sec</CardTitle></CardHeader>
          <CardContent>
            <div className="h-48">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={insertData}>
                  <Bar dataKey="value" fill="#18181b" radius={[4, 4, 0, 0]} />
                  <XAxis dataKey="name" tick={{ fontSize: 11 }} axisLine={false} tickLine={false} />
                  <YAxis tick={{ fontSize: 11 }} axisLine={false} tickLine={false} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader><CardTitle className="text-sm">Memory</CardTitle></CardHeader>
          <CardContent>
            <div className="h-32 flex items-end gap-2">
              {[85, 92, 78, 88, 95, 87, 82].map((v, i) => (
                <div key={i} className="flex-1 bg-zinc-200 dark:bg-zinc-700 rounded-t-md" style={{ height: `${v}%` }} />
              ))}
            </div>
            <p className="text-xs text-zinc-500 mt-2">280 MB current · 320 MB peak</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader><CardTitle className="text-sm">Storage</CardTitle></CardHeader>
          <CardContent>
            <div className="h-32 flex items-end gap-2">
              {[60, 65, 70, 75, 80, 85, 90].map((v, i) => (
                <div key={i} className="flex-1 bg-emerald-200 dark:bg-emerald-800 rounded-t-md" style={{ height: `${v}%` }} />
              ))}
            </div>
            <p className="text-xs text-zinc-500 mt-2">512 MB used · 1 TB available</p>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
