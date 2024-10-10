package cache

import (
	"strings"

	"github.com/tarantool/go-tarantool/v2"
)

func GetDealer(url string) tarantool.NetDialer {
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
