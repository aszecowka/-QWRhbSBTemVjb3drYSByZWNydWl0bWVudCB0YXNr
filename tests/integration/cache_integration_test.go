// +build integration

package integration_test

import (
	"context"
	"github.com/aszecowka/QWRhbSBTemVjb3drYSByZWNydWl0bWVudCB0YXNr/internal/weather"
	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestCacheIntegration(t *testing.T) {
	// GIVEN
	address := os.Getenv("REDIS_ADDRESS")
	require.NotEmpty(t, address)
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       13,
	})

	_, err := client.Ping().Result()
	require.NoError(t, err)

	defer func() {
		client.FlushDB()
	}()

	givenTTL := time.Millisecond*10
	cache := weather.NewCache(client, givenTTL)
	ctx := context.Background()
	cityName := "randomCityName"
	givenWeather := weather.OpenWeatherResponse{ID: 1, Name: cityName, Main: weather.Main{Temp: 55.55}}

	// WHEN & THEN
	err = cache.Set(ctx, cityName, givenWeather)
	require.NoError(t, err)

	actual, err := cache.Get(ctx, cityName)
	require.NoError(t, err)

	assert.Equal(t, givenWeather, actual)

	<-time.After(givenTTL + time.Millisecond*10)

	_, err = cache.Get(ctx, cityName)
	assert.Equal(t, weather.NotFoundError, err)
}
