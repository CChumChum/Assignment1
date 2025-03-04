package handlers

const EMPTY_PATH = "/"
const DEFAULT_PATH = "/countryinfo/v1/"
const INFO_PATH = DEFAULT_PATH + "info/"
const POPULATION_PATH = DEFAULT_PATH + "population/"
const STATUS_PATH = DEFAULT_PATH + "status/"

const CountriesNowAPI = "http://129.241.150.113:3500/api/v0.1/"
const RestCountriesAPI = "http://129.241.150.113:8080/v3.1/"
const JsonHeader = "application/json"

const GENERIC_SERVER_ERROR = "An unexpected error occurred. Please try again.!"
const REQUEST_SERVER_ERROR = "Unable to process the request due to an internal error."
const RESPONSE_SERVER_ERROR = "The server encountered an issue while processing the response."
const ENCODE_SERVER_ERROR = "Failed to encode the response payload."
const DECODE_SERVER_ERROR = "Failed to decode the response from the external service."
const INTEGER_FAULT = "Invalid 'limit' parameter. It must be a positive integer."
