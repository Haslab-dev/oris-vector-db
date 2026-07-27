package main

import (
	"fmt"
	"os"

	"github.com/hasdev/oris/api"
)

func main() {
	// Example: create/open a collection, insert points, and search.
	path := "./example-data"
	defer os.RemoveAll(path)

	cfg := api.DefaultConfig("my-collection", 128)
	col, err := api.Open(path, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening collection: %v\n", err)
		os.Exit(1)
	}
	defer col.Close()

	// Insert some random-ish points.
	for i := uint64(0); i < 100; i++ {
		vec := make([]float32, 128)
		vec[i%128] = float32(i) / 100.0
		if err := col.Insert(i, vec, nil, nil, []byte(fmt.Sprintf("point-%d", i))); err != nil {
			fmt.Fprintf(os.Stderr, "error inserting: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Inserted %d points\n", col.Count())

	// Search.
	query := make([]float32, 128)
	query[0] = 0.5
	results, err := col.Search(query, 5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error searching: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Search results:")
	for _, r := range results {
		fmt.Printf("  ID=%d Score=%.4f\n", r.ID, r.Score)
	}
}
