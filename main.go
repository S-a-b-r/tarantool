package main

import (
	"context"
	"fmt"
	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/pool"
	"time"
)

func main() {
	// Connect to the database
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	dialer := tarantool.NetDialer{
		Address:  "127.0.0.1:3302",
		User:     "sampleuser",
		Password: "123456",
	}
	opts := tarantool.Opts{
		Timeout: time.Second,
	}

	instances := []pool.Instance{{"instance-001", dialer, opts}}

	conn, err := pool.Connect(ctx, instances)

	//conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		fmt.Println("Connection refused:", err)
		return
	}

	// Select by primary key
	data, err := conn.DoInstance(
		tarantool.NewSelectRequest("bands").
			Limit(1).
			Key([]interface{}{uint(1)}), "instance-001",
	).Get()

	if err != nil {
		fmt.Println("Got an error:", err)
	}
	fmt.Println("Tuple selected by the primary key value:", data)
}
