package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"parse-graph/data"
	"parse-graph/utils"

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

	var graph data.Graph
	err = parseXML(filePath, &graph)
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

func ConnectToDB() (*pgxpool.Pool, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return nil, err
	}

	connectionString := fmt.Sprintf("postgres://%s@%s:%s/%s", os.Getenv("DB_USER"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	fmt.Println(os.Getenv("DATABASE_URL"))
	dbpool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		return nil, err
	}

	// var greeting string
	// err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 	return nil, err
	// }

	// fmt.Println(greeting)
	return dbpool, nil
}

func OpenDB() {}
