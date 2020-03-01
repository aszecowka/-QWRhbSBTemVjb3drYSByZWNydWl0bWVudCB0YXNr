package weather_test

import (
	"context"
	"encoding/json"
	"github.com/alicebob/miniredis"
	"github.com/aszecowka/QWRhbSBTemVjb3drYSByZWNydWl0bWVudCB0YXNr/internal/weather"
	"github.com/aszecowka/QWRhbSBTemVjb3drYSByZWNydWl0bWVudCB0YXNr/internal/weather/automock"
	"github.com/go-redis/redis/v7"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	ttl := time.Second
	ctx := context.Background()

	t.Run("returns not found when cache is empty", func(t *testing.T) {
		// GIVEN
		sut, mockRedis := configureCacheWithDeps(t, ttl)
		defer mockRedis.Close()
		// WHEN
		_, err := sut.Get(ctx, "city")
		// THEN
		assert.Equal(t, weather.NotFoundError, err)
	})

	t.Run("store and get element from cache", func(t *testing.T) {
		// GIVEN
		sut, mockRedis := configureCacheWithDeps(t, ttl)
		defer mockRedis.Close()
		// WHEN & THEN
		err := sut.Set(ctx, "London", *fixWeatherInLondon())
		require.NoError(t, err)

		actualResp, err := sut.Get(ctx, "London")
		require.NoError(t, err)
		assert.Equal(t, *fixWeatherInLondon(), actualResp)

		actualID, err := mockRedis.Get("cityNameToID:london")
		require.NoError(t, err)
		assert.Equal(t, "2", actualID)
		actualWeatherInLondon, err := mockRedis.Get("weatherByCityID:2")
		require.NoError(t, err)
		actualWeatherObject := weather.OpenWeatherResponse{}
		err = json.Unmarshal([]byte(actualWeatherInLondon), &actualWeatherObject)
		require.NoError(t, err)
		assert.Equal(t, *fixWeatherInLondon(), actualWeatherObject)
		actualTTL := mockRedis.TTL("weatherByCityID:2")
		assert.Equal(t, ttl, actualTTL)
	})

	t.Run("city name is case-insensitive", func(t *testing.T) {
		// GIVEN
		sut, mockRedis := configureCacheWithDeps(t, ttl)
		defer mockRedis.Close()
		// WHEN & THEN
		err := sut.Set(ctx, "London", *fixWeatherInLondon())
		require.NoError(t, err)

		actualResp, err := sut.Get(ctx, "london")
		require.NoError(t, err)
		assert.Equal(t, *fixWeatherInLondon(), actualResp)
	})

	t.Run("Get: returns error on getting city ID", func(t *testing.T) {
		// GIVEN
		mockRedis := &automock.RedisClient{}

		defer mockRedis.AssertExpectations(t)
		mockRedis.OnGetReturnsError("cityNameToID:london", errors.New("some error"))
		sut := weather.NewCache(mockRedis, ttl)
		// WHEN
		_, err := sut.Get(ctx, "London")
		// THEN
		assert.EqualError(t, err, "while getting ID of the city: some error")
	})

	t.Run("Get: returns error on getting weather", func(t *testing.T) {
		// GIVEN
		mockRedis := &automock.RedisClient{}

		defer mockRedis.AssertExpectations(t)

		mockRedis.On("Get", "cityNameToID:london").Return(redis.NewStringResult("2", nil))
		mockRedis.OnGetReturnsError("weatherByCityID:2", errors.New("some error"))
		sut := weather.NewCache(mockRedis, ttl)
		// WHEN
		_, err := sut.Get(ctx, "London")
		// THEN
		assert.EqualError(t, err, "while reading weather: some error")
	})

	t.Run("Set: returns error on storing city ID", func(t *testing.T) {
		// GIVEN
		mockRedis := &automock.RedisClient{}

		defer mockRedis.AssertExpectations(t)

		mockRedis.OnSetNXReturnsError("cityNameToID:london", errors.New("some error"))
		sut := weather.NewCache(mockRedis, ttl)
		// WHEN
		err := sut.Set(ctx, "London", *fixWeatherInLondon())
		// THEN
		assert.EqualError(t, err, "while storing mapping city to ID: some error")

	})

	t.Run("Set: returns error on storing weather", func(t *testing.T) {
		// GIVEN
		mockRedis := &automock.RedisClient{}

		defer mockRedis.AssertExpectations(t)

		mockRedis.On("SetNX", "cityNameToID:london", mock.Anything, mock.Anything).Return(redis.NewBoolResult(false, nil))
		mockRedis.OnSetReturnsError("weatherByCityID:2", errors.New("some error"))
		sut := weather.NewCache(mockRedis, ttl)
		// WHEN
		err := sut.Set(ctx, "London", *fixWeatherInLondon())
		// THEN
		assert.EqualError(t, err, "while storing weather object: some error")
	})

}

func configureCacheWithDeps(t *testing.T, ttl time.Duration) (weather.Cache, *miniredis.Miniredis) {
	mockRedis, err := miniredis.Run()
	require.NoError(t, err)

	redisCli := redis.NewClient(&redis.Options{
		Addr: mockRedis.Addr(),
	})

	_, err = redisCli.Ping().Result()
	require.NoError(t, err)

	sut := weather.NewCache(redisCli, ttl)
	return sut, mockRedis

}
