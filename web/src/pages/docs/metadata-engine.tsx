import { Card, CardContent } from "../../components/ui/card"

export default function MetadataEnginePage() {
  return (
    <div className="space-y-8 max-w-3xl">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Metadata Engine</h1>
        <p className="text-sm text-zinc-500 mt-1">Filter expressions evaluated against indexed metadata using Roaring Bitmaps.</p>
      </div>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Filter AST</h2>
        <p className="text-sm text-zinc-600 dark:text-zinc-400">
          Filters are expressed as an AST (Abstract Syntax Tree) of composable filter nodes.
          Each filter evaluates to a Roaring Bitmap of matching document IDs.
        </p>
        <Card>
          <CardContent className="p-0">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-zinc-200 dark:border-zinc-700">
                  <th className="text-left p-3 font-medium text-zinc-500">Filter</th>
                  <th className="text-left p-3 font-medium text-zinc-500">Operation</th>
                  <th className="text-left p-3 font-medium text-zinc-500">Example</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-zinc-100 dark:divide-zinc-800">
                {[
                  ["AND", "Bitmap intersection", `metadata.And{
  Filters: []metadata.Filter{filter1, filter2}
}`],
                  ["OR", "Bitmap union", `metadata.Or{
  Filters: []metadata.Filter{filter1, filter2}
}`],
                  ["NOT", "Bitmap negation (AllDocs \\ filter)", `metadata.Not{
  Filter: filter
}`],
                  ["IN", "Multi-value bitmap union", `metadata.In{
  Field: "country",
  Values: []string{"japan", "usa"}
}`],
                  ["Range", "Numeric range scan", `metadata.Range{
  Field: "price",
  Min: &min, Max: &max
}`],
                  ["Exists", "Bitmap of docs with field", `metadata.Exists{
  Field: "category"
}`],
                ].map(([name, op, example]) => (
                  <tr key={name}>
                    <td className="p-3 font-mono text-xs text-zinc-900 dark:text-zinc-50">{name}</td>
                    <td className="p-3 text-xs text-zinc-600 dark:text-zinc-400">{op}</td>
                    <td className="p-3">
                      <pre className="text-xs font-mono bg-zinc-100 dark:bg-zinc-800 p-2 rounded">{example}</pre>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">How It Works</h2>
        <Card>
          <CardContent className="p-4 space-y-3">
            <div className="space-y-2">
              <h3 className="text-sm font-medium text-zinc-900 dark:text-zinc-50">Indexing</h3>
              <p className="text-sm text-zinc-600 dark:text-zinc-400">
                When a document is indexed with a metadata field, three indexes are updated:
              </p>
              <ul className="list-disc list-inside text-sm text-zinc-600 dark:text-zinc-400 space-y-1">
                <li><strong className="text-zinc-900 dark:text-zinc-50">String field index</strong>: maps (field, value) → bitmap of doc IDs. Used for IN and equality filters.</li>
                <li><strong className="text-zinc-900 dark:text-zinc-50">Numeric index</strong>: sorted list of (value, docID) pairs. Used for Range filters on numeric fields.</li>
                <li><strong className="text-zinc-900 dark:text-zinc-50">Exists index</strong>: maps field → bitmap of doc IDs that have the field. Used for Exists filters.</li>
              </ul>
            </div>
            <div className="space-y-2">
              <h3 className="text-sm font-medium text-zinc-900 dark:text-zinc-50">Evaluation</h3>
              <p className="text-sm text-zinc-600 dark:text-zinc-400">
                Filter evaluation produces a Roaring Bitmap that can be intersected with search results.
                The Planner applies filters as a post-filter step: results that don't match the filter bitmap are removed.
              </p>
            </div>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Usage Example</h2>
        <Card>
          <CardContent className="p-4">
            <pre className="text-sm font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md overflow-x-auto">
{`import "github.com/hasdev/oris/engine/metadata"

meta := metadata.New()

// Index metadata.
meta.IndexField(1, "country", "japan")
meta.IndexField(1, "price", float64(100))
meta.IndexField(2, "country", "usa")

// Evaluate: country = japan AND price < 50.
priceMax := 50.0
result, _ := meta.Evaluate(&metadata.And{
  Filters: []metadata.Filter{
    &metadata.In{Field: "country", Values: []string{"japan"}},
    &metadata.Range{Field: "price", Min: nil, Max: &priceMax},
  },
})

// result is a Roaring Bitmap of matching doc IDs.`}
            </pre>
          </CardContent>
        </Card>
      </section>
    </div>
  )
}
