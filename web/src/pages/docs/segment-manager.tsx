import { Card, CardContent } from "../../components/ui/card"

export default function SegmentManagerPage() {
  return (
    <div className="space-y-8 max-w-3xl">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Segment Manager</h1>
        <p className="text-sm text-zinc-500 mt-1">Segment lifecycle, compaction, and search fan-out.</p>
      </div>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Segment Lifecycle</h2>
        <Card>
          <CardContent className="p-4">
            <pre className="text-xs font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md leading-relaxed">
{`Mutable → Seal → Immutable → Compaction → Merge → Delete

Each segment is a self-contained unit holding:
  - Dense vectors
  - Payload store
  - Metadata index`}
            </pre>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">States</h2>
        <Card>
          <CardContent className="p-0">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-zinc-200 dark:border-zinc-700">
                  <th className="text-left p-3 font-medium text-zinc-500">State</th>
                  <th className="text-left p-3 font-medium text-zinc-500">Writable</th>
                  <th className="text-left p-3 font-medium text-zinc-500">On Disk</th>
                  <th className="text-left p-3 font-medium text-zinc-500">Description</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-zinc-100 dark:divide-zinc-800">
                {[
                  ["Mutable", "Yes", "No", "Accepting new writes in memory. Data is not yet persisted to disk."],
                  ["Sealing", "No", "No", "Transitioning to immutable. Flushing data to disk."],
                  ["Immutable", "No", "Yes", "Read-only segment on disk. Searched via flat index loaded from dense.idx."],
                  ["Compacting", "No", "Yes", "Being merged into a new segment. Old files still exist until merge completes."],
                  ["Deleted", "No", "No", "Removed after successful compaction."],
                ].map(([state, writable, disk, desc]) => (
                  <tr key={state}>
                    <td className="p-3 font-mono text-xs text-zinc-900 dark:text-zinc-50">{state}</td>
                    <td className="p-3 text-xs">{writable === "Yes" ? <span className="text-emerald-600">✓</span> : <span className="text-zinc-400">—</span>}</td>
                    <td className="p-3 text-xs">{disk === "Yes" ? <span className="text-emerald-600">✓</span> : <span className="text-zinc-400">—</span>}</td>
                    <td className="p-3 text-xs text-zinc-600 dark:text-zinc-400">{desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Mutable Segment</h2>
        <Card>
          <CardContent className="p-4 space-y-3">
            <p className="text-sm text-zinc-600 dark:text-zinc-400">
              The active writable segment lives in memory. All inserts and deletes go to the mutable segment first.
              When it reaches the configured maxPoints threshold, it auto-seals to disk and a new mutable segment is created.
            </p>
            <p className="text-sm text-zinc-600 dark:text-zinc-400">
              Searches within the mutable segment use brute-force (flat scan) since the data is in memory.
              The segment also maintains a metadata engine for filter support.
            </p>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Seal (Mutable → Immutable)</h2>
        <Card>
          <CardContent className="p-4">
            <pre className="text-xs font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md leading-relaxed">
{`When a mutable segment reaches capacity:
1. State changes to Sealing (rejects new writes)
2. payload.dat written: [ID(8B) + payloadLen(4B) + payload(N)]
3. dense.idx written: [ID(8B) + vecLen(4B) + vec(N*4B)]
4. Directory: segments/segment_NNNNNN/
5. State changes to Immutable
6. New mutable segment created`}
            </pre>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Compaction</h2>
        <Card>
          <CardContent className="p-4 space-y-3">
            <p className="text-sm text-zinc-600 dark:text-zinc-400">
              When the number of immutable segments exceeds maxSegments, compaction triggers automatically.
              The compaction process:
            </p>
            <ol className="list-decimal list-inside text-sm text-zinc-600 dark:text-zinc-400 space-y-1">
              <li>Creates a new mutable segment to merge into</li>
              <li>Copies all vectors and payloads from each immutable segment</li>
              <li>Seals the merged segment to disk as a new immutable</li>
              <li>Replaces the immutable list with the single merged segment</li>
            </ol>
            <p className="text-sm text-zinc-500 mt-2">
              This is inspired by Lucene's segment merging approach. Compaction improves search performance
              by reducing the number of segments that need to be searched and queried.
            </p>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Search Flow</h2>
        <Card>
          <CardContent className="p-4">
            <pre className="text-xs font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md leading-relaxed">
{`SegmentManager.SearchAcross(query, topK):
  1. Search mutable segment (brute force in memory)
  2. For each immutable segment:
     Load flat index from dense.idx
     Brute-force search
  3. Merge results from all segments:
     If same ID appears multiple times, keep lowest score
  4. Sort by score ascending
  5. Return top-K`}
            </pre>
          </CardContent>
        </Card>
      </section>
    </div>
  )
}
