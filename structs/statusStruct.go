package structs

type StatusResponse struct {
	CountriesNowApi  int    `json:"countriesNowApi"`
	RestCountriesApi int    `json:"restCountriesApi"`
	Version          string `json:"version"`
	Uptime           string `json:"uptime"`
}
