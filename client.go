package main

import ( 
		"github.com/go-redis/redis"
		"fmt"
		)

var client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

func main() {

	client.Set("key", "value", 0)

	val, err := client.Get("key").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("key", val)
}

