package arbs

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/robinmuhia/arbitrageClient/pkg/oddsWrapper/application/enums"
	"github.com/robinmuhia/arbitrageClient/pkg/oddsWrapper/domain"
	"github.com/robinmuhia/arbitrageClient/pkg/oddsWrapper/infrastructure/services"
)

func setUpUsecase() *UseCasesArbsImpl {
	var client services.ArbClient

	svc, err := services.NewServiceOddsAPI()
	if err != nil {
		log.Panic("error: %w", err)
	}

	client = svc

	return &UseCasesArbsImpl{
		OddsApiClient: client,
	}
}

var (
	twoWayArbodd = domain.Odds{
		HomeTeam:     "Vuvu",
		AwayTeam:     "Zela",
		SportKey:     "baseball_np",
		SportTitle:   "Baseball",
		CommenceTime: "2024-01-01T08:29:59Z",
		Bookmakers: []domain.Bookmaker{
			{
				Title: "Betway",
				Markets: []domain.Market{
					{
						Outcomes: []domain.Outcome{
							{Price: 1.44},
							{Price: 8.5},
						}}}},
			{
				Title: "Betway",
				Markets: []domain.Market{
					{
						Outcomes: []domain.Outcome{
							{Price: 0},
							{Price: 0},
						}}}},
		},
	}
	twoWayNotArbodd = domain.Odds{
		HomeTeam:     "Vuvu",
		AwayTeam:     "Zela",
		SportKey:     "baseball_np",
		SportTitle:   "Baseball",
		CommenceTime: "2024-01-01T08:29:59Z",
		Bookmakers: []domain.Bookmaker{
			{
				Title: "Betway",
				Markets: []domain.Market{
					{
						Outcomes: []domain.Outcome{
							{Price: 1.44},
							{Price: 1.44},
						}}}}},
	}
	threeWayArbodd = domain.Odds{
		HomeTeam:     "Vuvu",
		AwayTeam:     "Zela",
		SportKey:     "baseball_np",
		SportTitle:   "Baseball",
		CommenceTime: "2024-01-01T08:29:59Z",
		Bookmakers: []domain.Bookmaker{
			{
				Title: "Betway",
				Markets: []domain.Market{
					{
						Outcomes: []domain.Outcome{
							{Price: 4.9},
							{Price: 17},
							{Price: 1.57},
						}}}},
			{
				Title: "Betway",
				Markets: []domain.Market{
					{
						Outcomes: []domain.Outcome{
							{Price: 0.0},
							{Price: 0.0},
							{Price: 0.0},
						}}}},
		},
	}
	threeWayNotArbodd = domain.Odds{
		HomeTeam:     "Vuvu",
		AwayTeam:     "Zela",
		SportKey:     "baseball_np",
		SportTitle:   "Baseball",
		CommenceTime: "2024-01-01T08:29:59Z",
		Bookmakers: []domain.Bookmaker{
			{
				Title: "Betway",
				Markets: []domain.Market{
					{
						Outcomes: []domain.Outcome{
							{Price: 4.9},
							{Price: 1.87},
							{Price: 1.57},
						}}}}},
	}
)

func TestUseCasesArbsImpl_composeTwoArbsBet(t *testing.T) {
	type args struct {
		odd domain.Odds
		i   int
		j   int
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "happy case: arb is found",
			args: args{
				odd: twoWayArbodd,
				i:   0,
				j:   0,
			},
			want: true,
		},
		{
			name: "sad case: arb is not found",
			args: args{
				odd: twoWayNotArbodd,
				i:   0,
				j:   0,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := setUpUsecase()
			_, got := us.composeTwoArbsBet(tt.args.odd, tt.args.i, tt.args.j)

			if got != tt.want {
				t.Errorf("UseCasesArbsImpl.composeTwoArbsBet() got1 = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesArbsImpl_composeThreeArbsBet(t *testing.T) {
	type args struct {
		odd domain.Odds
		i   int
		j   int
		k   int
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "happy case: arb is found",
			args: args{
				odd: threeWayArbodd,
				i:   0,
				j:   0,
				k:   0,
			},
			want: true,
		},
		{
			name: "sad case: arb is not found",
			args: args{
				odd: threeWayNotArbodd,
				i:   0,
				j:   0,
				k:   0,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := setUpUsecase()
			_, got := us.composeThreeArbsBet(tt.args.odd, tt.args.i, tt.args.j, tt.args.k)

			if got != tt.want {
				t.Errorf("UseCasesArbsImpl.composeThreeArbsBet() got1 = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesArbsImpl_GetArbs(t *testing.T) {
	type args struct {
		ctx        context.Context
		oddsParams services.OddsParams
	}

	tests := []struct {
		name    string
		args    args
		want    int
		want1   int
		wantErr bool
	}{
		{
			name: "happy case: get arbs",
			args: args{
				ctx:        context.Background(),
				oddsParams: services.OddsParams{},
			},
			want:    1,
			want1:   1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := setUpUsecase()

			baseUrl := os.Getenv(enums.BaseURL.String())

			sportPath := fmt.Sprintf("%s/%s", baseUrl, enums.Sport.String())

			resp := []domain.Sport{
				{
					Key:    "foo",
					Title:  "bar",
					Active: true,
				},
			}

			httpmock.RegisterResponder(http.MethodGet, sportPath, func(r *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(http.StatusOK, resp)
			})

			oddPath := fmt.Sprintf("%s/%s/%s/%s", baseUrl, enums.Sport.String(), resp[0].Key, enums.Odds.String())
			httpmock.RegisterResponder(http.MethodGet, oddPath, func(r *http.Request) (*http.Response, error) {
				oddResp := []domain.Odds{}
				oddResp = append(oddResp, twoWayArbodd)
				oddResp = append(oddResp, twoWayNotArbodd)
				oddResp = append(oddResp, threeWayArbodd)
				oddResp = append(oddResp, threeWayNotArbodd)

				return httpmock.NewJsonResponse(http.StatusOK, oddResp)
			})

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			got, got1, err := us.GetArbs(tt.args.ctx, tt.args.oddsParams)

			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesArbsImpl.GetArbs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("UseCasesArbsImpl.GetArbs() got = %v, want %v", len(got), tt.want)
			}

			if !reflect.DeepEqual(len(got1), tt.want1) {
				t.Errorf("UseCasesArbsImpl.GetArbs() got1 = %v, want %v", len(got1), tt.want1)
			}
		})
	}
}
