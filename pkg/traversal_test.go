package pkg_test

import (
	"github.com/BrunoSilvaFreire/gographs/pkg"
	"github.com/BrunoSilvaFreire/gographs/pkg/algorithms"
	"testing"
)

func TestTraversal_Cycles(t *testing.T) {
	g := pkg.NewAdjacencyList[string, int]()
	a := g.AddVertex("A")
	b := g.AddVertex("B")
	c := g.AddVertex("C")

	// Cycle: A -> B -> C -> A
	g.Connect(a, b, 1)
	g.Connect(b, c, 1)
	g.Connect(c, a, 1)

	visited := make(map[pkg.GraphIndex]int)
	for v := range algorithms.BFS(g, a) {
		visited[v]++
	}

	if len(visited) != 3 {
		t.Errorf("BFS should visit 3 nodes in cycle, visited %d", len(visited))
	}
	for v, count := range visited {
		if count > 1 {
			t.Errorf("Node %d visited %d times, expected 1", v, count)
		}
	}
}

func TestTraversal_Disconnected(t *testing.T) {
	g := pkg.NewAdjacencyList[string, int]()
	a := g.AddVertex("A")
	b := g.AddVertex("B")
	c := g.AddVertex("C")
	d := g.AddVertex("D")

	// Cluster 1: A -> B
	g.Connect(a, b, 1)
	// Cluster 2: C -> D
	g.Connect(c, d, 1)

	visited := make(map[pkg.GraphIndex]bool)
	for v := range algorithms.BFS(g, a) {
		visited[v] = true
	}

	if len(visited) != 2 {
		t.Errorf("BFS from A should only visit cluster 1 (A, B), visited %d nodes", len(visited))
	}
	if visited[c] || visited[d] {
		t.Errorf("BFS from A should not reach cluster 2")
	}
}

func TestTraversal_EmptyAndMultiStart(t *testing.T) {
	g := pkg.NewAdjacencyList[string, int]()
	
	// Empty graph
	count := 0
	for v := range algorithms.BFS(g) {
		_ = v
		count++
	}
	if count != 0 {
		t.Errorf("BFS on empty graph should yield 0, got %d", count)
	}

	// Multi-start overlapping
	a := g.AddVertex("A")
	b := g.AddVertex("B")
	c := g.AddVertex("C")
	g.Connect(a, b, 1)
	g.Connect(b, c, 1)

	visited := make(map[pkg.GraphIndex]bool)
	// Start from both A and B. B is reachable from A.
	for v := range algorithms.Flood(g, a, b) {
		visited[v] = true
	}
	if len(visited) != 3 {
		t.Errorf("Flood from A and B should visit all 3 nodes exactly once, visited %d", len(visited))
	}
}
