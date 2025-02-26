package handlers

type yearlyPop struct {
	Year  string `json:"year,omitempty"`
	Value int    `json:"value,omitempty"`
}

type Info struct {
	Name       string `json:"name"`
	Continent  string `json:"continents"`
	Population string `json:"population"`
	Languages  string `json:"languages,"`
	Bordering  string `json:"bordering,omitempty"`
	Flag       string `json:"flag"`
	Capital    string `json:"capital"`
	Cities     string `json:"cities"`
}

type Population struct {
	Mean   int       `json:"mean"`
	Values yearlyPop `json:"values"`
}
