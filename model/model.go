package model

type ImportConfiguration struct {
	Packages []Package `json:"packages"`
}

type ImportImage struct {
	Name string `json:"name"`
	From string `json:"from"`
}

type Package struct {
	File    string        `json:"file"`
	Content []ImportImage `json:"content"`
}
