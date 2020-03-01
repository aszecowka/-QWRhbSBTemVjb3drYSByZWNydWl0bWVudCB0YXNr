package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"strings"
	"time"
)

//go:generate mockery -name=RedisClient -output=automock -outpkg=automock -case=underscore
type RedisClient interface {
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd
}

type redisCache struct {
	ttl      time.Duration
	redisCli RedisClient
}

func NewCache(redisCli RedisClient, ttl time.Duration) *redisCache {
	return &redisCache{
		redisCli: redisCli,
		ttl:      ttl,
	}
}

func (c *redisCache) Get(ctx context.Context, city string) (OpenWeatherResponse, error) {
	cityID, err := c.getCityID(city)
	if err != nil {
		return OpenWeatherResponse{}, err
	}
	return c.getWeather(cityID)
}

func (c *redisCache) Set(ctx context.Context, queriedCity string, weather OpenWeatherResponse) error {
	if err := c.storeCityID(weather.ID, queriedCity); err != nil {
		return err
	}
	if err := c.storeWeather(weather); err != nil {
		return err
	}

	return nil
}

func (c *redisCache) getCityID(city string) (int, error) {
	id, err := c.redisCli.Get(c.keyForCityNameToID(city)).Int()
	switch {
	case err == redis.Nil:
		return -1, NotFoundError
	case err != nil:
		return -1, fmt.Errorf("while getting ID of the city: %w", err)
	}
	return id, nil
}

func (c *redisCache) getWeather(cityID int) (OpenWeatherResponse, error) {
	asJSON, err := c.redisCli.Get(c.keyForWeather(cityID)).Result()
	switch {
	case err == redis.Nil:
		return OpenWeatherResponse{}, NotFoundError
	case err != nil:
		return OpenWeatherResponse{}, fmt.Errorf("while reading weather: %w", err)
	}

	out := OpenWeatherResponse{}
	if err := json.Unmarshal([]byte(asJSON), &out); err != nil {
		return OpenWeatherResponse{}, fmt.Errorf("while unmarshaling weather object: %w", err)
	}
	return out, nil
}

func (c *redisCache) storeCityID(cityID int, queriedCity string) error {
	if _, err := c.redisCli.SetNX(c.keyForCityNameToID(queriedCity), cityID, 0).Result(); err != nil {
		return fmt.Errorf("while storing mapping city to ID: %w", err)
	}
	return nil
}

func (c *redisCache) storeWeather(weather OpenWeatherResponse) error {
	b, err := json.Marshal(weather)
	if err != nil {
		return fmt.Errorf("while marshaling weather object: %w", err)
	}

	_, err = c.redisCli.Set(c.keyForWeather(weather.ID), string(b), c.ttl).Result()
	if err != nil {
		return fmt.Errorf("while storing weather object: %w", err)
	}
	return nil
}

func (c *redisCache) keyForWeather(cityID int) string {
	return fmt.Sprintf("weatherByCityID:%d", cityID)
}

func (c *redisCache) keyForCityNameToID(cityName string) string {
	return fmt.Sprintf("cityNameToID:%s", strings.ToLower(cityName))
}

var NotFoundError = errors.New("not found")
