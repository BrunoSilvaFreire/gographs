package algorithms

import (
	"container/heap"
	"errors"
	"github.com/BrunoSilvaFreire/gographs/pkg"
	"math"
)

var ErrPathNotFound = errors.New("path not found")

type pqItem struct {
	index    pkg.GraphIndex
	priority float64
}

type priorityQueue []pqItem

func (pq priorityQueue) Len() int           { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].priority < pq[j].priority }
func (pq priorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x any)        { *pq = append(*pq, x.(pqItem)) }
func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// AStar performs an A* search to find the shortest path from 'from' to 'to'.
func AStar[V any, E any](
	g pkg.Graph[V, E],
	from, to pkg.GraphIndex,
	heuristic func(pkg.GraphIndex, pkg.GraphIndex) float64,
	distance func(pkg.GraphIndex, pkg.GraphIndex, *E) float64,
) ([]pkg.GraphIndex, error) {
	cameFrom := make(map[pkg.GraphIndex]pkg.GraphIndex)
	gScore := make(map[pkg.GraphIndex]float64)
	fScore := make(map[pkg.GraphIndex]float64)

	for i, _ := range g.AllVertices() {
		gScore[i] = math.Inf(1)
		fScore[i] = math.Inf(1)
	}

	gScore[from] = 0
	fScore[from] = heuristic(from, to)

	pq := &priorityQueue{}
	heap.Init(pq)
	heap.Push(pq, pqItem{index: from, priority: fScore[from]})

	openSet := make(map[pkg.GraphIndex]struct{})
	openSet[from] = struct{}{}

	for pq.Len() > 0 {
		current := heap.Pop(pq).(pqItem).index
		delete(openSet, current)

		if current == to {
			return reconstructPath(cameFrom, current), nil
		}

		for neighbor, edge := range g.EdgesFrom(current) {
			tentativeGScore := gScore[current] + distance(current, neighbor, edge)
			if tentativeGScore < gScore[neighbor] {
				cameFrom[neighbor] = current
				gScore[neighbor] = tentativeGScore
				fScore[neighbor] = tentativeGScore + heuristic(neighbor, to)
				if _, ok := openSet[neighbor]; !ok {
					heap.Push(pq, pqItem{index: neighbor, priority: fScore[neighbor]})
					openSet[neighbor] = struct{}{}
				}
			}
		}
	}

	return nil, ErrPathNotFound
}

func reconstructPath(cameFrom map[pkg.GraphIndex]pkg.GraphIndex, current pkg.GraphIndex) []pkg.GraphIndex {
	path := []pkg.GraphIndex{current}
	for {
		parent, ok := cameFrom[current]
		if !ok {
			break
		}
		current = parent
		path = append([]pkg.GraphIndex{current}, path...)
	}
	return path
}
