package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tarantool/go-tarantool/v2"
)

func main() {
	// Connect to the database
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	dialer := tarantool.NetDialer{
		Address:  "localhost:3302",
		User:     "sampleuser",
		Password: "123456",
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
		tarantool.NewSelectRequest("bands").
			Limit(10).
			Iterator(tarantool.IterEq).
			Key([]interface{}{uint(1)}),
	).Get()
	if err != nil {
		fmt.Println("Got an error:", err)
	}
	fmt.Println("Tuple selected by the primary key value:", data)

	// Insert data
	// tuples := [][]interface{}{
	// 	{1, "Roxette", 1986},
	// 	{2, "Scorpions", 1965},
	// 	{3, "Ace of Base", 1987},
	// 	{4, "The Beatles", 1960},
	// }
	// var futures []*tarantool.Future
	// for _, tuple := range tuples {
	// 	request := tarantool.NewInsertRequest("bands").Tuple(tuple)
	// 	futures = append(futures, conn.Do(request))
	// }
	//
	// fmt.Println("Inserted tuples:")
	// for _, future := range futures {
	// 	result, err := future.Get()
	// 	if err != nil {
	// 		fmt.Println("Got an error:", err)
	// 	} else {
	// 		fmt.Println(result)
	// 	}
	// }

}
