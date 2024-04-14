package main

type Graph struct {
	ID    string
	Name  string
	Nodes []Node
	Edges []Edge
}

type Node struct {
	ID   string
	Name string
}

type Edge struct {
	ID     string
	FromID string
	ToID   string
	Cost   float64
}

func main() {

}
