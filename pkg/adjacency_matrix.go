package pkg

import (
	"iter"
)

type AdjacencyMatrix[V any, E any] struct {
	vertices []V
	edges    []E
	present  []bool
	size     int
}

func NewAdjacencyMatrix[V any, E any](size int) *AdjacencyMatrix[V, E] {
	return &AdjacencyMatrix[V, E]{
		vertices: make([]V, size),
		edges:    make([]E, size*size),
		present:  make([]bool, size*size),
		size:     size,
	}
}

func (g *AdjacencyMatrix[V, E]) Size() int {
	return g.size
}

func (g *AdjacencyMatrix[V, E]) IsEmpty() bool {
	return g.size == 0
}

func (g *AdjacencyMatrix[V, E]) Vertex(index GraphIndex) (*V, error) {
	if int(index) >= g.size {
		return nil, ErrVertexNotFound
	}
	return &g.vertices[index], nil
}

func (g *AdjacencyMatrix[V, E]) Edge(from, to GraphIndex) (*E, error) {
	if int(from) >= g.size || int(to) >= g.size {
		return nil, ErrVertexNotFound
	}
	idx := int(from)*g.size + int(to)
	if !g.present[idx] {
		return nil, ErrEdgeNotFound
	}
	return &g.edges[idx], nil
}

func (g *AdjacencyMatrix[V, E]) Connect(from, to GraphIndex, edge E) error {
	if int(from) >= g.size || int(to) >= g.size {
		return ErrVertexNotFound
	}
	idx := int(from)*g.size + int(to)
	g.edges[idx] = edge
	g.present[idx] = true
	return nil
}

func (g *AdjacencyMatrix[V, E]) Disconnect(from, to GraphIndex) error {
	if int(from) >= g.size || int(to) >= g.size {
		return ErrVertexNotFound
	}
	idx := int(from)*g.size + int(to)
	g.present[idx] = false
	var zeroE E
	g.edges[idx] = zeroE
	return nil
}

func (g *AdjacencyMatrix[V, E]) EdgesFrom(index GraphIndex) iter.Seq2[GraphIndex, *E] {
	return func(yield func(GraphIndex, *E) bool) {
		if int(index) >= g.size {
			return
		}
		start := int(index) * g.size
		for i := 0; i < g.size; i++ {
			idx := start + i
			if g.present[idx] {
				if !yield(GraphIndex(i), &g.edges[idx]) {
					return
				}
			}
		}
	}
}

func (g *AdjacencyMatrix[V, E]) AllVertices() iter.Seq2[GraphIndex, *V] {
	return func(yield func(GraphIndex, *V) bool) {
		for i := 0; i < g.size; i++ {
			if !yield(GraphIndex(i), &g.vertices[i]) {
				return
			}
		}
	}
}
