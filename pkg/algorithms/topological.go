package algorithms

import (
	"github.com/BrunoSilvaFreire/gographs/pkg"
)

// TopologicalSort performs a topological sort on the graph using Kahn's algorithm.
// It returns a slice of vertex indices in topological order, or ErrCycleDetected if a cycle is found.
func TopologicalSort[V any, E any](g pkg.Graph[V, E]) ([]pkg.GraphIndex, error) {
	inDegree := make(map[pkg.GraphIndex]int)
	allVertices := make([]pkg.GraphIndex, 0, g.Size())

	// Initialize in-degrees
	for idx, _ := range g.AllVertices() {
		allVertices = append(allVertices, idx)
		if _, ok := inDegree[idx]; !ok {
			inDegree[idx] = 0
		}
		for neighbor, _ := range g.EdgesFrom(idx) {
			inDegree[neighbor]++
		}
	}

	// Queue for vertices with 0 in-degree
	var queue []pkg.GraphIndex
	for _, idx := range allVertices {
		if inDegree[idx] == 0 {
			queue = append(queue, idx)
		}
	}

	var result []pkg.GraphIndex
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		result = append(result, u)

		for v, _ := range g.EdgesFrom(u) {
			inDegree[v]--
			if inDegree[v] == 0 {
				queue = append(queue, v)
			}
		}
	}

	if len(result) != g.Size() {
		return nil, pkg.ErrCycleDetected
	}

	return result, nil
}
