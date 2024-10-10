package main

import (
	"context"

	"github.com/spf13/viper"
	"gitlab.globars.ru/shared/config"
	"gitlab.globars.ru/shared/logger"

	"tarantool/cache"
)

func main() {
	// Connect to the database
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := logger.Init(viper.GetString(config.LogLevel), viper.GetString(config.LogFile), viper.GetInt(config.LogAge), viper.GetInt(config.LogBackups))

	cacheT := cache.Init(ctx, log, config.RedisURL)

	cacheT.Close()
}
