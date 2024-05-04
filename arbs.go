package arbs

import (
	"context"
	"log"

	"github.com/robinmuhia/arbitrageClient/pkg/oddsWrapper/domain"
	"github.com/robinmuhia/arbitrageClient/pkg/oddsWrapper/infrastructure/services"
	"github.com/robinmuhia/arbitrageClient/pkg/oddsWrapper/usecases/arbs"
)

var client services.ArbClient

type ArbsParams struct {
	Region     string
	Markets    string
	OddsFormat string
	DateFormat string
}

func init() {
	svc, err := services.NewServiceOddsAPI()
	if err != nil {
		log.Panic("error: %w", err)
	}

	client = svc
}

// GetAllArbs returns all possible arbitrage opportunities from the odds API.
func GetAllArbs(ctx context.Context, arbParams ArbsParams) ([]domain.ThreeOddsArb, []domain.TwoOddsArb, error) {
	us := arbs.UseCasesArbsImpl{
		OddsApiClient: client,
	}

	params := services.OddsParams{
		Region:     arbParams.Region,
		Markets:    arbParams.Markets,
		OddsFormat: arbParams.OddsFormat,
		DateFormat: arbParams.DateFormat,
	}

	return us.GetArbs(ctx, params)
}
