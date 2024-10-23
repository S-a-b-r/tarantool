package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gitlab.globars.ru/shared/cache/v2"
	"gitlab.globars.ru/shared/config"
	"gitlab.globars.ru/shared/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := logger.Init(viper.GetString(config.LogLevel), viper.GetString(config.LogFile), viper.GetInt(config.LogAge), viper.GetInt(config.LogBackups))

	tdb := cache.InitTarantool(ctx, log, "tarantool://admin:admin@localhost:3302")
	defer tdb.Close()

	cl, cancel := context.WithTimeout(ctx, 5*time.Second)

	// report := map[string]string{
	// 	"ID":           "testReport2",
	// 	"Hash":         "hello",
	// 	"AccountId":    "0001",
	// 	"UserID":       "testUser",
	// 	"Created":      time.Now().String(),
	// 	"Requested":    time.Now().String(),
	// 	"TemplateName": "testTmp",
	// }
	//
	// res, err := tdb.Set(cl, "report:002", report, time.Hour*2).Result()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	data, err := tdb.Keys(cl, ".*").Result()
	if err != nil {
		fmt.Println(err)
	}
	res, err := tdb.MGet(cl, data...).Result()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)

}
