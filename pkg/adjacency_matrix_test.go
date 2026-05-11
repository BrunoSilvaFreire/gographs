package pkg_test

import (
	"github.com/BrunoSilvaFreire/gographs/pkg"
	"testing"
	"errors"
)

func TestAdjacencyMatrix_Errors(t *testing.T) {
	g := pkg.NewAdjacencyMatrix[string, int](5)

	// Out of bounds
	_, err := g.Vertex(5)
	if !errors.Is(err, pkg.ErrVertexNotFound) {
		t.Errorf("Expected ErrVertexNotFound for index 5, got %v", err)
	}

	// Connect out of bounds
	err = g.Connect(0, 5, 10)
	if !errors.Is(err, pkg.ErrVertexNotFound) {
		t.Errorf("Expected ErrVertexNotFound for destination 5, got %v", err)
	}

	// Edge not found (but within bounds)
	_, err = g.Edge(0, 1)
	if !errors.Is(err, pkg.ErrEdgeNotFound) {
		t.Errorf("Expected ErrEdgeNotFound for unconnected indices, got %v", err)
	}
}

func TestAdjacencyMatrix_MutationAndGrowth(t *testing.T) {
	// AdjacencyMatrix is fixed size, but we can verify pointer stability within that size
	g := pkg.NewAdjacencyMatrix[string, int](10)
	v0, v1 := pkg.GraphIndex(0), pkg.GraphIndex(1)

	g.Connect(v0, v1, 50)
	edge, _ := g.Edge(v0, v1)
	*edge = 100

	// Disconnect and reconnect
	g.Disconnect(v0, v1)
	_, err := g.Edge(v0, v1)
	if !errors.Is(err, pkg.ErrEdgeNotFound) {
		t.Errorf("Edge should be gone")
	}

	g.Connect(v0, v1, 200)
	edgeAgain, _ := g.Edge(v0, v1)
	if *edgeAgain != 200 {
		t.Errorf("Expected 200, got %d", *edgeAgain)
	}
}

func TestAdjacencyMatrix_Iteration(t *testing.T) {
	g := pkg.NewAdjacencyMatrix[string, int](3)
	g.Connect(0, 1, 10)
	g.Connect(0, 2, 20)

	count := 0
	for to, edge := range g.EdgesFrom(0) {
		count++
		if to == 1 && *edge != 10 {
			t.Errorf("Wrong edge value for 1")
		}
		if to == 2 && *edge != 20 {
			t.Errorf("Wrong edge value for 2")
		}
	}
	if count != 2 {
		t.Errorf("Expected 2 edges, got %d", count)
	}
}
