package automock

import (
	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/mock"
)

func (_m *RedisClient) OnGetReturnsError(expectedKey string, returnedError error) {
	cmd := redis.StringCmd{}
	cmd.SetErr(returnedError)
	_m.On("Get", expectedKey).Return(&cmd).Once()
}

func (_m *RedisClient) OnSetReturnsError(key string, returnedError error) {
	cmd := redis.StatusCmd{}
	cmd.SetErr(returnedError)
	_m.On("Set", key, mock.Anything, mock.Anything).Return(&cmd).Once()
}

//SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd
func (_m *RedisClient) OnSetNXReturnsError(key string, returnedError error) {
	cmd := redis.BoolCmd{}
	cmd.SetErr(returnedError)
	_m.On("SetNX", key, mock.Anything, mock.Anything).Return(&cmd).Once()
}
