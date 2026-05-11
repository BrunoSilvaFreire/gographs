package algorithms

import (
	"github.com/BrunoSilvaFreire/gographs/pkg"
	"iter"
	"slices"
)

// BFS returns an iterator that performs a Breadth-First Search starting from the given indices.
func BFS[V any, E any](g pkg.Graph[V, E], starts ...pkg.GraphIndex) iter.Seq[pkg.GraphIndex] {
	return func(yield func(pkg.GraphIndex) bool) {
		visited := make(map[pkg.GraphIndex]struct{})
		queue := make([]pkg.GraphIndex, 0, len(starts))

		for _, start := range starts {
			if _, ok := visited[start]; !ok {
				visited[start] = struct{}{}
				queue = append(queue, start)
			}
		}

		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]

			if !yield(current) {
				return
			}

			for neighbor, _ := range g.EdgesFrom(current) {
				if _, ok := visited[neighbor]; !ok {
					visited[neighbor] = struct{}{}
					queue = append(queue, neighbor)
				}
			}
		}
	}
}

// DFS returns an iterator that performs a Depth-First Search starting from the given indices.
func DFS[V any, E any](g pkg.Graph[V, E], starts ...pkg.GraphIndex) iter.Seq[pkg.GraphIndex] {
	return func(yield func(pkg.GraphIndex) bool) {
		visited := make(map[pkg.GraphIndex]struct{})
		stack := make([]pkg.GraphIndex, 0, len(starts))

		// To maintain order, we push starts in reverse if we want the first start to be popped first,
		// but since starts is variadic, let's just push them.
		for i := len(starts) - 1; i >= 0; i-- {
			stack = append(stack, starts[i])
		}

		for len(stack) > 0 {
			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if _, ok := visited[current]; ok {
				continue
			}

			if !yield(current) {
				return
			}

			visited[current] = struct{}{}

			for neighbor, _ := range g.EdgesFrom(current) {
				if _, ok := visited[neighbor]; !ok {
					stack = append(stack, neighbor)
				}
			}
		}
	}
}

// Flood is an alias for BFS with multiple starting points.
func Flood[V any, E any](g pkg.Graph[V, E], starts ...pkg.GraphIndex) iter.Seq[pkg.GraphIndex] {
	return BFS(g, starts...)
}

// ReverseLevelOrder returns an iterator that performs a Reverse Level Order traversal
// (BFS order reversed) starting from the given indices.
func ReverseLevelOrder[V any, E any](g pkg.Graph[V, E], starts ...pkg.GraphIndex) iter.Seq[pkg.GraphIndex] {
	return func(yield func(pkg.GraphIndex) bool) {
		var order []pkg.GraphIndex
		for vertex := range BFS(g, starts...) {
			order = append(order, vertex)
		}

		slices.Reverse(order)
		for _, vertex := range order {
			if !yield(vertex) {
				return
			}
		}
	}
}
