package main

import (
	"context"
	"fmt"
	"github.com/tarantool/go-tarantool/v2"
	"time"
)

func main() {
	// Connect to the database
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	dialer := tarantool.NetDialer{
		Address: "127.0.0.1:3301",
		User:    "guest",
	}
	opts := tarantool.Opts{
		Timeout: time.Second,
	}
	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		fmt.Println("Connection refused:", err)
		return
	}

	// Select by primary key
	data, err := conn.Do(
		tarantool.NewSelectRequest("users").
			Limit(10).
			Iterator(tarantool.IterEq).
			Key([]interface{}{uint(1)}),
	).Get()

	if err != nil {
		fmt.Println("Got an error:", err)
	}
	fmt.Println("Tuple selected by the primary key value:", data)
}
