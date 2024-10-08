package cache

import (
	"context"
	"strings"

	"github.com/rs/zerolog"
	"github.com/tarantool/go-tarantool/v2"
)

type Session struct {
	*tarantool.Connection
	ctx context.Context
	l   *zerolog.Logger
}

func getUserPassword(url string) (user string, password string) {
	userData := strings.Split(url, ":")
	switch len(userData) {
	case 1:
		user = userData[0]
	case 2:
		user = userData[0]
		password = userData[1]
	}
	return user, password
}

func GetDealer(url string) tarantool.NetDialer {
	url = "tarantool://globars:bLtybc84@192.168.169.170:5672/"

	url = strings.Replace(url, "tarantool://", "", -1)
	urlData := strings.Split(url, "@")

	var password, user, address string

	switch len(urlData) {
	case 1:
		address = urlData[0]
	case 2:
		address = urlData[1]
		user, password = getUserPassword(urlData[0])
	}

	return tarantool.NetDialer{
		Address:  address,
		User:     user,
		Password: password,
	}
}

func Init(ctx context.Context, logger *zerolog.Logger, url string) *Session {
	l := logger.With().Str("address", url).Logger()

	dealer := GetDealer(url)

	return &Session{}
}
