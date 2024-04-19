// database queries and utilities

package data

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const dbTimeout = time.Second * 5

// TODO handle the case that the graph already exists in the db
func (g *Graph) SaveGraph(db *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// executes the insertion of graph, nodes, and edges as a single transaction (to avoid partial insertions)
	transaction, err := db.Begin(ctx)
	if err != nil {
		fmt.Println("Failed to begin transaction")
		return err
	}
	defer transaction.Rollback(ctx)

	stmt := `
	INSERT INTO graphs (id, graph_name)
	VALUES ($1, $2)
	`

	_, err = transaction.Exec(ctx, stmt, g.ID, g.Name)
	if err != nil {
		fmt.Println("Saving graph failed")
		return err
	}

	stmt = `
	INSERT INTO nodes (id, node_name, graph_id)
	VALUES ($1, $2, $3)
	`

	batch := &pgx.Batch{}

	for _, node := range g.Nodes {
		batch.Queue(stmt, node.ID, node.Name, g.ID)
	}

	batchResults := transaction.SendBatch(ctx, batch)

	if err := batchResults.Close(); err != nil {
		fmt.Println("Failed to complete batch execution: saving nodes")
		return err
	}

	stmt = `
	INSERT INTO edges (id, graph_id, from_node, to_node, cost)
	VALUES ($1, $2, $3, $4, $5)
	`

	batch = &pgx.Batch{}

	for _, edge := range g.Edges {
		batch.Queue(stmt, edge.ID, g.ID, edge.FromID, edge.ToID, edge.Cost)
	}

	batchResults = transaction.SendBatch(ctx, batch)

	if err := batchResults.Close(); err != nil {
		fmt.Println("Failed to complete batch execution: saving edges")
		return err
	}

	if err := transaction.Commit(ctx); err != nil {
		fmt.Println("Failed to commit transaction")
		return err
	}

	fmt.Println("Saving graph was successful!")
	return nil
}

// TODO test this
func LoadGraph(db *pgxpool.Pool, id string) (*Graph, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var graph Graph

	query := `SELECT * FROM graphs WHERE id = $1`
	err := db.QueryRow(ctx, query, id).Scan(&graph.ID, &graph.Name)
	if err != nil {
		return nil, err
	}

	query = `SELECT id, node_name FROM nodes WHERE graph_id = $1`
	rows, err := db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var node Node

		err := rows.Scan(&node.ID, &node.Name)
		if err != nil {
			return nil, err
		}

		graph.Nodes = append(graph.Nodes, &node)
	}

	query = `SELECT id, from_node, to_node, cost FROM edges WHERE graph_id = $1`
	rows, err = db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var edge Edge

		err := rows.Scan(&edge.ID, &edge.FromID, &edge.ToID, &edge.Cost)
		if err != nil {
			return nil, err
		}

		graph.Edges = append(graph.Edges, &edge)
	}

	return &graph, nil
}

// TODO
func FindCycles(id string) [][]string {
	return [][]string{}
}
