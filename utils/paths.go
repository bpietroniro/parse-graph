package utils

import (
	"container/heap"
	"math"
	"parse-graph/models"
)

type vertex struct {
	Node string
	Cost float64
}

type PriorityQueue []*vertex

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Cost < pq[j].Cost
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(v any) {
	*pq = append(*pq, v.(*vertex))
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	vertex := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return vertex
}

func reconstructPath(cameFrom map[string]string, current string) []string {
	totalPath := []string{current}
	for current != "" {
		current = cameFrom[current]
		totalPath = append([]string{current}, totalPath...)
	}
	return totalPath
}

func Dijkstra(graph *models.Graph, startID, goalID string) []string {
	adjList := getAdjacencyList(graph)

	start := startID
	goal := goalID

	visited := make(PriorityQueue, 0)
	heap.Init(&visited)
	heap.Push(&visited, start)

	cameFrom := make(map[string]string)

	minCosts := make(map[string]float64)
	for nodeID := range adjList {
		minCosts[nodeID] = math.Inf(1)
	}
	minCosts[start] = 0

	for visited.Len() > 0 {
		current := heap.Pop(&visited).(*vertex)

		if current.Node == goal {
			return reconstructPath(cameFrom, current.Node)
		}

		for _, vertex := range adjList[current.Node] {
			tentativeMinCost := minCosts[current.Node] + vertex.Cost
			if tentativeMinCost < minCosts[vertex.Node] {
				cameFrom[vertex.Node] = current.Node
				minCosts[vertex.Node] = tentativeMinCost
				if !seen(&visited, vertex.Node) {
					heap.Push(&visited, vertex)
				}
			}
		}
	}

	return nil
}

func seen(visited *PriorityQueue, nodeID string) bool {
	for _, n := range *visited {
		if n.Node == nodeID {
			return true
		}
	}
	return false
}

func getAdjacencyList(graph *models.Graph) map[string][]*vertex {
	adjacencies := make(map[string][]*vertex)

	for _, node := range graph.Nodes {
		adjacencies[node.ID] = []*vertex{}
	}

	for _, edge := range graph.Edges {
		v := vertex{Node: edge.ToID, Cost: edge.Cost}
		adjacencies[edge.FromID] = append(adjacencies[edge.FromID], &v)
	}

	return adjacencies
}
