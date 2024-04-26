package arbs

import (
	"context"
	"fmt"
	"sync"

	"github.com/robinmuhia/arbitrageClient/pkg/oddsWrapper/domain"
	"github.com/robinmuhia/arbitrageClient/pkg/oddsWrapper/infrastructure/services"
)

type UseCasesArbsImpl struct {
	OddsApiClient services.ArbClient
}

// composeTwoArbsBet process an Odd to check if an arbitrage oppurtunity exists
func (us *UseCasesArbsImpl) composeTwoArbsBet(odd domain.Odds, i int, j int) (domain.TwoOddsArb, bool) {
	homeOdd := odd.Bookmakers[i].Markets[0].Outcomes[0].Price
	awayOdd := odd.Bookmakers[j].Markets[0].Outcomes[1].Price
	arb := (1 / homeOdd) + (1 / awayOdd)

	if arb < 1.0 {
		profit := (1 - arb) * 100
		twowayArb := domain.TwoOddsArb{
			Title:    fmt.Sprintf("%s - %s", odd.HomeTeam, odd.AwayTeam),
			Home:     odd.Bookmakers[i].Title,
			HomeOdds: homeOdd,
			HomeStake: 1/homeOdd,
			Away:     odd.Bookmakers[j].Title,
			AwayOdds: awayOdd,
			AwayStake: 1/awayOdd,
			GameType: odd.SportKey,
			League:   odd.SportTitle,
			Profit:   profit,
			GameTime: odd.CommenceTime,
		}

		return twowayArb, true
	}

	return domain.TwoOddsArb{}, false
}

// composeThreeArbsBet process an Odd to check if an arbitrage oppurtunity exists
func (us *UseCasesArbsImpl) composeThreeArbsBet(odd domain.Odds, i int, j int, k int) (domain.ThreeOddsArb, bool) {
	homeOdd := odd.Bookmakers[i].Markets[0].Outcomes[0].Price
	awayOdd := odd.Bookmakers[j].Markets[0].Outcomes[1].Price
	drawOdd := odd.Bookmakers[k].Markets[0].Outcomes[2].Price
	arb := (1 / homeOdd) + (1 / awayOdd) + (1 / drawOdd)

	if arb < 1.0 {
		profit := (1 - arb) * 100
		threewayArb := domain.ThreeOddsArb{
			Title:    fmt.Sprintf("%s - %s", odd.HomeTeam, odd.AwayTeam),
			Home:     odd.Bookmakers[i].Title,
			HomeOdds: homeOdd,
			HomeStake: 1/homeOdd,
			Away:     odd.Bookmakers[j].Title,
			AwayOdds: awayOdd,
			AwayStake: 1/awayOdd,
			Draw:     odd.Bookmakers[k].Title,
			DrawOdds: drawOdd,
			DrawStake: 1/drawOdd,
			GameType: odd.SportKey,
			League:   odd.SportTitle,
			Profit:   profit,
			GameTime: odd.CommenceTime,
		}

		return threewayArb, true
	}

	return domain.ThreeOddsArb{}, false
}

// checkIfMarketHasEnoughGames checks whether a market has enough games to analyze an arbitrage oppurtunity
func (us *UseCasesArbsImpl) checkIfMarketHasEnoughGames(bookmarker domain.Bookmaker) bool {
	return len(bookmarker.Markets) >= 1
}

func (us *UseCasesArbsImpl) findPossibleArbOppurtunity(odd domain.Odds,
	threeOddsCh chan<- domain.ThreeOddsArb,
	twoOddsCh chan<- domain.TwoOddsArb,
	wg *sync.WaitGroup) {
	defer wg.Done()

	if len(odd.Bookmakers) < 2 {
		return // Skip if there are not enough bookmakers for comparison
	}

	for i := 0; i < len(odd.Bookmakers); i++ {
		if !us.checkIfMarketHasEnoughGames(odd.Bookmakers[i]) {
			return
		}

		for j := 0; j < len(odd.Bookmakers); j++ {
			if !us.checkIfMarketHasEnoughGames(odd.Bookmakers[j]) {
				return
			}

			switch {
			case len(odd.Bookmakers[i].Markets[0].Outcomes) == 2 && len(odd.Bookmakers[j].Markets[0].Outcomes) == 2:
				twoWayArb, isArb := us.composeTwoArbsBet(odd, i, j)
				if isArb {
					twoOddsCh <- twoWayArb
				}

			case len(odd.Bookmakers[i].Markets[0].Outcomes) == 3 && len(odd.Bookmakers[j].Markets[0].Outcomes) == 3:
				for k := 0; k < len(odd.Bookmakers); k++ {
					if !us.checkIfMarketHasEnoughGames(odd.Bookmakers[k]) {
						return
					}

					if len(odd.Bookmakers[k].Markets[0].Outcomes) == 3 {
						threeWayArb, isArb := us.composeThreeArbsBet(odd, i, j, k)
						if isArb {
							threeOddsCh <- threeWayArb
						}
					}
				}
			}
		}
	}
}

func (us *UseCasesArbsImpl) GetArbs(ctx context.Context, oddsParams services.OddsParams) ([]domain.ThreeOddsArb, []domain.TwoOddsArb, error) {
	odds, err := us.OddsApiClient.GetAllOdds(ctx, oddsParams)
	if err != nil {
		return nil, nil, err
	}

	var ThreeOddsArbs []domain.ThreeOddsArb

	var TwoOddsArbs []domain.TwoOddsArb

	// Create channels to receive arbitrage results
	threeOddsCh := make(chan domain.ThreeOddsArb)
	twoOddsCh := make(chan domain.TwoOddsArb)

	// Create a wait group to ensure all goroutines finish before returning
	var wg sync.WaitGroup

	var once sync.Once

	for _, odd := range odds {
		wg.Add(1)

		go us.findPossibleArbOppurtunity(odd, threeOddsCh, twoOddsCh, &wg)
	}

	// Close the channels once all goroutines finish processing
	go func() {
		wg.Wait()
		once.Do(func() {
			close(threeOddsCh)
			close(twoOddsCh)
		})
	}()

	// Collect the results from channels
	for {
		select {
		case arb, ok := <-threeOddsCh:
			if !ok {
				threeOddsCh = nil // Set to nil to exit the loop when both channels are closed
			} else {
				ThreeOddsArbs = append(ThreeOddsArbs, arb)
			}
		case arb, ok := <-twoOddsCh:
			if !ok {
				twoOddsCh = nil // Set to nil to exit the loop when both channels are closed
			} else {
				TwoOddsArbs = append(TwoOddsArbs, arb)
			}
		}
		// Exit the loop when both channels are closed
		if threeOddsCh == nil && twoOddsCh == nil {
			break
		}
	}

	return ThreeOddsArbs, TwoOddsArbs, nil
}
