package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/robinmuhia/arbitrageClient/pkg/oddsWrapper/application/enums"
	"github.com/robinmuhia/arbitrageClient/pkg/oddsWrapper/domain"
)

var (
	baseURL = os.Getenv(enums.BaseURL.String())
	apiKey  = os.Getenv(enums.ApiKeyEnv.String())
)

// oddsoddsAPIHTTPClient instantiates a client to call the odds API url
type oddsAPIHTTPClient struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

// ArbClient implements methods intended to be exposed by the oddsHTTPClient
type ArbClient interface {
	GetAllOdds(ctx context.Context, oddsParams OddsParams) ([]domain.Odds, error)
}

// OddsParams represent the parameters required to query for specific odds
type OddsParams struct {
	Region     string
	Markets    string
	OddsFormat string
	DateFormat string
}

// NewServiceOddsAPI returns a new instance of an OddsAPI service
func NewServiceOddsAPI() (*oddsAPIHTTPClient, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("empty env variables, %s", enums.BaseURL.String())
	}

	if apiKey == "" {
		return nil, fmt.Errorf("empty env variables, %s", enums.ApiKeyEnv.String())
	}

	return &oddsAPIHTTPClient{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		baseURL: baseURL,
		apiKey:  apiKey,
	}, nil
}

// makeRequest calls the Odds API endpoint
func (s *oddsAPIHTTPClient) makeRequest(ctx context.Context, method, urlPath string, queryParams url.Values, _ interface{}) (*http.Response, error) {
	var request *http.Request

	switch method {
	case http.MethodGet:
		req, err := http.NewRequestWithContext(ctx, method, urlPath, nil)
		if err != nil {
			return nil, err
		}

		request = req

	default:
		return nil, fmt.Errorf("unsupported http method: %s", method)
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	if queryParams != nil {
		request.URL.RawQuery = queryParams.Encode()
	}

	return s.client.Do(request)
}

// getSports returns a list of sports and an error
func (s oddsAPIHTTPClient) getSports(ctx context.Context) ([]domain.Sport, error) {
	urlPath := fmt.Sprintf("%s/%s", s.baseURL, enums.Sport.String())

	queryParams := url.Values{}
	queryParams.Add(enums.ApiKey.String(), s.apiKey)

	resp, err := s.makeRequest(ctx, http.MethodGet, urlPath, queryParams, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get sports data: %s", resp.Status)
	}

	var sports []domain.Sport
	if err := json.NewDecoder(resp.Body).Decode(&sports); err != nil {
		return nil, fmt.Errorf("failed to get sports: %w", err)
	}

	return sports, nil
}

// getOdd returns all odds from one sport
func (s oddsAPIHTTPClient) getOdd(ctx context.Context, oddParams OddsParams, sport domain.Sport, wg *sync.WaitGroup) ([]domain.Odds, error) {
	defer wg.Done()

	urlPath := fmt.Sprintf("%s/%s/%s/%s", s.baseURL, enums.Sport.String(), sport.Key, enums.Odds.String())

	queryParams := url.Values{}
	queryParams.Add(enums.ApiKey.String(), s.apiKey)
	queryParams.Add(enums.Region.String(), oddParams.Region)
	queryParams.Add(enums.Markets.String(), oddParams.Markets)
	queryParams.Add(enums.OddsFormat.String(), oddParams.OddsFormat)
	queryParams.Add(enums.DateFormat.String(), oddParams.DateFormat)

	resp, err := s.makeRequest(ctx, http.MethodGet, urlPath, queryParams, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get sports data: %s", resp.Status)
	}

	var odds []domain.Odds
	if err := json.NewDecoder(resp.Body).Decode(&odds); err != nil {
		return nil, fmt.Errorf("failed to decode odds data for %s: %w", sport.Title, err)
	}

	return odds, nil
}

// GetOdds returns a list of all available odds given various parameters across all sports
func (s oddsAPIHTTPClient) GetAllOdds(ctx context.Context, oddsParams OddsParams) ([]domain.Odds, error) {
	sports, err := s.getSports(ctx)
	if err != nil {
		return nil, err
	}

	ticker := time.NewTicker(1 * time.Second)

	var wg sync.WaitGroup

	var mu sync.Mutex

	var allOdds []domain.Odds

	for _, sport := range sports {
		if sport.Active {
			wg.Add(1)

			go func() {
				odds, err := s.getOdd(ctx, oddsParams, sport, &wg)
				if err != nil {
					log.Print(err.Error())
				}

				mu.Lock()
				allOdds = append(allOdds, odds...)
				mu.Unlock()
			}()
		}

		<-ticker.C // waits a second to send next goroutine, intended to prevent ddosing and rate limiting
	}

	wg.Wait()

	ticker.Stop()

	return allOdds, nil
}
