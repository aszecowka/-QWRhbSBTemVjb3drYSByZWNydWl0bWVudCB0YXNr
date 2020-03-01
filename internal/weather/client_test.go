package weather_test

import (
	"context"
	"encoding/json"
	"github.com/aszecowka/QWRhbSBTemVjb3drYSByZWNydWl0bWVudCB0YXNr/internal/weather"
	"github.com/aszecowka/QWRhbSBTemVjb3drYSByZWNydWl0bWVudCB0YXNr/internal/weather/automock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient(t *testing.T) {

	givenTimeout := time.Second
	ctx := context.Background()
	t.Run("Success", func(t *testing.T) {
		// GIVEN
		testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			assert.Equal(t, http.MethodGet, req.Method)
			assert.Equal(t, "/data/2.5/weather", req.URL.Path)
			assert.Equal(t, fixAPIKey(), req.URL.Query().Get("appid"))
			assert.Equal(t, "Paris", req.URL.Query().Get("q"))
			assert.Equal(t, "application/json", req.Header.Get("Accept"))

			err := json.NewEncoder(rw).Encode(fixWeatherInParis())
			require.NoError(t, err)

		}))
		defer testServer.Close()
		sut := weather.NewClient(http.DefaultClient, testServer.URL, fixAPIKey(), givenTimeout)
		// WHEN
		actualResponse, err := sut.Get(ctx, "Paris")
		// THEN
		require.NoError(t, err)
		assert.Equal(t, *fixWeatherInParis(), actualResponse)
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		// GIVEN
		mockHTTP := automock.NewClientThatReturnsNotFoundStatusCode()
		defer mockHTTP.AssertExpectations(t)
		sut := weather.NewClient(mockHTTP, "", "", givenTimeout)
		// WHEN
		_, err := sut.Get(ctx, "Paris")
		// THEN
		assert.Equal(t, err, weather.NotFoundError)
	})

	t.Run("Error - Wrong status code", func(t *testing.T) {
		// GIVEN
		mockHTTP := automock.NewClientThatReturnsWrongStatusCode()
		defer mockHTTP.AssertExpectations(t)
		sut := weather.NewClient(mockHTTP, "", "", givenTimeout)
		// WHEN
		_, err := sut.Get(ctx, "Paris")
		// THEN
		assert.EqualError(t, err, "wrong status code: got [500], expected [200]")
	})

	t.Run("Error - HTTP client returns error", func(t *testing.T) {
		// GIVEN
		mockHTTP := automock.NewClientThatReturnsError()
		defer mockHTTP.AssertExpectations(t)
		sut := weather.NewClient(mockHTTP, "", "", givenTimeout)
		// WHEN
		_, err := sut.Get(ctx, "Paris")
		// THEN
		assert.EqualError(t, err, "while executing request: some error")
	})

	t.Run("Error on draining and closing body", func(t *testing.T) {
		// GIVEN
		mockHTTP := automock.NewClientThatReturnsErrorOnClosingBody(http.StatusOK)
		defer mockHTTP.AssertExpectations(t)
		sut := weather.NewClient(mockHTTP, "", "", givenTimeout)
		// WHEN
		_, err := sut.Get(ctx, "Paris")
		// THEN
		assert.EqualError(t, err, "while closing response body: close error")

	})

	t.Run("Error on draining and closing body does not override previous errors", func(t *testing.T) {
		// GIVEN
		mockHTTP := automock.NewClientThatReturnsErrorOnClosingBody(http.StatusInternalServerError)
		defer mockHTTP.AssertExpectations(t)
		sut := weather.NewClient(mockHTTP, "", "", givenTimeout)

		// WHEN
		_, err := sut.Get(ctx, "Paris")
		// THEN
		assert.EqualError(t, err, "wrong status code: got [500], expected [200]")
	})

	t.Run("Request has defined timeout", func(t *testing.T) {
		// GIVEN
		mockHTTP := automock.NewClientThatChecksIfRequestHasDefinedDeadline(http.StatusOK)
		defer mockHTTP.AssertExpectations(t)
		sut := weather.NewClient(mockHTTP, "", "", givenTimeout)

		// WHEN
		_, err := sut.Get(ctx, "Paris")

		// THEN
		require.NoError(t, err)
	})

	t.Run("Fails fast on cancelled context", func(t *testing.T) {
		// GIVEN
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()
		sut := weather.NewClient(nil, "", "", givenTimeout)

		// WHEN
		_, err := sut.Get(cancelledCtx, "Paris")
		// THEN
		require.Equal(t, err, context.Canceled)
	})
}

func fixAPIKey() string {
	return "api-key-123"
}
