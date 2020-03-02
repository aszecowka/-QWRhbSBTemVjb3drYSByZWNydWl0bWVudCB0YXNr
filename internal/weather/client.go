package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

//go:generate mockery -name=HTTPClient -output=automock -outpkg=automock -case=underscore
type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}

func NewClient(httpClient HTTPClient, log logrus.FieldLogger, baseURL string, apiKey string, reqTimeout time.Duration) *restClient {
	return &restClient{
		httpClient:     httpClient,
		baseURL:        baseURL,
		apiKey:         apiKey,
		requestTimeout: reqTimeout,
		log:            log.WithField("service", "openweathermap-client"),
	}
}

type restClient struct {
	baseURL        string
	apiKey         string
	requestTimeout time.Duration
	httpClient     HTTPClient
	log            logrus.FieldLogger
}

func (c *restClient) Get(ctx context.Context, city string) (_ OpenWeatherResponse, err error) {
	if ctx.Err() != nil {
		return OpenWeatherResponse{}, ctx.Err()
	}
	ctx, releaseCtx := context.WithTimeout(ctx, c.requestTimeout)
	defer releaseCtx()
	url := fmt.Sprintf("%s/data/2.5/weather", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return OpenWeatherResponse{}, fmt.Errorf("while creating request: %w", err)
	}

	q := req.URL.Query()
	q.Set("q", city)
	q.Set("appid", c.apiKey)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return OpenWeatherResponse{}, fmt.Errorf("while executing request: %w", err)
	}
	defer func() {
		drainErr := c.drainAndClose(resp.Body)
		// ignore drainErr if there were other errors
		if err == nil {
			err = drainErr
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return OpenWeatherResponse{}, NotFoundError
	default:
		return OpenWeatherResponse{}, fmt.Errorf("wrong status code: got [%d], expected [%d]", resp.StatusCode, http.StatusOK)
	}

	fetchResponse := OpenWeatherResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&fetchResponse); err != nil {
		return OpenWeatherResponse{}, errors.Wrap(err, "while decoding response")
	}
	c.log.Infof("Successfully fetched data for city: [%s]", city) // only to show when cache is used
	return fetchResponse, nil

}

func (c *restClient) drainAndClose(r io.ReadCloser) error {
	_, copyErr := io.Copy(ioutil.Discard, r)
	closeErr := r.Close()
	// error from io.Copy likely more useful than the one from Close
	if copyErr != nil {
		return fmt.Errorf("while draining response body: %w", copyErr)
	}
	if closeErr != nil {
		return fmt.Errorf("while closing response body: %w", closeErr)
	}
	return nil
}
