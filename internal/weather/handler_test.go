package weather_test

import (
	"encoding/json"
	"fmt"
	"github.com/aszecowka/QWRhbSBTemVjb3drYSByZWNydWl0bWVudCB0YXNr/internal/weather"
	"github.com/aszecowka/QWRhbSBTemVjb3drYSByZWNydWl0bWVudCB0YXNr/internal/weather/automock"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	t.Run("got response for many cities", func(t *testing.T) {
		// GIVEN
		mockGetter := &automock.BulkGetter{}
		defer mockGetter.AssertExpectations(t)
		ctx := fixCtx()

		givenResponse := map[string]*weather.OpenWeatherResponse{
			"London":       fixWeatherInLondon(),
			"Paris":        fixWeatherInParis(),
			"MythicalCity": nil,
		}
		mockGetter.On("GetWeatherForCities", ctx, []string{"London", "Paris", "MythicalCity"}).Return(givenResponse, nil)
		sut := weather.NewHandler(mockGetter, nil)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/weather?city=London&city=Paris&city=MythicalCity", nil)
		req = req.WithContext(ctx)
		// WHEN
		sut.ServeHTTP(rr, req)
		// THEN
		assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
		assert.Equal(t, "application/json", rr.Result().Header.Get("Content-Type"))
		require.NotNil(t, rr.Body)
		actual := make(map[string]*weather.OpenWeatherResponse)
		err := json.NewDecoder(rr.Body).Decode(&actual)
		require.NoError(t, err)
		assert.Equal(t, givenResponse, actual)
	})

	t.Run("got client error when no cities specified", func(t *testing.T) {
		// GIVEN
		sut := weather.NewHandler(nil, nil)
		rr := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, "/weather", nil)
		// WHEN
		sut.ServeHTTP(rr, req)
		// THEN
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		require.NotNil(t, rr.Body)
		badRequestResponse := weather.BadRequestResponse{}
		err := json.NewDecoder(rr.Body).Decode(&badRequestResponse)
		require.NoError(t, err)
		assert.Equal(t, "`city` query parameter is required", badRequestResponse.Message)

	})

	t.Run("got client error when wrong HTTP method", func(t *testing.T) {
		// GIVEN
		sut := weather.NewHandler(nil, nil)
		rr := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodPost, "/weather", nil)
		// WHEN
		sut.ServeHTTP(rr, req)
		// THEN
		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("got internal server error when cannot fetch weather", func(t *testing.T) {
		// GIVEN
		mockGetter := &automock.BulkGetter{}
		defer mockGetter.AssertExpectations(t)
		log, logHook := test.NewNullLogger()
		fmt.Println(logHook.Levels())

		mockGetter.On("GetWeatherForCities", mock.Anything, []string{"London"}).Return(nil, errors.New("some error"))
		sut := weather.NewHandler(mockGetter, log)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/weather?city=London", nil)
		// WHEN
		sut.ServeHTTP(rr, req)
		// THEN
		assert.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
		require.Equal(t, 0, rr.Body.Len())
		require.Len(t, logHook.AllEntries(), 1)
		assert.Equal(t, logrus.WarnLevel, logHook.LastEntry().Level)
		assert.Equal(t, "Got error on getting weather for cities: some error", logHook.LastEntry().Message)

	})
}
