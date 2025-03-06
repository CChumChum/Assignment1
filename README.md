# Country Info API

This project provides a RESTful API to fetch country-related information, including general details, population statistics, and service status.

## Features
- **Country Info**: Retrieves details such as name, flag, population, capital, languages, and bordering countries.
- **Population Data**: Fetches historical population data with optional filtering by year range.
- **Service Status**: Reports the uptime and availability of external APIs.

## Endpoints

### 1. Get Country Information

**Endpoint:** `GET /info/{isoCode}`

Fetches country details using a 2-letter ISO code.

**Query Parameters:**

- `limit` (optional) - Limits the number of cities returned

**Response:**

```json
{
  "name": "Norway",
  "continents": ["Europe"],
  "population": 5379475,
  "languages": {"nor": "Norwegian"},
  "bordering": ["SWE", "FIN", "RUS"],
  "flag": "https://flagcdn.com/w320/no.png",
  "capital": "Oslo",
  "cities": ["Bergen", "Trondheim"]
}
```

---

### 2. Get Population Data

**Endpoint:** `GET /population/{isoCode}`

Retrieves population data for a given country. Optional limit to specify a year range.

**Query Parameters:**

- `limit` (optional) - Define a year range in the format `YYYY-YYYY`

**Response:**

```json
{
  "mean": 5320000,
  "values": [
    {"year": 2000, "value": 4800000},
    {"year": 2010, "value": 5000000}
  ]
}
```

---

### 3. API Status Check

**Endpoint:** `GET /status`

Returns API availability and uptime.

**Response:**

```json
{
  "CountriesNowApi": 200,
  "RestCountriesApi": 200,
  "Version": "1.0.0",
  "Uptime": "120.45s"
}
```



##


## Notes
- The API relies on external services (`RestCountriesAPI`, `CountriesNowAPI`). Ensure they are reachable.
- Error handling and validation checks are included to manage invalid inputs.
- The comments in the code are written short by ChatGPT

