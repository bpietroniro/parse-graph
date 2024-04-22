package utils

import (
	_ "container/heap"
	_ "math"
	"parse-graph/models"
)

func GetAdjacencyList(graph *models.Graph) map[string][]float64 {
	adjacencies := make(map[string][]float64)

	for _, node := range graph.Nodes {
		adjacencies[node.ID] = []float64{}
	}

	for _, edge := range graph.Edges {
		adjacencies[edge.FromID] = append(adjacencies[edge.FromID], edge.Cost)
	}

	return adjacencies
}

// TODO (Dijkstra)
// func CheapestPath(graphID, startID, endID string) {
// adjacencies := GetAdjacencyList(graph)

// var visited []*data.Node

// distances := make(map[string]float64)
// for _, node := range graph.Nodes {
// 	distances[node.ID] = math.Inf(1)
// }
// distances[startID] = 0

// for len(visited) < len(graph.Nodes) {
// 	minDistance := math.Inf(1)
// 	var currentNode *data.Node

// 	for _, node := range graph.Nodes {
// 		if !slices.Contains(visited, node) && distances[node.ID] < minDistance {
// 			minDistance = distances[node.ID]
// 			currentNode = node
// 		}
// 	}

// 	if currentNode == nil {
// 		break
// 	}

// 	visited = append(visited, currentNode)

// 	for _, edge := range graph.Edges {
// 		if edge.FromID == currentNode.ID {
// 			newDistance := distances[currentNode.ID] + edge.Cost
// 			if newDistance < distances[edge.ToID] {
// 				distances[edge.ToID] = newDistance
// 			}
// 		}
// 	}
// }
// }

func Cheapest(results []models.PathResult) models.PathResult {
	min := results[0]

	for _, result := range results {
		if result.Cost < min.Cost {
			min = result
		}
	}

	return min
}
