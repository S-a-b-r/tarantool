package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/crud"
)

func main() {
	// Connect to the database
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	dialer := tarantool.NetDialer{
		Address:  "localhost:3302",
		User:     "admin",
		Password: "admin",
	}
	opts := tarantool.Opts{
		Timeout: time.Second,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		fmt.Println("Connection refused:", err)
		return
	}

	// req := crud.MakeSelectRequest("bands").
	// 	Opts(crud.SelectOpts{
	// 		First: crud.MakeOptInt(2),
	// 	})

	req := crud.MakeGetRequest("bands").Key(4) // getReq

	ret := crud.Result{}
	if err = conn.Do(req).GetTyped(&ret); err != nil {
		fmt.Printf("Failed to execute request: %s", err)
		return
	}

	fmt.Println("Tuple selected by the primary key value:", ret.Rows)
}
