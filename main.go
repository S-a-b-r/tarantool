package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gitlab.globars.ru/shared/config"
	"gitlab.globars.ru/shared/logger"

	"tarantool/cache"
)

func main() {
	config.Init()

	// Connect to the database
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := logger.Init(viper.GetString(config.LogLevel), viper.GetString(config.LogFile), viper.GetInt(config.LogAge), viper.GetInt(config.LogBackups))
	log.Info().Msg("starting...")

	cacheT, err := cache.Init(ctx, log, viper.GetString(config.RedisURL))
	if err != nil {
		panic(err)
	}

	cl, _ := context.WithTimeout(ctx, time.Second*1)
	data, err := cacheT.HGet(cl, "unitInfo:01", "uid").Result()
	if err != nil {
		fmt.Println("err: ", err)
	}
	fmt.Println(data)

	// cacheT.Subscriber(10, cache.PChKeyEventsHSet, func(str string) {
	// 	fmt.Println(str)
	// })

	cacheT.Close()
}
