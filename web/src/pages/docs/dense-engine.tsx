import { Card, CardContent } from "../../components/ui/card"

export default function DenseEnginePage() {
  return (
    <div className="space-y-8 max-w-3xl">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Dense Engine (HNSW)</h1>
        <p className="text-sm text-zinc-500 mt-1">Hierarchical Navigable Small World graphs for approximate nearest neighbor search.</p>
      </div>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Overview</h2>
        <p className="text-sm text-zinc-600 dark:text-zinc-400">
          HNSW builds a multi-layer graph where each layer is a progressively sparser approximation of the full dataset.
          The top layer contains only a few nodes (long-range connections), while the bottom layer contains all nodes with
          fine-grained local connections. Search starts at the top and greedily descends, finding the nearest neighbors
          at each layer before proceeding to the next.
        </p>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Configuration</h2>
        <Card>
          <CardContent className="p-0">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-zinc-200 dark:border-zinc-700">
                  <th className="text-left p-3 font-medium text-zinc-500">Parameter</th>
                  <th className="text-left p-3 font-medium text-zinc-500">Default</th>
                  <th className="text-left p-3 font-medium text-zinc-500">Description</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-zinc-100 dark:divide-zinc-800">
                {[
                  ["M", "16", "Number of bi-directional connections per layer"],
                  ["efConstruction", "200", "Dynamic candidate list size during graph construction (higher = better recall, slower build)"],
                  ["efSearch", "100", "Dynamic candidate list size during search (higher = better recall, slower search)"],
                  ["Dimension", "—", "Vector dimension (all vectors must match)"],
                  ["Distance", "cosine", "Distance metric: cosine, dot, or euclidean"],
                ].map(([param, def, desc]) => (
                  <tr key={param}>
                    <td className="p-3 font-mono text-xs text-zinc-900 dark:text-zinc-50">{param}</td>
                    <td className="p-3 font-mono text-xs text-zinc-600 dark:text-zinc-400">{def}</td>
                    <td className="p-3 text-xs text-zinc-600 dark:text-zinc-400">{desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">How Insert Works</h2>
        <Card>
          <CardContent className="p-4 space-y-3">
            <div className="space-y-2">
              <h3 className="text-sm font-medium text-zinc-900 dark:text-zinc-50">1. Random Level Assignment</h3>
              <p className="text-sm text-zinc-600 dark:text-zinc-400">
                Each new node is assigned a random level using the distribution floor(-ln(uniform) × mL) where mL = 1/ln(M).
                This creates a logarithmic number of layers, with most nodes on the bottom layer.
              </p>
            </div>
            <div className="space-y-2">
              <h3 className="text-sm font-medium text-zinc-900 dark:text-zinc-50">2. Top-Level Descent</h3>
              <p className="text-sm text-zinc-600 dark:text-zinc-400">
                Starting from the entry point, greedily traverse from the top layer down to the new node's level,
                finding the single nearest neighbor at each layer.
              </p>
            </div>
            <div className="space-y-2">
              <h3 className="text-sm font-medium text-zinc-900 dark:text-zinc-50">3. efConstruction Search</h3>
              <p className="text-sm text-zinc-600 dark:text-zinc-400">
                At each layer from the insertion level down to 0, find efConstruction nearest neighbors using the
                searchLayer algorithm. This produces a candidate list of close nodes.
              </p>
            </div>
            <div className="space-y-2">
              <h3 className="text-sm font-medium text-zinc-900 dark:text-zinc-50">4. Connect and Shrink</h3>
              <p className="text-sm text-zinc-600 dark:text-zinc-400">
                Select the M closest candidates and create bi-directional edges. If a neighbor now exceeds Mmax
                connections, it shrinks its edge list by keeping only the closest Mmax connections.
              </p>
            </div>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">How Search Works</h2>
        <Card>
          <CardContent className="p-4">
            <pre className="text-xs font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md overflow-x-auto leading-relaxed">
{`searchLayer(entry, query, ef, layer):
  1. Initialize:
     - candidates = min-heap (closest on top), seeded with entry point
     - results = max-heap (farthest on top), seeded with entry point
     - visited = {entry}

  2. While candidates is not empty:
     Pop candidate from candidates (closest first)
     If candidate.dist > results[0].dist (farthest result):
       break  // no better candidates remain

     For each neighbor of candidate in this layer:
       If not visited:
         Mark visited
         Compute distance to query
         If result set not full OR distance < farthest result:
           Add to both heaps
           If results exceeds ef, pop farthest

  3. Return result set`}
            </pre>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Flat Backend</h2>
        <p className="text-sm text-zinc-600 dark:text-zinc-400">
          For collections with fewer than 1,000 vectors, the factory auto-selects a Flat (brute-force) index instead of HNSW.
          Flat linearly scans all vectors and computes exact distances. This avoids HNSW overhead for small datasets where
          exact search is faster than approximate.
        </p>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Performance Benchmarks</h2>
        <Card>
          <CardContent className="p-4">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-zinc-200 dark:border-zinc-700">
                  <th className="text-left pb-2 font-medium text-zinc-500">Benchmark</th>
                  <th className="text-right pb-2 font-medium text-zinc-500">Result</th>
                  <th className="text-right pb-2 font-medium text-zinc-500">Platform</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-zinc-100 dark:divide-zinc-800">
                {[
                  ["Insert (128-dim)", "120,239 ns/op", "Apple M3"],
                  ["Search 10K (128-dim, top-10)", "528,510 ns/op", "Apple M3"],
                  ["Cosine 128-dim", "~1.5 ns/op", "Apple M3 (NEON)"],
                ].map(([bench, result, platform]) => (
                  <tr key={bench}>
                    <td className="py-2 text-zinc-900 dark:text-zinc-50">{bench}</td>
                    <td className="py-2 text-right font-mono text-xs text-zinc-600 dark:text-zinc-400">{result}</td>
                    <td className="py-2 text-right text-xs text-zinc-500">{platform}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
        </Card>
      </section>
    </div>
  )
}
