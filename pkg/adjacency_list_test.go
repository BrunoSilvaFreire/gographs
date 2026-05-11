package pkg_test

import (
	"github.com/BrunoSilvaFreire/gographs/pkg"
	"testing"
	"errors"
)

func TestAdjacencyList_VertexRecycling(t *testing.T) {
	g := pkg.NewAdjacencyList[string, int]()

	// Add 3 vertices
	_ = g.AddVertex("V0")
	v1 := g.AddVertex("V1")
	_ = g.AddVertex("V2")

	if g.Size() != 3 {
		t.Errorf("Expected size 3, got %d", g.Size())
	}

	// Remove middle vertex
	err := g.RemoveVertex(v1)
	if err != nil {
		t.Fatalf("Failed to remove vertex: %v", err)
	}

	if g.Size() != 2 {
		t.Errorf("Expected size 2 after removal, got %d", g.Size())
	}

	// Verify v1 is not accessible
	_, err = g.Vertex(v1)
	if !errors.Is(err, pkg.ErrVertexNotFound) {
		t.Errorf("Expected ErrVertexNotFound for removed vertex, got %v", err)
	}

	// Add new vertex, it should recycle index v1
	v3 := g.AddVertex("V3")
	if v3 != v1 {
		t.Errorf("Expected recycled index %d, got %d", v1, v3)
	}

	if g.Size() != 3 {
		t.Errorf("Expected size 3 after recycling, got %d", g.Size())
	}
}

func TestAdjacencyList_Errors(t *testing.T) {
	g := pkg.NewAdjacencyList[string, int]()
	v0 := g.AddVertex("V0")

	// Invalid index access
	_, err := g.Vertex(999)
	if !errors.Is(err, pkg.ErrVertexNotFound) {
		t.Errorf("Expected ErrVertexNotFound, got %v", err)
	}

	// Connect non-existent
	err = g.Connect(v0, 999, 10)
	if !errors.Is(err, pkg.ErrVertexNotFound) {
		t.Errorf("Expected ErrVertexNotFound for destination, got %v", err)
	}

	// Edge not found
	v1 := g.AddVertex("V1")
	_, err = g.Edge(v0, v1)
	if !errors.Is(err, pkg.ErrEdgeNotFound) {
		t.Errorf("Expected ErrEdgeNotFound, got %v", err)
	}
}

func TestAdjacencyList_PointerStability(t *testing.T) {
	g := pkg.NewAdjacencyList[string, int]()
	v0 := g.AddVertex("V0")
	v1 := g.AddVertex("V1")
	
	g.Connect(v0, v1, 100)
	edge, _ := g.Edge(v0, v1)
	*edge = 200

	// Force slice growth in AdjacencyNode and potentially AdjacencyList
	for i := 0; i < 1000; i++ {
		g.AddVertex("Growth")
	}

	// Map pointers in Go are stable until the map grows, but since we store *E in the map,
	// the pointer to E itself is always stable regardless of map growth.
	edgeAgain, _ := g.Edge(v0, v1)
	if *edgeAgain != 200 {
		t.Errorf("Mutation lost after growth: expected 200, got %d", *edgeAgain)
	}
	if edge != edgeAgain {
		t.Errorf("Edge pointer changed. Stored pointers should be stable.")
	}
}
