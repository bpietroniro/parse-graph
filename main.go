package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"parse-graph/data"
	"parse-graph/utils"
	"strconv"

	"github.com/beevik/etree"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	dbpool, err := ConnectToDB()
	if err != nil {
		panic("Could not connect to PostgreSQL")
	}
	defer dbpool.Close()

	// TODO need more test cases
	filePath := "examples/graph_1.xml"

	graph, err := parseXML(filePath)
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

	loaded_graph, err := data.LoadGraph(dbpool, "g0")
	if err != nil {
		fmt.Println(err)
		return
	}

	utils.FindAllPaths(loaded_graph)
}

func parseXML(filePath string) (*data.Graph, error) {
	document := etree.NewDocument()
	if err := document.ReadFromFile(filePath); err != nil {
		return nil, err
	}

	graphElement := document.Root()
	if graphElement == nil {
		return nil, errors.New("empty document")
	}
	if graphElement.Tag != "graph" {
		return nil, errors.New("document root element is not a graph")
	}

	var graph data.Graph

	graphID, err := validateUniqueChild(graphElement, "id")
	if err != nil {
		return nil, err
	}
	graph.ID = graphID.Text()

	graphName, err := validateUniqueChild(graphElement, "name")
	if err != nil {
		return nil, err
	}
	graph.Name = graphName.Text()

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
			return nil, errors.New("found duplicate node ID")
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
		return nil, errors.New("nodes group is empty")
	}

	edgeListElement, err := validateUniqueChild(graphElement, "edges")
	if err != nil {
		return nil, err
	}

	for _, edge := range edgeListElement.SelectElements("node") {
		edgeID, err := validateUniqueChild(edge, "id")
		if err != nil {
			return nil, err
		}

		fromID, err := validateUniqueChild(edge, "from")
		if err != nil {
			return nil, err
		}
		if _, exists := idSet[fromID.Text()]; !exists {
			return nil, fmt.Errorf("start node doesn't exist in the graph: %s", fromID.Text())
		}

		toID, err := validateUniqueChild(edge, "to")
		if err != nil {
			return nil, err
		}
		if _, exists := idSet[toID.Text()]; !exists {
			return nil, fmt.Errorf("end node doesn't exist in the graph: %s", toID.Text())
		}

		costStr := edge.SelectElement("cost").Text()
		var cost float64
		if costStr == "" {
			cost = 0
		} else {
			cost, err = strconv.ParseFloat(costStr, 64)
			if err != nil {
				return nil, err
			}
		}

		graph.Edges = append(graph.Edges, &data.Edge{
			ID:     edgeID.Text(),
			FromID: fromID.Text(),
			ToID:   toID.Text(),
			Cost:   cost,
		})
	}

	return &graph, nil
}

func validateUniqueChild(e *etree.Element, tag string) (*etree.Element, error) {
	elements := e.SelectElements(tag)

	if elements == nil {
		return nil, fmt.Errorf("missing %s in %s", tag, e.Tag)
	}

	if len(elements) > 1 {
		return nil, fmt.Errorf("found duplicate %s in %s", tag, e.Tag)
	}

	return elements[0], nil
}

func ConnectToDB() (*pgxpool.Pool, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return nil, err
	}

	connectionString := fmt.Sprintf("postgres://%s@%s:%s/%s", os.Getenv("DB_USER"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	dbpool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		return nil, err
	}

	return dbpool, nil
}

func OpenDB() {}
