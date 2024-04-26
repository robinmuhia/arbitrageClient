package enums

type env string

const (
	BaseURL   env = "ODDS_API_BASE_URL"
	ApiKeyEnv env = "ODDS_API_KEY" //nolint: gosec
)

func (e env) String() string {
	return string(e)
}
