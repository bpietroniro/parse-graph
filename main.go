package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"parse-graph/data"
)

func main() {
	filePath := "examples/graph_1.xml"

	// Open the XML file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %v", err)
		return
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Failed to read file: %v", err)
		return
	}

	// Unmarshal the XML data into the struct
	var graph data.Graph
	err = xml.Unmarshal(content, &graph)
	if err != nil {
		fmt.Printf("Failed to parse XML: %v", err)
		return
	}

	fmt.Println("XML parsing successful!")
	fmt.Println(graph)
	fmt.Println(graph.Nodes)
	fmt.Println(graph.Edges)
}
