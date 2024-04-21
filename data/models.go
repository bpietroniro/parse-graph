package data

type Graph struct {
	ID    string  `xml:"id"`
	Name  string  `xml:"name"`
	Nodes []*Node `xml:"nodes>node"`
	Edges []*Edge `xml:"edges>node"`
}

type Node struct {
	ID   string `xml:"id"`
	Name string `xml:"name"`
}

type Edge struct {
	ID     string  `xml:"id"`
	FromID string  `xml:"from"`
	ToID   string  `xml:"to"`
	Cost   float64 `xml:"cost"`
}

type Inputs struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Query struct {
	Paths    *Inputs `json:"paths"`
	Cheapest *Inputs `json:"cheapest"`
}

type QueryList struct {
	Queries []Query `json:"queries"`
}

type Result struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type PathResult struct {
	Result
	Paths [][]string `json:"paths"`
}

type CheapestResult struct {
	Result
	Path []string `json:"path"`
}
