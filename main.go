package main

import (
	"encoding/json"
	"fmt"
	"os"
	"parse-graph/data"
	"parse-graph/models"
	"parse-graph/utils"
	"path/filepath"
)

func main() {
	var filePath string

	// If no arguments are given, prompt for file input
	if len(os.Args) < 2 {
		for {
			fmt.Println("Welcome! To proceed, please enter the path of the file you'd like to process. This can be either:")
			fmt.Println("1. A file containing an XML representation of a graph,")
			fmt.Println("2. A file containing JSON queries to be performed on an already-saved graph.")
			fmt.Print("-> ")
			fmt.Scan(&filePath)
			if filePath != "" {
				break
			}
		}
	} else {
		filePath = os.Args[1]
	}

	// Connect to the database
	dbpool, err := data.ConnectToDB()
	if err != nil {
		fmt.Printf("Could not connect to PostgreSQL: %v\n", err)
		return
	}
	defer dbpool.Close()

	ext := filepath.Ext(filePath)

	switch ext {
	case ".xml":
		handleXML(filePath)

	// Parse a list of JSON queries and output a list of corresponding JSON results
	case ".json":
		handleJSON(filePath)
	}

}

// Parse a graph from XML and save it to the database
func handleXML(filePath string) {
	graph, err := utils.ParseAndValidateXML(filePath)
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
}

// Parse a list of queries, fulfill them, and output the results as JSON
func handleJSON(filePath string) {
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

	queries, err := utils.ParseJSON(filePath)
	if err != nil {
		fmt.Printf("failed to parse JSON file: %v\n", err)
		return
	}

	queryCache := make(map[models.QueryInputs][]models.PathResult)
	answers := map[string][]interface{}{"answers": {}}

	for _, q := range queries.Queries {
		if q.Paths != nil {
			start := q.Paths.Start
			end := q.Paths.End
			qi := models.QueryInputs{Start: start, End: end}

			allPaths, ok := queryCache[qi]
			if !ok {
				allPaths, err = data.FindAllPaths(graphID, start, end)
				queryCache[models.QueryInputs{Start: start, End: end}] = allPaths
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			var ans models.PathResultList
			ans.From = start
			ans.To = end
			ans.Paths = [][]string{}
			for _, path := range allPaths {
				ans.Paths = append(ans.Paths, path.Path)
			}

			answers["answers"] = append(answers["answers"], map[string]interface{}{"paths": ans})

		} else if q.Cheapest != nil {
			start := q.Cheapest.Start
			end := q.Cheapest.End
			qi := models.QueryInputs{Start: start, End: end}

			allPaths, ok := queryCache[qi]
			if !ok {
				allPaths, err = data.FindAllPaths(graphID, start, end)
				queryCache[models.QueryInputs{Start: start, End: end}] = allPaths
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			var cheapestPath models.PathResult
			if len(allPaths) > 0 {
				cheapestPath = allPaths[0]

				// Find the past with minimum cost
				for _, p := range allPaths {
					if p.Cost < cheapestPath.Cost {
						cheapestPath = p
					}
				}
			}

			ans := models.CheapestResult{From: start, To: end, Path: false}
			if cheapestPath.Path != nil {
				ans.Path = cheapestPath.Path
			}

			answers["answers"] = append(answers["answers"], map[string]interface{}{"cheapest": ans})
		}

	}

	jsonBytes, err := json.MarshalIndent(answers, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %v", err)
		return
	}
	fmt.Println(string(jsonBytes))
}
