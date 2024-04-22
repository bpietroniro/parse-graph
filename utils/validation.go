package utils

import (
	"errors"
	"fmt"
	"parse-graph/models"
	"strconv"

	"github.com/beevik/etree"
)

func ValidateGraphTags(graphElement *etree.Element, graph *models.Graph) error {
	graphID, err := ValidateUniqueChild(graphElement, "id")
	if err != nil {
		return err
	}
	graph.ID = graphID.Text()

	graphName, err := ValidateUniqueChild(graphElement, "name")
	if err != nil {
		return err
	}
	graph.Name = graphName.Text()

	return nil
}

func ValidateNodesAndGetIDs(graphElement *etree.Element, graph *models.Graph) (map[string]struct{}, error) {
	nodeListElement, err := ValidateUniqueChild(graphElement, "nodes")
	if err != nil {
		return nil, err
	}

	idSet := make(map[string]struct{})
	for _, node := range nodeListElement.SelectElements("node") {
		nodeID, err := ValidateUniqueChild(node, "id")
		if err != nil {
			return nil, err
		}
		if _, exists := idSet[nodeID.Text()]; exists {
			return nil, errors.New("invalid input: found duplicate node ID")
		}
		idSet[nodeID.Text()] = struct{}{}

		nodeName, err := ValidateUniqueChild(node, "name")
		if err != nil {
			return nil, err
		}

		graph.Nodes = append(graph.Nodes, &models.Node{
			ID:   nodeID.Text(),
			Name: nodeName.Text(),
		})
	}

	if len(graph.Nodes) == 0 {
		return nil, errors.New("invalid input: nodes group is empty")
	}

	return idSet, nil
}

func ValidateEdges(graphElement *etree.Element, graph *models.Graph, idSet *map[string]struct{}) error {
	edgeListElement, err := ValidateUniqueChild(graphElement, "edges")
	if err != nil {
		return err
	}

	for _, edge := range edgeListElement.SelectElements("node") {
		edgeID, err := ValidateUniqueChild(edge, "id")
		if err != nil {
			return err
		}

		fromID, err := ValidateUniqueChild(edge, "from")
		if err != nil {
			return err
		}
		if _, exists := (*idSet)[fromID.Text()]; !exists {
			return fmt.Errorf("invalid input: edge %s's start node doesn't exist in the graph: %s", edgeID.Text(), fromID.Text())
		}

		toID, err := ValidateUniqueChild(edge, "to")
		if err != nil {
			return err
		}
		if _, exists := (*idSet)[toID.Text()]; !exists {
			return fmt.Errorf("invalid input: edge %s's end node (%s) doesn't exist in the graph", edgeID.Text(), toID.Text())
		}

		costStr := edge.SelectElement("cost").Text()
		var cost float64
		if costStr == "" {
			cost = 0
		} else {
			cost, err = strconv.ParseFloat(costStr, 64)
			if err != nil {
				return err
			}
		}

		graph.Edges = append(graph.Edges, &models.Edge{
			ID:     edgeID.Text(),
			FromID: fromID.Text(),
			ToID:   toID.Text(),
			Cost:   cost,
		})
	}

	return nil
}

func ValidateUniqueChild(e *etree.Element, tag string) (*etree.Element, error) {
	elements := e.SelectElements(tag)

	if elements == nil {
		return nil, fmt.Errorf("invalid input: missing %s tag in %s", tag, e.Tag)
	}

	if len(elements) > 1 {
		return nil, fmt.Errorf("invalid input: found duplicate %s tag in %s", tag, e.Tag)
	}

	return elements[0], nil
}
