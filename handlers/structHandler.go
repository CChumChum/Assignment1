package handlers

type CountryName struct {
	CountryName string `json:"common"`
}

type CountryFlag struct {
	Png string `json:"png"`
}

type Cities struct {
	Cities []string `json:"data"`
}

type RestCountries struct {
	Name       CountryName       `json:"name"`
	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Languages  map[string]string `json:"languages"`
	Borders    []string          `json:"borders"`
	Capital    []string          `json:"capital"`
	Flag       CountryFlag       `json:"flags"`
}

type InfoResponse struct {
	Name       CountryName       `json:"name"`
	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Languages  map[string]string `json:"languages"`
	Bordering  []string          `json:"bordering"`
	Flag       string            `json:"flag"`
	Capital    string            `json:"capital"`
	Cities     []string          `json:"cities"`
}

type PopulationData struct {
	PopulationInfo []map[string]int `json:"populationCounts"`
}

type GetPopulationData struct {
	Data PopulationData `json:"data"`
}

type PopulationInfoResponse struct {
	Mean   int              `json:"mean"`
	Values []map[string]int `json:"values"`
}

type StatusResponse struct {
	CountriesNowApi  int    `json:"countriesNowApi"`
	RestCountriesApi int    `json:"restCountriesApi"`
	Version          string `json:"version"`
	Uptime           string `json:"uptime"`
}
