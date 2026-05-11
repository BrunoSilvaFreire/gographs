package pkg

import (
	"errors"
	"iter"
)

var (
	ErrVertexNotFound = errors.New("vertex not found")
	ErrEdgeNotFound   = errors.New("edge not found")
	ErrCycleDetected  = errors.New("cycle detected")
)

type GraphIndex = uint32

type Graph[V any, E any] interface {
	// Size returns the number of vertices in the graph.
	Size() int

	// IsEmpty returns true if the graph contains no vertices.
	IsEmpty() bool

	// Vertex returns a pointer to the vertex at the given index, or ErrVertexNotFound.
	Vertex(index GraphIndex) (*V, error)

	// Edge returns a pointer to the edge connecting from and to, or ErrEdgeNotFound.
	Edge(from, to GraphIndex) (*E, error)

	// Connect creates an edge between from and to with the given data.
	Connect(from, to GraphIndex, edge E) error

	// Disconnect removes the edge between from and to.
	Disconnect(from, to GraphIndex) error

	// EdgesFrom returns an iterator over all edges originating from the given index.
	EdgesFrom(index GraphIndex) iter.Seq2[GraphIndex, *E]

	// AllVertices returns an iterator over all valid vertex indices and their data.
	AllVertices() iter.Seq2[GraphIndex, *V]
}
