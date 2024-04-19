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

// TODO
func (g *Graph) SaveGraph(db *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

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

// TODO
func LoadGraph(id string) (*Graph, error) {
	var graph Graph

	return &graph, nil
}

// TODO
func FindCycles(id string) [][]string {
	return [][]string{}
}
