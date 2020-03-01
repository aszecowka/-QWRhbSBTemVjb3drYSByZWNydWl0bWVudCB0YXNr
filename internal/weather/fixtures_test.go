package weather_test

import (
	"context"
	"github.com/aszecowka/QWRhbSBTemVjb3drYSByZWNydWl0bWVudCB0YXNr/internal/weather"
)

func fixCtx() context.Context {
	return context.WithValue(context.Background(), key("key"), "value")
}

type key string

func fixWeatherInParis() *weather.OpenWeatherResponse {
	return &weather.OpenWeatherResponse{
		ID:   1,
		Name: "Paris",
	}
}

func fixWeatherInLondon() *weather.OpenWeatherResponse {
	return &weather.OpenWeatherResponse{
		ID:   2,
		Name: "London",
	}

}
