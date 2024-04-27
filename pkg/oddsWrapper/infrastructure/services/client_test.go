package services

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/robinmuhia/arbitrageClient/pkg/oddsWrapper/application/enums"
	"github.com/robinmuhia/arbitrageClient/pkg/oddsWrapper/domain"
)

func Test_oddsAPIHTTPClient_makeRequest(t *testing.T) {
	type args struct {
		ctx         context.Context
		method      string
		urlPath     string
		queryParams url.Values
		in4         interface{}
	}

	queryParams := url.Values{}
	queryParams.Add("foo", "bar")

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: successful request",
			args: args{
				ctx:         context.Background(),
				method:      http.MethodGet,
				urlPath:     "https://www.foo.com",
				queryParams: queryParams,
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid http method",
			args: args{
				ctx:         context.Background(),
				method:      http.MethodPost,
				urlPath:     "https://www.foo.com",
				queryParams: queryParams,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "happy case: successful request" {
				httpmock.RegisterResponder(http.MethodGet, tt.args.urlPath, func(r *http.Request) (*http.Response, error) { //nolint:all
					return httpmock.NewJsonResponse(http.StatusOK, nil)
				})
			}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			s, _ := NewServiceOddsAPI()
			resp, err := s.makeRequest(tt.args.ctx, tt.args.method, tt.args.urlPath, tt.args.queryParams, tt.args.in4)

			if err == nil {
				defer resp.Body.Close()
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("oddsAPIHTTPClient.makeRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_oddsAPIHTTPClient_getSports(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: successful retrieval of sport",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to decode sport",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "sad case: unsuccessful retrieval of sport",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewServiceOddsAPI()
			urlPath := fmt.Sprintf("%s/%s", s.baseURL, enums.Sport.String())

			if tt.name == "happy case: successful retrieval of sport" {
				httpmock.RegisterResponder(http.MethodGet, urlPath, func(r *http.Request) (*http.Response, error) {
					resp := []domain.Sport{
						{
							Key:   "foo",
							Title: "bar",
						},
					}

					return httpmock.NewJsonResponse(http.StatusOK, resp)
				})
			}

			if tt.name == "sad case: unsuccessful retrieval of sport" {
				httpmock.RegisterResponder(http.MethodGet, urlPath, func(r *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(http.StatusUnauthorized, nil)
				})
			}

			if tt.name == "sad case: unable to decode sport" {
				httpmock.RegisterResponder(http.MethodGet, urlPath, func(r *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(http.StatusOK, "nana")
				})
			}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			_, err := s.getSports(tt.args.ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("oddsAPIHTTPClient.makeRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_oddsAPIHTTPClient_getOdd(t *testing.T) {
	type args struct {
		ctx       context.Context
		oddParams OddsParams
		sport     domain.Sport
		wg        *sync.WaitGroup
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: successful retrieval of odd",
			args: args{
				ctx: context.Background(),
				oddParams: OddsParams{
					Region:     "foo",
					Markets:    "bar",
					OddsFormat: "foo",
					DateFormat: "bar",
				},
				sport: domain.Sport{
					Key: "foo",
				},
				wg: &sync.WaitGroup{},
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to decode odd",
			args: args{
				ctx: context.Background(),
				oddParams: OddsParams{
					Region:     "foo",
					Markets:    "bar",
					OddsFormat: "foo",
					DateFormat: "bar",
				},
				sport: domain.Sport{
					Key: "foo",
				},
				wg: &sync.WaitGroup{},
			},
			wantErr: true,
		},
		{
			name: "sad case: unsuccessful retrieval of odd",
			args: args{
				ctx: context.Background(),
				oddParams: OddsParams{
					Region:     "foo",
					Markets:    "bar",
					OddsFormat: "foo",
					DateFormat: "bar",
				},
				sport: domain.Sport{
					Key: "foo",
				},
				wg: &sync.WaitGroup{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewServiceOddsAPI()

			urlPath := fmt.Sprintf("%s/%s/%s/%s", s.baseURL, enums.Sport.String(), tt.args.sport.Key, enums.Odds.String())

			tt.args.wg.Add(1)

			if tt.name == "happy case: successful retrieval of odd" {
				httpmock.RegisterResponder(http.MethodGet, urlPath, func(r *http.Request) (*http.Response, error) {
					resp := []domain.Odds{{
						ID:       "foo",
						SportKey: "bar",
					},
					}

					return httpmock.NewJsonResponse(http.StatusOK, resp)
				})
			}

			if tt.name == "sad case: unsuccessful retrieval odd" {
				httpmock.RegisterResponder(http.MethodGet, urlPath, func(r *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(http.StatusUnauthorized, nil)
				})
			}

			if tt.name == "sad case: unable to decode odd" {
				httpmock.RegisterResponder(http.MethodGet, urlPath, func(r *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(http.StatusOK, "nana")
				})
			}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			_, err := s.getOdd(tt.args.ctx, tt.args.oddParams, tt.args.sport, tt.args.wg)
			if (err != nil) != tt.wantErr {
				t.Errorf("oddsAPIHTTPClient.getOdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_oddsAPIHTTPClient_GetAllOdds(t *testing.T) {
	type args struct {
		ctx        context.Context
		oddsParams OddsParams
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: successfully get all odds",
			args: args{
				ctx: context.Background(),
				oddsParams: OddsParams{
					Region:     "foo",
					Markets:    "bar",
					OddsFormat: "foo",
					DateFormat: "bar",
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: unsuccessfully get all odds",
			args: args{
				ctx: context.Background(),
				oddsParams: OddsParams{
					Region:     "foo",
					Markets:    "bar",
					OddsFormat: "foo",
					DateFormat: "bar",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewServiceOddsAPI()

			sportPath := fmt.Sprintf("%s/%s", s.baseURL, enums.Sport.String())

			if tt.name == "happy case: successfully get all odds" {
				resp := []domain.Sport{
					{
						Key:    "foo",
						Title:  "bar",
						Active: true,
					},
					{
						Key:    "bar",
						Title:  "foo",
						Active: true,
					},
				}

				httpmock.RegisterResponder(http.MethodGet, sportPath, func(r *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(http.StatusOK, resp)
				})

				oddPath1 := fmt.Sprintf("%s/%s/%s/%s", s.baseURL, enums.Sport.String(), resp[0].Key, enums.Odds.String())
				httpmock.RegisterResponder(http.MethodGet, oddPath1, func(r *http.Request) (*http.Response, error) {
					oddResp1 := []domain.Odds{
						{
							ID:       "foo",
							SportKey: "bar",
						},
					}

					return httpmock.NewJsonResponse(http.StatusOK, oddResp1)
				})

				oddPath2 := fmt.Sprintf("%s/%s/%s/%s", s.baseURL, enums.Sport.String(), resp[1].Key, enums.Odds.String())
				httpmock.RegisterResponder(http.MethodGet, oddPath2, func(r *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(http.StatusNotFound, nil)
				})
			}

			if tt.name == "sad case: unsuccessfully get all odds" {
				httpmock.RegisterResponder(http.MethodGet, sportPath, func(r *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(http.StatusUnauthorized, nil)
				})
			}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			_, err := s.GetAllOdds(tt.args.ctx, tt.args.oddsParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("oddsAPIHTTPClient.GetAllOdds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewServiceOddsAPI(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "happy case: successful instantiation",
			wantErr: false,
		},
		{
			name:    "sad case: no api key provided",
			wantErr: true,
		},
		{
			name:    "sad case: no base URL provided",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad case: no base URL provided" {
				baseURL = ""
			}

			if tt.name == "sad case: no api key provided" {
				apiKey = ""
			}

			_, err := NewServiceOddsAPI()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServiceOddsAPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
