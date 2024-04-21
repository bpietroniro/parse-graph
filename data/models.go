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

type PathInputs struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type PathQuery struct {
	Paths    *PathInputs `json:"paths"`
	Cheapest *PathInputs `json:"cheapest"`
}

type QueryList struct {
	Queries []PathQuery `json:"queries"`
}
