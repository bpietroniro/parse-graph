package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"parse-graph/data"
	"parse-graph/models"
	"parse-graph/utils"
	"path/filepath"

	"github.com/beevik/etree"
)

func main() {
	var filename string

	// If no arguments are given, prompt for file input
	if len(os.Args) < 2 {
		for {
			fmt.Println("Welcome! To proceed, please enter the path of the file you'd like to process. This can be either:")
			fmt.Println("1. A file containing an XML representation of a graph,")
			fmt.Println("2. A file containing JSON queries to be performed on an already-saved graph.")
			fmt.Print("-> ")
			fmt.Scan(&filename)
			if filename != "" {
				break
			}
		}
	} else {
		filename = os.Args[1]
	}

	// Connect to the database
	dbpool, err := data.ConnectToDB()
	if err != nil {
		fmt.Printf("Could not connect to PostgreSQL: %v\n", err)
		return
	}
	defer dbpool.Close()

	ext := filepath.Ext(filename)

	switch ext {
	// Parse a graph from XML and save it to the database
	case ".xml":
		graph, err := parseAndValidateXML(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("XML parsing successful!")

		err = data.SaveGraph(graph)
		if err != nil {
			fmt.Println(err)
			return
		}

		for {
			fmt.Println(`Do you want to find the graph's cycles? (y)`)
			var ans string
			fmt.Scan(&ans)

			if ans == "y" || ans == "Y" {
				cycles, err := data.FindCycles(graph.ID)
				if err != nil {
					fmt.Println(err)
					return
				}
				if len(cycles) == 0 {
					fmt.Println("The graph contains no cycles.")
				} else {
					fmt.Println("The graph contains the following cyclic paths:")
					fmt.Println(cycles)
				}
				break
			} else if ans == "n" || ans == "N" {
				break
			}
		}

	// Parse a list of JSON queries and output a list of corresponding JSON results
	case ".json":
		var graphID string

		if len(os.Args) < 3 {
			for {
				fmt.Println("To execute queries on an existing graph, enter the graph ID: ")
				fmt.Scan(&graphID)
				if graphID != "" {
					break
				}
			}
		} else {
			graphID = os.Args[2]
		}

		// may wind up being unnecessary
		// graph, err := data.LoadGraph(dbpool, graphID)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

		queries, err := parseJSON(filename)
		if err != nil {
			fmt.Println("failed to parse JSON file")
			fmt.Println(err)
			return
		}

		queryMap := make(map[models.QueryInputs][]models.PathResult)

		for _, q := range queries.Queries {
			if q.Paths != nil {
				start := q.Paths.Start
				end := q.Paths.End
				qi := models.QueryInputs{Start: start, End: end}

				allPaths, ok := queryMap[qi]
				if !ok {
					allPaths, err = data.FindAllPaths(graphID, start, end)
					queryMap[models.QueryInputs{Start: start, End: end}] = allPaths
					if err != nil {
						fmt.Println(err)
						return
					}
				}

				// TODO marshal as JSON
				fmt.Print("All paths: ")
				fmt.Println(allPaths)
			} else if q.Cheapest != nil {
				start := q.Cheapest.Start
				end := q.Cheapest.End
				qi := models.QueryInputs{Start: start, End: end}

				allPaths, ok := queryMap[qi]
				if !ok {
					allPaths, err = data.FindAllPaths(graphID, start, end)
					queryMap[models.QueryInputs{Start: start, End: end}] = allPaths
					if err != nil {
						fmt.Println(err)
						return
					}
				}

				var cheapestPath models.PathResult
				if len(allPaths) > 0 {
					cheapestPath = utils.Cheapest(allPaths)
					// TODO marshal as JSON
					fmt.Print("Cheapest path: ")
					fmt.Println(cheapestPath)
				} else {
					fmt.Println("no paths found")
				}
			}
		}
	}

}

func parseJSON(filePath string) (*models.QueryList, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ql models.QueryList
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ql)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal(ql)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsonBytes))

	return &ql, nil
}

func parseAndValidateXML(filePath string) (*models.Graph, error) {
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

	var graph models.Graph

	err := utils.ValidateGraphTags(graphElement, &graph)
	if err != nil {
		return nil, err
	}

	idSet, err := utils.ValidateNodesAndGetIDs(graphElement, &graph)
	if err != nil {
		return nil, err
	}

	err = utils.ValidateEdges(graphElement, &graph, &idSet)
	if err != nil {
		return nil, err
	}

	return &graph, nil
}
