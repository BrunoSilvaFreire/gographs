package pkg_test

import (
	"errors"
	"github.com/BrunoSilvaFreire/gographs/pkg"
	"github.com/BrunoSilvaFreire/gographs/pkg/algorithms"
	"reflect"
	"testing"
)

func TestAStar_WeightedPaths(t *testing.T) {
	g := pkg.NewAdjacencyList[string, float64]()
	s := g.AddVertex("Start")
	a := g.AddVertex("A")
	b := g.AddVertex("B")
	e := g.AddVertex("End")

	// Two paths:
	// 1. S --(10)--> A --(10)--> E (Total 20, 3 nodes)
	// 2. S --(1)--> B --(1)--> E  (Total 2, 3 nodes)
	g.Connect(s, a, 10.0)
	g.Connect(a, e, 10.0)
	g.Connect(s, b, 1.0)
	g.Connect(b, e, 1.0)

	heuristic := func(from, to pkg.GraphIndex) float64 { return 0 }
	distance := func(from, to pkg.GraphIndex, edge *float64) float64 { return *edge }

	path, err := algorithms.AStar(g, s, e, heuristic, distance)
	if err != nil {
		t.Fatalf("AStar failed: %v", err)
	}

	expected := []pkg.GraphIndex{s, b, e}
	if !reflect.DeepEqual(path, expected) {
		t.Errorf("AStar failed to pick cheaper path. Expected %v, got %v", expected, path)
	}
}

func TestAStar_LongerButCheaper(t *testing.T) {
	g := pkg.NewAdjacencyList[string, float64]()
	s := g.AddVertex("Start")
	a := g.AddVertex("A")
	b1 := g.AddVertex("B1")
	b2 := g.AddVertex("B2")
	e := g.AddVertex("End")

	// Path 1: Shortest in nodes, heavy weight
	// S --(10)--> A --(10)--> E (Total 20)
	g.Connect(s, a, 10.0)
	g.Connect(a, e, 10.0)

	// Path 2: Longer in nodes, light weight
	// S --(1)--> B1 --(1)--> B2 --(1)--> E (Total 3)
	g.Connect(s, b1, 1.0)
	g.Connect(b1, b2, 1.0)
	g.Connect(b2, e, 1.0)

	heuristic := func(from, to pkg.GraphIndex) float64 { return 0 }
	distance := func(from, to pkg.GraphIndex, edge *float64) float64 { return *edge }

	path, err := algorithms.AStar(g, s, e, heuristic, distance)
	if err != nil {
		t.Fatalf("AStar failed: %v", err)
	}

	expected := []pkg.GraphIndex{s, b1, b2, e}
	if !reflect.DeepEqual(path, expected) {
		t.Errorf("AStar failed to pick longer-but-cheaper path. Expected %v, got %v", expected, path)
	}
}

func TestAStar_Unreachable(t *testing.T) {
	g := pkg.NewAdjacencyList[string, float64]()
	s := g.AddVertex("Start")
	e := g.AddVertex("End")

	heuristic := func(from, to pkg.GraphIndex) float64 { return 0 }
	distance := func(from, to pkg.GraphIndex, edge *float64) float64 { return *edge }

	_, err := algorithms.AStar(g, s, e, heuristic, distance)
	if !errors.Is(err, algorithms.ErrPathNotFound) {
		t.Errorf("Expected ErrPathNotFound for disconnected nodes, got %v", err)
	}
}

func TestAStar_SameStartEnd(t *testing.T) {
	g := pkg.NewAdjacencyList[string, float64]()
	s := g.AddVertex("Start")

	heuristic := func(from, to pkg.GraphIndex) float64 { return 0 }
	distance := func(from, to pkg.GraphIndex, edge *float64) float64 { return *edge }

	path, err := algorithms.AStar(g, s, s, heuristic, distance)
	if err != nil {
		t.Fatalf("AStar failed on same start/end: %v", err)
	}

	expected := []pkg.GraphIndex{s}
	if !reflect.DeepEqual(path, expected) {
		t.Errorf("Expected path [s], got %v", path)
	}
}
