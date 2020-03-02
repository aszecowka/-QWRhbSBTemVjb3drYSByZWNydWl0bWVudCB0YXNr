package weather

import (
	"context"
	"errors"
	"fmt"
)

//go:generate mockery -name=Cache -output=automock -outpkg=automock -case=underscore
type Cache interface {
	Get(ctx context.Context, city string) (OpenWeatherResponse, error)
	Set(ctx context.Context, queriedCity string, weather OpenWeatherResponse) error
}

//go:generate mockery -name=OpenWeatherClient -output=automock -outpkg=automock -case=underscore
type OpenWeatherClient interface {
	Get(ctx context.Context, city string) (OpenWeatherResponse, error)
}

func NewService(cache Cache, client OpenWeatherClient) *service {
	return &service{cache: cache, client: client}
}

type service struct {
	cache  Cache
	client OpenWeatherClient
}

func (s *service) GetWeatherForCities(ctx context.Context, cities []string) (map[string]*OpenWeatherResponse, error) {
	out := make(map[string]*OpenWeatherResponse)
	for _, city := range cities {
		resp, err := s.getWeatherForCity(ctx, city)
		switch {
		case errors.Is(err, NotFoundError):
			out[city] = nil
		case err != nil:
			return nil, fmt.Errorf("while fetching weather for city %s: %w", city, err)
		default:
			out[city] = &resp
		}

	}
	return out, nil
}

func (s *service) getWeatherForCity(ctx context.Context, city string) (OpenWeatherResponse, error) {
	w, err := s.cache.Get(ctx, city)
	switch {
	case errors.Is(err, NotFoundError):
	case err != nil:
		return OpenWeatherResponse{}, fmt.Errorf("while getting weather from cache: %w", err)
	default:
		return w, nil
	}
	w, err = s.client.Get(ctx, city)
	if err != nil {
		return OpenWeatherResponse{}, fmt.Errorf("while getting weather from the REST service: %w", err)
	}

	if err := s.cache.Set(ctx, city, w); err != nil {
		return OpenWeatherResponse{}, fmt.Errorf("while updating cache: %w", err)
	}

	return w, nil
}
