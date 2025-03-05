package structs

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
