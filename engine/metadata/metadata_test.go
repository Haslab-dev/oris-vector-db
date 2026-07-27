package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexAndExists(t *testing.T) {
	e := New()
	e.IndexField(1, "country", "japan")
	e.IndexField(2, "country", "usa")
	e.IndexField(3, "country", "japan")

	bm, err := e.Evaluate(&Exists{Field: "country"})
	require.NoError(t, err)
	assert.Equal(t, uint64(3), bm.GetCardinality())
	assert.True(t, bm.Contains(uint32(1)))
	assert.True(t, bm.Contains(uint32(2)))
	assert.True(t, bm.Contains(uint32(3)))
}

func TestInFilter(t *testing.T) {
	e := New()
	e.IndexField(1, "country", "japan")
	e.IndexField(2, "country", "usa")
	e.IndexField(3, "country", "japan")

	bm, err := e.Evaluate(&In{Field: "country", Values: []string{"japan"}})
	require.NoError(t, err)
	assert.Equal(t, uint64(2), bm.GetCardinality())
	assert.True(t, bm.Contains(uint32(1)))
	assert.True(t, bm.Contains(uint32(3)))
}

func TestAndFilter(t *testing.T) {
	e := New()
	e.IndexField(1, "country", "japan")
	e.IndexField(1, "city", "tokyo")
	e.IndexField(2, "country", "japan")
	e.IndexField(2, "city", "osaka")
	e.IndexField(3, "country", "japan")
	e.IndexField(3, "city", "tokyo")

	bm, err := e.Evaluate(&And{
		Filters: []Filter{
			&In{Field: "country", Values: []string{"japan"}},
			&In{Field: "city", Values: []string{"tokyo"}},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, uint64(2), bm.GetCardinality())
	assert.True(t, bm.Contains(uint32(1)))
	assert.True(t, bm.Contains(uint32(3)))
}

func TestOrFilter(t *testing.T) {
	e := New()
	e.IndexField(1, "country", "japan")
	e.IndexField(2, "country", "usa")
	e.IndexField(3, "country", "france")

	bm, err := e.Evaluate(&Or{
		Filters: []Filter{
			&In{Field: "country", Values: []string{"japan"}},
			&In{Field: "country", Values: []string{"usa"}},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, uint64(2), bm.GetCardinality())
	assert.True(t, bm.Contains(uint32(1)))
	assert.True(t, bm.Contains(uint32(2)))
}

func TestNotFilter(t *testing.T) {
	e := New()
	e.IndexField(1, "country", "japan")
	e.IndexField(2, "country", "usa")

	bm, err := e.Evaluate(&Not{
		Filter: &In{Field: "country", Values: []string{"japan"}},
	})
	require.NoError(t, err)
	assert.Equal(t, uint64(1), bm.GetCardinality())
	assert.True(t, bm.Contains(uint32(2)))
}

func TestRangeFilter(t *testing.T) {
	e := New()
	e.IndexField(1, "price", float64(10))
	e.IndexField(2, "price", float64(20))
	e.IndexField(3, "price", float64(30))

	t.Run("min only", func(t *testing.T) {
		min := 15.0
		bm, err := e.Evaluate(&Range{Field: "price", Min: &min, Max: nil})
		require.NoError(t, err)
		assert.Equal(t, uint64(2), bm.GetCardinality())
		assert.True(t, bm.Contains(uint32(2)))
		assert.True(t, bm.Contains(uint32(3)))
	})

	t.Run("max only", func(t *testing.T) {
		max := 25.0
		bm, err := e.Evaluate(&Range{Field: "price", Min: nil, Max: &max})
		require.NoError(t, err)
		assert.Equal(t, uint64(2), bm.GetCardinality())
		assert.True(t, bm.Contains(uint32(1)))
		assert.True(t, bm.Contains(uint32(2)))
	})

	t.Run("both bounds", func(t *testing.T) {
		min, max := 5.0, 25.0
		bm, err := e.Evaluate(&Range{Field: "price", Min: &min, Max: &max})
		require.NoError(t, err)
		assert.Equal(t, uint64(2), bm.GetCardinality())
		assert.True(t, bm.Contains(uint32(1)))
		assert.True(t, bm.Contains(uint32(2)))
	})
}

func TestRemoveField(t *testing.T) {
	e := New()
	e.IndexField(1, "country", "japan")
	e.IndexField(2, "country", "usa")

	e.RemoveField(1, "country", "japan")

	bm, err := e.Evaluate(&In{Field: "country", Values: []string{"japan"}})
	require.NoError(t, err)
	assert.Equal(t, uint64(0), bm.GetCardinality())
}

func TestNestedFilters(t *testing.T) {
	e := New()
	e.IndexField(1, "country", "japan")
	e.IndexField(1, "city", "tokyo")
	e.IndexField(1, "price", float64(100))
	e.IndexField(2, "country", "japan")
	e.IndexField(2, "city", "osaka")
	e.IndexField(2, "price", float64(50))
	e.IndexField(3, "country", "usa")
	e.IndexField(3, "city", "nyc")
	e.IndexField(3, "price", float64(200))

	// japan AND (tokyo OR osaka) AND price < 150
	priceMax := 150.0
	bm, err := e.Evaluate(&And{
		Filters: []Filter{
			&In{Field: "country", Values: []string{"japan"}},
			&Or{
				Filters: []Filter{
					&In{Field: "city", Values: []string{"tokyo"}},
					&In{Field: "city", Values: []string{"osaka"}},
				},
			},
			&Range{Field: "price", Min: nil, Max: &priceMax},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, uint64(2), bm.GetCardinality())
	assert.True(t, bm.Contains(uint32(1)))
	assert.True(t, bm.Contains(uint32(2)))
}

func TestEmptyIndex(t *testing.T) {
	e := New()

	bm, err := e.Evaluate(&In{Field: "nonexistent", Values: []string{"x"}})
	require.NoError(t, err)
	assert.Equal(t, uint64(0), bm.GetCardinality())

	bm, err = e.Evaluate(&Exists{Field: "nonexistent"})
	require.NoError(t, err)
	assert.Equal(t, uint64(0), bm.GetCardinality())
}

func TestAllDocs(t *testing.T) {
	e := New()
	e.IndexField(1, "a", "x")
	e.IndexField(2, "a", "y")
	e.IndexField(3, "a", "z")

	assert.Equal(t, uint64(3), e.AllDocs().GetCardinality())
}

func TestEvaluateEmptyAnd(t *testing.T) {
	e := New()
	e.IndexField(1, "a", "x")

	bm, err := e.Evaluate(&And{Filters: nil})
	require.NoError(t, err)
	assert.Equal(t, uint64(1), bm.GetCardinality())
}

func TestEvaluateEmptyOr(t *testing.T) {
	e := New()
	bm, err := e.Evaluate(&Or{Filters: nil})
	require.NoError(t, err)
	assert.Equal(t, uint64(0), bm.GetCardinality())
}
