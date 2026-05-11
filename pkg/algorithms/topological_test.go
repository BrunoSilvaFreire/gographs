package algorithms

import (
	"github.com/BrunoSilvaFreire/gographs/pkg"
	"testing"
)

func TestTopologicalSort(t *testing.T) {
	g := pkg.NewAdjacencyList[string, struct{}]()

	v1 := g.AddVertex("v1")
	v2 := g.AddVertex("v2")
	v3 := g.AddVertex("v3")

	// v1 -> v2 -> v3
	g.Connect(v1, v2, struct{}{})
	g.Connect(v2, v3, struct{}{})

	order, err := TopologicalSort(g)
	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}

	if len(order) != 3 {
		t.Fatalf("expected 3 vertices, got %d", len(order))
	}

	// Order should be v1, v2, v3
	if order[0] != v1 || order[1] != v2 || order[2] != v3 {
		t.Errorf("incorrect order: %v", order)
	}
}

func TestTopologicalSortCycle(t *testing.T) {
	g := pkg.NewAdjacencyList[string, struct{}]()

	v1 := g.AddVertex("v1")
	v2 := g.AddVertex("v2")
	v3 := g.AddVertex("v3")

	// v1 -> v2 -> v3 -> v1
	g.Connect(v1, v2, struct{}{})
	g.Connect(v2, v3, struct{}{})
	g.Connect(v3, v1, struct{}{})

	_, err := TopologicalSort(g)
	if err != pkg.ErrCycleDetected {
		t.Fatalf("expected ErrCycleDetected, got %v", err)
	}
}
