package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"parse-graph/data"
	"parse-graph/utils"
	"path/filepath"
	"strconv"

	"github.com/beevik/etree"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file path as a command line argument.")
	}

	filename := os.Args[1]

	dbpool, err := data.ConnectToDB()
	if err != nil {
		panic("Could not connect to PostgreSQL")
	}
	defer dbpool.Close()

	ext := filepath.Ext(filename)

	switch ext {
	case ".xml":
		graph, err := parseAndValidateXML(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("XML parsing successful!")

		err = graph.SaveGraph(dbpool)
		if err != nil {
			fmt.Println(err)
			return
		}

		cycles, err := data.FindCycles(dbpool, graph.ID)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(cycles)

	case ".json":
		if len(os.Args) < 3 {
			fmt.Println("To query a graph, provide the graph ID as an additional argument after the JSON filename.")
			return
		}
		graphID := os.Args[2]

		// may wind up being unnecessary
		graph, err := data.LoadGraph(dbpool, graphID)
		if err != nil {
			fmt.Println(err)
			return
		}

		queries, err := parseJSON(filename)
		if err != nil {
			fmt.Println("failed to parse JSON file")
			fmt.Println(err)
			return
		}
		fmt.Println(queries)

		utils.FindAllPaths(graph)
	}

}

// TODO
func parseJSON(filePath string) (*data.QueryList, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ql data.QueryList
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ql)
	if err != nil {
		return nil, err
	}

	return &ql, nil
}

func parseAndValidateXML(filePath string) (*data.Graph, error) {
	document := etree.NewDocument()
	if err := document.ReadFromFile(filePath); err != nil {
		return nil, err
	}

	graphElement := document.Root()
	if graphElement == nil {
		return nil, errors.New("empty document")
	}
	if graphElement.Tag != "graph" {
		return nil, errors.New("invalid input: document root element is not a graph")
	}

	var graph data.Graph

	err := validateGraphTags(graphElement, &graph)
	if err != nil {
		return nil, err
	}

	idSet, err := validateNodesAndGetIDs(graphElement, &graph)
	if err != nil {
		return nil, err
	}

	err = validateEdges(graphElement, &graph, &idSet)
	if err != nil {
		return nil, err
	}

	return &graph, nil
}

func validateGraphTags(graphElement *etree.Element, graph *data.Graph) error {
	graphID, err := validateUniqueChild(graphElement, "id")
	if err != nil {
		return err
	}
	graph.ID = graphID.Text()

	graphName, err := validateUniqueChild(graphElement, "name")
	if err != nil {
		return err
	}
	graph.Name = graphName.Text()

	return nil
}

func validateNodesAndGetIDs(graphElement *etree.Element, graph *data.Graph) (map[string]struct{}, error) {
	nodeListElement, err := validateUniqueChild(graphElement, "nodes")
	if err != nil {
		return nil, err
	}

	idSet := make(map[string]struct{})
	for _, node := range nodeListElement.SelectElements("node") {
		nodeID, err := validateUniqueChild(node, "id")
		if err != nil {
			return nil, err
		}
		if _, exists := idSet[nodeID.Text()]; exists {
			return nil, errors.New("invalid input: found duplicate node ID")
		}
		idSet[nodeID.Text()] = struct{}{}

		nodeName, err := validateUniqueChild(node, "name")
		if err != nil {
			return nil, err
		}

		graph.Nodes = append(graph.Nodes, &data.Node{
			ID:   nodeID.Text(),
			Name: nodeName.Text(),
		})
	}

	if len(graph.Nodes) == 0 {
		return nil, errors.New("invalid input: nodes group is empty")
	}

	return idSet, nil
}

func validateEdges(graphElement *etree.Element, graph *data.Graph, idSet *map[string]struct{}) error {
	edgeListElement, err := validateUniqueChild(graphElement, "edges")
	if err != nil {
		return err
	}

	for _, edge := range edgeListElement.SelectElements("node") {
		edgeID, err := validateUniqueChild(edge, "id")
		if err != nil {
			return err
		}

		fromID, err := validateUniqueChild(edge, "from")
		if err != nil {
			return err
		}
		if _, exists := (*idSet)[fromID.Text()]; !exists {
			return fmt.Errorf("invalid input: edge %s's start node doesn't exist in the graph: %s", edgeID.Text(), fromID.Text())
		}

		toID, err := validateUniqueChild(edge, "to")
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

		graph.Edges = append(graph.Edges, &data.Edge{
			ID:     edgeID.Text(),
			FromID: fromID.Text(),
			ToID:   toID.Text(),
			Cost:   cost,
		})
	}

	return nil
}

func validateUniqueChild(e *etree.Element, tag string) (*etree.Element, error) {
	elements := e.SelectElements(tag)

	if elements == nil {
		return nil, fmt.Errorf("invalid input: missing %s tag in %s", tag, e.Tag)
	}

	if len(elements) > 1 {
		return nil, fmt.Errorf("invalid input: found duplicate %s tag in %s", tag, e.Tag)
	}

	return elements[0], nil
}
