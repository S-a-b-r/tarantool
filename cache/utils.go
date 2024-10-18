package cache

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/crud"
)

func getDialer(url string) tarantool.NetDialer {
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

func getFieldOnHashMap(ret crud.Result, field string) (string, error) {
	rows, ok := ret.Rows.([]interface{})
	if !ok || len(rows) == 0 {
		return "", errors.New("invalid rows data")
	}

	row, ok := rows[0].([]interface{})
	if !ok || len(row) == 0 {
		return "", errors.New("invalid row data")
	}

	data, ok := row[0].(map[interface{}]interface{})
	if !ok {
		return "", errors.New("invalid data type")
	}
	for key, val := range data {
		kStr, ok := key.(string)
		if !ok {
			continue
		}
		if kStr == field {
			return getStringVal(val), nil
		}
	}

	return "", errors.New("not found field")
}

func getValue(ret crud.Result) (string, error) {
	rows, ok := ret.Rows.([]interface{})
	if !ok || len(rows) == 0 {
		return "", errors.New("invalid rows data")
	}

	row, ok := rows[0].([]interface{})
	if !ok || len(row) < 2 {
		return "", errors.New("invalid row data")
	}

	switch d := row[1].(type) {
	case map[interface{}]interface{}:
		return fmt.Sprintf("%+v", d), nil
	case string:
		return d, nil
	case int:
		return strconv.Itoa(d), nil
	}

	return "", errors.New("undefined type")
}

func getStringVal(val interface{}) string {
	switch vv := val.(type) {
	case string:
		return vv
	case int:
		return strconv.Itoa(vv)
	case int8:
		return strconv.Itoa(int(vv))
	case int16:
		return strconv.Itoa(int(vv))
	case int32:
		return strconv.Itoa(int(vv))
	case int64:
		return strconv.Itoa(int(vv))
	case uint8:
		return strconv.Itoa(int(vv))
	case uint16:
		return strconv.Itoa(int(vv))
	case uint32:
		return strconv.Itoa(int(vv))
	case uint64:
		return strconv.Itoa(int(vv))
	}
	return ""
}

func getStringSliceRes(data interface{}) ([]string, error) {
	res := make([]string, 0)

	dSlice, ok := data.([]interface{})
	if !ok {
		return nil, errors.New("data is not array")
	}

	ddSlice, ok := dSlice[0].([]interface{})
	if !ok {
		return nil, errors.New("data is not array")
	}

	for _, val := range ddSlice {
		valStr, ok := val.(string)
		if !ok {
			return []string{}, errors.New("can't set string value")
		}
		res = append(res, valStr)
	}

	return res, nil
}
