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

type storeReport struct {
	ID           string    `json:"id"`
	Hash         string    `json:"hash"`
	AccountId    string    `json:"accountId"`
	UserID       string    `json:"userId"`
	Created      time.Time `json:"created"`
	Requested    time.Time `json:"requested"`
	TemplateName string    `json:"templateName"`
	Status       string    `json:"status"`
	Path         string    `json:"path"`
	Error        string    `json:"error,omitempty"`
	ObjectName   string    `json:"objectName,omitempty"`
	Report       string    `json:"report"`
	IsNoTrack    bool      `json:"isNoTrack"`
	IsMobile     bool      `json:"isMobile"`
}

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

	// report := map[string]string{
	// 	"ID":           "testRep",
	// 	"Hash":         "hello",
	// 	"AccountId":    "0001",
	// 	"UserID":       "testUser",
	// 	"Created":      time.Now().String(),
	// 	"Requested":    time.Now().String(),
	// 	"TemplateName": "testTmp",
	// }

	// data, err := cacheT.HGet(cl, "report:001", "ID").Result()
	// data, err := cacheT.Get(cl, "report:001").Result()
	data, err := cacheT.Set(cl, "report:001", report, time.Hour*2).Result()

	if err != nil {
		fmt.Println("err: ", err)
	}
	fmt.Println(data)

	// cacheT.Subscriber(10, cache.PChKeyEventsHSet, func(str string) {
	// 	fmt.Println(str)
	// })

	cacheT.Close()
}
