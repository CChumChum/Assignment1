package structs

type CountryName struct {
	CountryName string `json:"common"`
}

type CountryFlag struct {
	Png string `json:"png"`
}

type Cities struct {
	Cities []string `json:"data"`
}

type RestCountriesResponse struct {
	Name       CountryName       `json:"name"`
	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Languages  map[string]string `json:"languages"`
	Bordering  []string          `json:"borders"`
	Capital    []string          `json:"capital"`
	Flag       CountryFlag       `json:"flags"`
}

type InfoResponse struct {
	Name       string            `json:"name"`
	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Languages  map[string]string `json:"languages"`
	Bordering  []string          `json:"bordering"`
	Flag       string            `json:"flag"`
	Capital    string            `json:"capital"`
	Cities     []string          `json:"cities"`
}
