package models

// for JSON input

type QueryInputs struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Query struct {
	Paths    *QueryInputs `json:"paths,omitempty"`
	Cheapest *QueryInputs `json:"cheapest,omitempty"`
}

type QueryList struct {
	Queries []Query `json:"queries"`
}

// for JSON output

type Result struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Cost float64 `json:"-"`
}

type PathResultList struct {
	From  string     `json:"from"`
	To    string     `json:"to"`
	Paths [][]string `json:"paths"`
}

type PathResult struct {
	From string   `json:"from"`
	To   string   `json:"to"`
	Cost string   `json:"-"`
	Path []string `json:"path"`
}

type CheapestResult struct {
	From string      `json:"from"`
	To   string      `json:"to"`
	Path interface{} `json:"path"`
}
