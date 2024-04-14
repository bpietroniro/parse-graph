package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"parse-graph/data"
	"parse-graph/utils"
)

func main() {
	// TODO need more test cases
	filePath := "examples/graph_1.xml"

	var graph data.Graph
	err := parseXML(filePath, &graph)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("XML parsing successful!")

	utils.FindAllPaths(&graph)
}

func parseXML(filePath string, graph *data.Graph) error {
	// Open the XML file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %v", err)
		return err
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Failed to read file: %v", err)
		return err
	}

	// Unmarshal the XML data into the Graph struct
	err = xml.Unmarshal(content, &graph)
	if err != nil {
		fmt.Printf("Failed to parse XML: %v", err)
		return err
	}

	return nil
}
