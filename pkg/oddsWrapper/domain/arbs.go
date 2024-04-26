package domain

// ThreeOddsArb represents the structure of how a match with three possible outcomes
// i.e a win, draw or loss will be represented in the response
type ThreeOddsArb struct {
	Title    string
	Home     string
	HomeOdds float64
	HomeStake float64
	Draw     string
	DrawOdds float64
	DrawStake float64
	Away     string
	AwayStake float64
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
	HomeOdds float64
	HomeStake float64
	Away     string
	AwayStake float64
	AwayOdds float64
	GameType string
	League   string
	Profit   float64
	GameTime string
}
