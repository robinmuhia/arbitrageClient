package domain

// ThreeOddsArb represents the structure of how a match with three possible outcomes
// i.e a win, draw or loss will be represented in the response
type ThreeOddsArb struct {
	Title    string
	Home     string
	Draw     string
	Away     string
	HomeOdds float64
	DrawOdds float64
	AwayOdds float64
	GameType string
	League   string
	Profit   float64
	GameTime string
}

// TwoOddsArb represents the structure of how a match with three possible outcomes
// i.e a win or loss will be represented in the response
type TwoOddsArb struct {
	Title    string
	Home     string
	Away     string
	HomeOdds float64
	AwayOdds float64
	GameType string
	League   string
	Profit   float64
	GameTime string
}
