package enums

type endpoint string

const (
	Sport endpoint = "sports"
	Odds  endpoint = "odds"
)

func (e endpoint) String() string {
	return string(e)
}
