package weather

import "time"

type Config struct {
	Port  int `envconfig:"default=8080"`
	Cache struct {
		Address  string        `envconfig:"default=localhost:6379"`
		Password string        `envconfig:"optional"`
		DB       int           `envconfig:"default=0"`
		TTL      time.Duration `envconfig:"default=5m"`
	}
	WeatherAPI struct {
		URL     string `envconfig:"default=https://api.openweathermap.org"`
		Key     string
		Timeout time.Duration `envconfig:"default=1s"`
	}
}
