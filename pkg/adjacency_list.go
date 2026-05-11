package pkg

import (
	"iter"
)

type AdjacencyNode[V any, E any] struct {
	vertex V
	edges  map[GraphIndex]*E
	active bool
}

type AdjacencyList[V any, E any] struct {
	nodes       []AdjacencyNode[V, E]
	freeIndices []GraphIndex
	activeCount int
}

func NewAdjacencyList[V any, E any]() *AdjacencyList[V, E] {
	return &AdjacencyList[V, E]{
		nodes:       make([]AdjacencyNode[V, E], 0),
		freeIndices: make([]GraphIndex, 0),
	}
}

func (g *AdjacencyList[V, E]) Size() int {
	return g.activeCount
}

func (g *AdjacencyList[V, E]) IsEmpty() bool {
	return g.activeCount == 0
}

func (g *AdjacencyList[V, E]) AddVertex(vertex V) GraphIndex {
	var index GraphIndex
	if len(g.freeIndices) > 0 {
		index = g.freeIndices[len(g.freeIndices)-1]
		g.freeIndices = g.freeIndices[:len(g.freeIndices)-1]
		g.nodes[index].vertex = vertex
		g.nodes[index].edges = make(map[GraphIndex]*E)
		g.nodes[index].active = true
	} else {
		index = GraphIndex(len(g.nodes))
		g.nodes = append(g.nodes, AdjacencyNode[V, E]{
			vertex: vertex,
			edges:  make(map[GraphIndex]*E),
			active: true,
		})
	}
	g.activeCount++
	return index
}

func (g *AdjacencyList[V, E]) RemoveVertex(index GraphIndex) error {
	if index >= GraphIndex(len(g.nodes)) || !g.nodes[index].active {
		return ErrVertexNotFound
	}

	g.nodes[index].active = false
	g.nodes[index].edges = nil // Clear edges
	var zeroV V
	g.nodes[index].vertex = zeroV // Clear vertex data
	g.freeIndices = append(g.freeIndices, index)
	g.activeCount--
	return nil
}

func (g *AdjacencyList[V, E]) Vertex(index GraphIndex) (*V, error) {
	if index >= GraphIndex(len(g.nodes)) || !g.nodes[index].active {
		return nil, ErrVertexNotFound
	}
	return &g.nodes[index].vertex, nil
}

func (g *AdjacencyList[V, E]) Edge(from, to GraphIndex) (*E, error) {
	if from >= GraphIndex(len(g.nodes)) || !g.nodes[from].active {
		return nil, ErrVertexNotFound
	}
	if to >= GraphIndex(len(g.nodes)) || !g.nodes[to].active {
		return nil, ErrVertexNotFound
	}
	edge, ok := g.nodes[from].edges[to]
	if !ok {
		return nil, ErrEdgeNotFound
	}
	return edge, nil
}

func (g *AdjacencyList[V, E]) Connect(from, to GraphIndex, edge E) error {
	if from >= GraphIndex(len(g.nodes)) || !g.nodes[from].active {
		return ErrVertexNotFound
	}
	if to >= GraphIndex(len(g.nodes)) || !g.nodes[to].active {
		return ErrVertexNotFound
	}
	g.nodes[from].edges[to] = &edge
	return nil
}

func (g *AdjacencyList[V, E]) Disconnect(from, to GraphIndex) error {
	if from >= GraphIndex(len(g.nodes)) || !g.nodes[from].active {
		return ErrVertexNotFound
	}
	delete(g.nodes[from].edges, to)
	return nil
}

func (g *AdjacencyList[V, E]) EdgesFrom(index GraphIndex) iter.Seq2[GraphIndex, *E] {
	return func(yield func(GraphIndex, *E) bool) {
		if index >= GraphIndex(len(g.nodes)) || !g.nodes[index].active {
			return
		}
		for to, edge := range g.nodes[index].edges {
			if !yield(to, edge) {
				return
			}
		}
	}
}

func (g *AdjacencyList[V, E]) AllVertices() iter.Seq2[GraphIndex, *V] {
	return func(yield func(GraphIndex, *V) bool) {
		for i := range g.nodes {
			if g.nodes[i].active {
				if !yield(GraphIndex(i), &g.nodes[i].vertex) {
					return
				}
			}
		}
	}
}
