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

	config.Init()

	log := logger.Init(viper.GetString(config.LogLevel), viper.GetString(config.LogFile), 1000, 1000)

	log.Info().Msg("starting")
	cacheDB := cache.InitTarantool(ctx, log, "tarantool://admin:admin@localhost:3302")
	defer cacheDB.Close()

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

	// _, err := tdb.Set(cl, "report:just_test", []byte("test string"), 0).Result()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// data, err := tdb.Keys(cl, ".*").Result()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// res, err := tdb.MGet(cl, data...).Result()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// key := fmt.Sprintf("cmd:%s:%s", "testUid", "testMsgId")
	// m := "testtesetetestsetstsetsetestestestsetsets"
	//
	// if err := cacheDB.HSet(cl, key, "command", m).Err(); err != nil {
	// 	log.Error().Err(err).Msg("failed to save command to redis data base")
	// 	return
	// }
	err := cacheDB.Set(cl, "cmd:testUid:testMsgId", []byte("test command"), time.Hour*1).Err()
	if err != nil {
		fmt.Println(err)
	}

	res, err := cacheDB.Get(cl, "cmd:testUid:testMsgId").Result()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)

}
