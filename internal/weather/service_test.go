package weather_test

import (
	"github.com/aszecowka/QWRhbSBTemVjb3drYSByZWNydWl0bWVudCB0YXNr/internal/weather"
	"github.com/aszecowka/QWRhbSBTemVjb3drYSByZWNydWl0bWVudCB0YXNr/internal/weather/automock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService(t *testing.T) {
	t.Run("returns data for all cities form cache and REST API", func(t *testing.T) {
		// GIVEN
		mockCache := &automock.Cache{}
		defer mockCache.AssertExpectations(t)
		mockCache.On("Get", fixCtx(), "London").Return(*fixWeatherInLondon(), nil)
		mockCache.On("Get", fixCtx(), "Paris").Return(weather.OpenWeatherResponse{}, weather.NotFoundError)
		mockCache.On("Set", fixCtx(), "Paris", *fixWeatherInParis()).Return(nil)

		mockClient := &automock.OpenWeatherClient{}
		mockClient.On("Get", fixCtx(), "Paris").Return(*fixWeatherInParis(), nil)
		defer mockClient.AssertExpectations(t)

		sut := weather.NewService(mockCache, mockClient)
		// WHEN
		actualResponse, err := sut.GetWeatherForCities(fixCtx(), []string{"London", "Paris"})
		// THEN
		require.NoError(t, err)
		expectedResponse := map[string]*weather.OpenWeatherResponse{
			"London": fixWeatherInLondon(),
			"Paris":  fixWeatherInParis(),
		}
		assert.Equal(t, expectedResponse, actualResponse)
	})

	t.Run("if city does not exist, nil weather is returned and weather for the rest cities is returned", func(t *testing.T) {
		// GIVEN
		mockCache := &automock.Cache{}
		defer mockCache.AssertExpectations(t)
		mockCache.On("Get", fixCtx(), "London").Return(*fixWeatherInLondon(), nil)
		mockCache.On("Get", fixCtx(), "ImaginaryCity").Return(weather.OpenWeatherResponse{}, weather.NotFoundError)

		mockClient := &automock.OpenWeatherClient{}
		mockClient.On("Get", fixCtx(), "ImaginaryCity").Return(weather.OpenWeatherResponse{}, weather.NotFoundError)
		defer mockClient.AssertExpectations(t)

		sut := weather.NewService(mockCache, mockClient)
		// WHEN
		actualResponse, err := sut.GetWeatherForCities(fixCtx(), []string{"London", "ImaginaryCity"})
		// THEN
		require.NoError(t, err)
		expectedResponse := map[string]*weather.OpenWeatherResponse{
			"London":        fixWeatherInLondon(),
			"ImaginaryCity": nil,
		}
		assert.Equal(t, expectedResponse, actualResponse)
	})

	t.Run("returns error if cannot access cache", func(t *testing.T) {
		// GIVEN
		mockCache := &automock.Cache{}
		defer mockCache.AssertExpectations(t)
		mockCache.On("Get", fixCtx(), "London").Return(weather.OpenWeatherResponse{}, errors.New("some error"))

		sut := weather.NewService(mockCache, nil)
		// WHEN
		_, err := sut.GetWeatherForCities(fixCtx(), []string{"London", "Paris"})
		// THEN
		require.EqualError(t, err, "while fetching weather for city London: while getting weather from cache: some error")
	})

	t.Run("returns error if cannot access REST API", func(t *testing.T) {
		// GIVEN
		mockCache := &automock.Cache{}
		defer mockCache.AssertExpectations(t)
		mockCache.On("Get", fixCtx(), "London").Return(weather.OpenWeatherResponse{}, weather.NotFoundError)

		mockClient := &automock.OpenWeatherClient{}
		mockClient.On("Get", fixCtx(), "London").Return(weather.OpenWeatherResponse{}, errors.New("some error"))
		defer mockClient.AssertExpectations(t)

		sut := weather.NewService(mockCache, mockClient)
		// WHEN
		_, err := sut.GetWeatherForCities(fixCtx(), []string{"London", "Paris"})
		// THEN
		require.EqualError(t, err, "while fetching weather for city London: while getting weather from the REST service: some error")

	})

	t.Run("returns error if cannot save data in the cache", func(t *testing.T) {
		// GIVEN
		mockCache := &automock.Cache{}
		defer mockCache.AssertExpectations(t)
		mockCache.On("Get", fixCtx(), "London").Return(weather.OpenWeatherResponse{}, weather.NotFoundError)
		mockCache.On("Set", fixCtx(), "London", *fixWeatherInLondon()).Return(errors.New("some error"))

		mockClient := &automock.OpenWeatherClient{}
		mockClient.On("Get", fixCtx(), "London").Return(*fixWeatherInLondon(), nil)
		defer mockClient.AssertExpectations(t)

		sut := weather.NewService(mockCache, mockClient)
		// WHEN
		_, err := sut.GetWeatherForCities(fixCtx(), []string{"London", "Paris"})
		// THEN
		require.EqualError(t, err, "while fetching weather for city London: while updating cache: some error")

	})
}
