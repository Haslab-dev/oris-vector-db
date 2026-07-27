import { Card, CardContent } from "../../components/ui/card"

export default function StorageEnginePage() {
  return (
    <div className="space-y-8 max-w-3xl">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Storage Engine</h1>
        <p className="text-sm text-zinc-500 mt-1">Persistence layer with WAL, snapshots, and pluggable backends.</p>
      </div>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Storage Interface</h2>
        <p className="text-sm text-zinc-600 dark:text-zinc-400">
          All storage operations go through a simple key-value interface, making it possible to swap backends without changing the retrieval engine.
        </p>
        <Card>
          <CardContent className="p-4">
            <pre className="text-sm font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md overflow-x-auto">
{`type Storage interface {
    Get(key []byte) ([]byte, error)
    Set(key, value []byte) error
    Delete(key []byte) error
    NewBatch() Batch
    Snapshot() (Snapshot, error)
    Close() error
}

type Batch interface {
    Set(key, value []byte)
    Delete(key []byte)
    Commit() error
    Close()
}

type Snapshot interface {
    Get(key []byte) ([]byte, error)
    NewIterator(start, end []byte) (Iterator, error)
    Close()
}`}
            </pre>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Pebble Backend</h2>
        <Card>
          <CardContent className="p-4 space-y-3">
            <p className="text-sm text-zinc-600 dark:text-zinc-400">
              The default storage backend uses <strong className="text-zinc-900 dark:text-zinc-50">Pebble</strong> (cockroachdb/pebble), a pure Go LSM tree key-value store.
              Pebble provides:
            </p>
            <ul className="list-disc list-inside text-sm text-zinc-600 dark:text-zinc-400 space-y-1">
              <li>High write throughput via LSM-tree architecture</li>
              <li>Write-ahead logging for crash safety</li>
              <li>Point-in-time snapshots</li>
              <li>Atomic batches</li>
              <li>Zero CGO dependencies (pure Go)</li>
              <li>RocksDB-compatible concepts and compression</li>
            </ul>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">In-Memory Backend</h2>
        <p className="text-sm text-zinc-600 dark:text-zinc-400">
          A thread-safe in-memory implementation is provided for testing. It implements the same Storage interface with
          concurrent read/write support via sync.RWMutex. All data is lost on Close().
        </p>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Write-Ahead Log (WAL)</h2>
        <Card>
          <CardContent className="p-4 space-y-3">
            <p className="text-sm text-zinc-600 dark:text-zinc-400">
              The WAL ensures durability by recording every write operation before it's applied to storage.
              On open, any uncommitted WAL entries are replayed to restore consistency, then the WAL is truncated.
            </p>
            <pre className="text-xs font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md overflow-x-auto">
{`WAL Entry Format:
┌─────┬──────────┬──────────┬──────────┬──────────────┬──────────┐
│ Op  │ Sequence │ Key Len  │   Key    │  Value Len   │  Value   │
│ 1B  │   8B     │   4B     │   K B    │     4B       │   V B    │
└─────┴──────────┴──────────┴──────────┴──────────────┴──────────┘

Op codes: 0x01 = Set, 0x02 = Delete`}
            </pre>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Snapshots</h2>
        <Card>
          <CardContent className="p-4 space-y-3">
            <p className="text-sm text-zinc-600 dark:text-zinc-400">
              Snapshotting creates a point-in-time copy of all key-value pairs, verified with a BLAKE3 checksum.
              Snapshots can be restored by replaying all key-value pairs through the Storage interface.
            </p>
            <pre className="text-sm font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md overflow-x-auto">
{`snap, _ := snapshot.Create(store, "./snapshots", "pre-upgrade")

// Verify integrity.
ok, _ := snapshot.Verify(snapPath, snap.Checksum)

// Restore to a new store.
snapshot.Restore(newStore, snapPath)`}
            </pre>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Store Layout</h2>
        <Card>
          <CardContent className="p-4">
            <pre className="text-xs font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md overflow-x-auto leading-relaxed">
{`Storage namespaces (key prefix):
  0x00 → Point keys: point ID → encoded point data
          Key format: 0x00 + uint64 LE ID (9 bytes total)

The Pebble database is organized as:
  {path}/data/      → Main collection data (Pebble DB)
  {path}/wal/       → WAL data (separate Pebble DB)
  {path}/snapshots/ → Snapshot files (.snap + BLAKE3 checksum)`}
            </pre>
          </CardContent>
        </Card>
      </section>
    </div>
  )
}
