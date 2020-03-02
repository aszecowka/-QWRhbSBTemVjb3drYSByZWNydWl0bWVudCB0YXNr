package main

import (
	"fmt"
	"github.com/aszecowka/QWRhbSBTemVjb3drYSByZWNydWl0bWVudCB0YXNr/internal/weather"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"github.com/vrischmann/envconfig"
	"net/http"
)

func main() {
	cfg := weather.Config{}
	err := envconfig.InitWithPrefix(&cfg, "APP")
	if err != nil {
		panic(err)
	}
	redisCli := redis.NewClient(&redis.Options{
		Addr:     cfg.Cache.Address,
		Password: cfg.Cache.Password,
		DB:       cfg.Cache.DB,
	})

	if _, err := redisCli.Ping().Result(); err != nil {
		panic(err)
	}
	log := logrus.StandardLogger()

	cache := weather.NewCache(redisCli, cfg.Cache.TTL)
	client := weather.NewClient(http.DefaultClient, log, cfg.WeatherAPI.URL, cfg.WeatherAPI.Key, cfg.WeatherAPI.Timeout)
	svc := weather.NewService(cache, client)
	h := weather.NewHandler(svc, log)

	log.Info("Starting server...")
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), h)
	if err != nil {
		panic(err)
	}

}
