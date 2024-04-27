package arbs

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/robinmuhia/arbitrageClient/pkg/oddsWrapper/application/enums"
	"github.com/robinmuhia/arbitrageClient/pkg/oddsWrapper/domain"
)

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

func TestGetAllArbs(t *testing.T) {
	type args struct {
		ctx       context.Context
		arbParams ArbsParams
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
				ctx:       context.Background(),
				arbParams: ArbsParams{},
			},
			want:    1,
			want1:   1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			got, got1, err := GetAllArbs(tt.args.ctx, tt.args.arbParams)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllArbs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("GetAllArbs() got = %v, want %v", len(got), tt.want)
			}

			if !reflect.DeepEqual(len(got1), tt.want1) {
				t.Errorf("GetAllArbs() got1 = %v, want %v", len(got1), tt.want1)
			}
		})
	}
}
