package domain

// Outcome the name and price of an individual outcome of a bet eg. Bayern 1.26
type Outcome struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Market represents the bookmarkers' odds for a game
type Market struct {
	Key        string    `json:"key"`
	LastUpdate string    `json:"last_update"`
	Outcomes   []Outcome `json:"outcomes"`
}

// Bookmaker describes the bookmarker such as bet365
type Bookmaker struct {
	Key        string   `json:"key"`
	Title      string   `json:"title"`
	LastUpdate string   `json:"last_update"`
	Markets    []Market `json:"markets"`
}

// Odds represent the odds structure of the games's odds
type Odds struct {
	ID           string      `json:"id"`
	SportKey     string      `json:"sport_key"`
	SportTitle   string      `json:"sport_title"`
	CommenceTime string      `json:"commence_time"`
	HomeTeam     string      `json:"home_team"`
	AwayTeam     string      `json:"away_team"`
	Bookmakers   []Bookmaker `json:"bookmakers"`
}
