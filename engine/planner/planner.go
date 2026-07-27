// Package planner orchestrates query execution across all Oris engines.
package planner

import (
	"sort"

	"github.com/hasdev/oris/engine/dense"
	"github.com/hasdev/oris/engine/metadata"
	"github.com/hasdev/oris/engine/sparse"
)

// Mode represents the query execution mode.
type Mode int

const (
	DenseOnly  Mode = iota
	SparseOnly
	Hybrid
)

// Query is a complete retrieval query.
type Query struct {
	DenseVector  []float32
	SparseTokens []string
	Filter       metadata.Filter
	TopK         int
	Mode         Mode
	Alpha        float64
}

// Result is a single combined search result.
type Result struct {
	ID          uint64
	DenseScore  float32
	SparseScore float32
	FinalScore  float32
}

// Planner executes queries across dense, sparse, and metadata engines.
type Planner struct {
	dense    dense.Engine
	sparse   sparse.Engine
	meta     *metadata.Engine
	denseDim int
}

// New creates a planner.
func New(denseEngine dense.Engine, sparseEngine sparse.Engine, metaEngine *metadata.Engine, denseDim int) *Planner {
	return &Planner{
		dense:    denseEngine,
		sparse:   sparseEngine,
		meta:     metaEngine,
		denseDim: denseDim,
	}
}

// Execute runs a query and returns ranked results.
func (p *Planner) Execute(q Query) ([]Result, error) {
	if q.TopK <= 0 {
		q.TopK = 10
	}
	if q.Alpha == 0 {
		q.Alpha = 0.5
	}
	switch q.Mode {
	case DenseOnly:
		return p.denseOnly(q)
	case SparseOnly:
		return p.sparseOnly(q)
	default:
		return p.hybrid(q)
	}
}

func (p *Planner) denseOnly(q Query) ([]Result, error) {
	if p.dense == nil || q.DenseVector == nil {
		return nil, nil
	}
	results, err := p.dense.Search(q.DenseVector, q.TopK)
	if err != nil {
		return nil, err
	}
	if q.Filter != nil {
		results = p.filterDense(results, q.Filter)
		if len(results) > q.TopK {
			results = results[:q.TopK]
		}
	}
	out := make([]Result, len(results))
	for i, r := range results {
		out[i] = Result{ID: r.ID, DenseScore: r.Score, FinalScore: r.Score}
	}
	return out, nil
}

func (p *Planner) sparseOnly(q Query) ([]Result, error) {
	if p.sparse == nil || len(q.SparseTokens) == 0 {
		return nil, nil
	}
	results, err := p.sparse.Search(q.SparseTokens, q.TopK)
	if err != nil {
		return nil, err
	}
	if q.Filter != nil {
		results = p.filterSparse(results, q.Filter)
		if len(results) > q.TopK {
			results = results[:q.TopK]
		}
	}
	out := make([]Result, len(results))
	for i, r := range results {
		out[i] = Result{ID: r.ID, SparseScore: r.Score, FinalScore: r.Score}
	}
	return out, nil
}

type hybridEntry struct {
	id     uint64
	dense  float32
	sparse float32
	final  float32
}

func (p *Planner) hybrid(q Query) ([]Result, error) {
	if p.dense == nil || q.DenseVector == nil || p.sparse == nil || len(q.SparseTokens) == 0 {
		if p.dense != nil && q.DenseVector != nil {
			return p.denseOnly(q)
		}
		return p.sparseOnly(q)
	}

	mult := 3
	denseResults, err := p.dense.Search(q.DenseVector, q.TopK*mult)
	if err != nil {
		return nil, err
	}
	sparseResults, err := p.sparse.Search(q.SparseTokens, q.TopK*mult)
	if err != nil {
		return nil, err
	}

	normDense := normalizeDense(denseResults)
	normSparse := normalizeSparse(sparseResults)

	alpha := q.Alpha
	merged := make(map[uint64]*hybridEntry)

	for _, r := range normDense {
		merged[r.ID] = &hybridEntry{id: r.ID, dense: r.Score}
	}
	for _, r := range normSparse {
		if e, ok := merged[r.ID]; ok {
			e.sparse = r.Score
		} else {
			merged[r.ID] = &hybridEntry{id: r.ID, sparse: r.Score}
		}
	}

	results := make([]hybridEntry, 0, len(merged))
	for _, e := range merged {
		e.final = float32(float64(e.dense)*alpha + float64(e.sparse)*(1-alpha))
		results = append(results, *e)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].final < results[j].final
	})
	if len(results) > q.TopK {
		results = results[:q.TopK]
	}
	if q.Filter != nil {
		results = p.filterHybrid(results, q.Filter)
		if len(results) > q.TopK {
			results = results[:q.TopK]
		}
	}

	out := make([]Result, len(results))
	for i, r := range results {
		out[i] = Result{ID: r.id, DenseScore: r.dense, SparseScore: r.sparse, FinalScore: r.final}
	}
	return out, nil
}

func (p *Planner) filterDense(r []dense.Result, f metadata.Filter) []dense.Result {
	bm, err := p.meta.Evaluate(f)
	if err != nil || bm == nil {
		return r
	}
	keep := make([]dense.Result, 0, len(r))
	for _, v := range r {
		if bm.Contains(uint32(v.ID)) {
			keep = append(keep, v)
		}
	}
	return keep
}

func (p *Planner) filterSparse(r []sparse.Result, f metadata.Filter) []sparse.Result {
	bm, err := p.meta.Evaluate(f)
	if err != nil || bm == nil {
		return r
	}
	keep := make([]sparse.Result, 0, len(r))
	for _, v := range r {
		if bm.Contains(uint32(v.ID)) {
			keep = append(keep, v)
		}
	}
	return keep
}

func (p *Planner) filterHybrid(r []hybridEntry, f metadata.Filter) []hybridEntry {
	bm, err := p.meta.Evaluate(f)
	if err != nil || bm == nil {
		return r
	}
	keep := make([]hybridEntry, 0, len(r))
	for _, v := range r {
		if bm.Contains(uint32(v.id)) {
			keep = append(keep, v)
		}
	}
	return keep
}

func normalizeDense(r []dense.Result) []dense.Result {
	if len(r) == 0 {
		return r
	}
	minS, maxS := r[0].Score, r[0].Score
	for _, v := range r {
		if v.Score > maxS {
			maxS = v.Score
		}
		if v.Score < minS {
			minS = v.Score
		}
	}
	rn := maxS - minS
	if rn == 0 {
		return r
	}
	for i := range r {
		r[i].Score = (r[i].Score - minS) / rn
	}
	return r
}

func normalizeSparse(r []sparse.Result) []sparse.Result {
	if len(r) == 0 {
		return r
	}
	minS, maxS := r[0].Score, r[0].Score
	for _, v := range r {
		if v.Score > maxS {
			maxS = v.Score
		}
		if v.Score < minS {
			minS = v.Score
		}
	}
	rn := maxS - minS
	if rn == 0 {
		return r
	}
	for i := range r {
		r[i].Score = (r[i].Score - minS) / rn
	}
	return r
}
