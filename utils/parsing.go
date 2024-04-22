package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"parse-graph/models"

	"github.com/beevik/etree"
)

func ParseAndValidateXML(filePath string) (*models.Graph, error) {
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

	err := ValidateGraphTags(graphElement, &graph)
	if err != nil {
		return nil, err
	}

	idSet, err := ValidateNodesAndGetIDs(graphElement, &graph)
	if err != nil {
		return nil, err
	}

	err = ValidateEdges(graphElement, &graph, &idSet)
	if err != nil {
		return nil, err
	}

	return &graph, nil
}

func ParseJSON(filePath string) (*models.QueryList, error) {
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
