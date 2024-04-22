package models

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

type Result struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Cost float64 `json:"-"`
}

type PathResultList struct {
	Result
	Paths [][]string `json:"paths"`
}

type PathResult struct {
	Result
	Path []string `json:"path"`
}
