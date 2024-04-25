package enums

type params string

const (
	ApiKey     params = "apiKey"
	Region     params = "regions"
	Markets    params = "markets"
	OddsFormat params = "oddsFormat"
	DateFormat params = "dateFormat"
)

func (e params) String() string {
	return string(e)
}
