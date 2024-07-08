package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

var NewRequestInterval = 10 * time.Millisecond

func main() {
	// Get Redis address from environment variable or use default
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	// Create a Redis connection pool
	pool := &redis.Pool{
		MaxActive:   512,
		MaxIdle:     3,
		IdleTimeout: 10 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisAddr)
		},
		TestOnBorrow: func(c redis.Conn, lastUsed time.Time) error {
			// Get the connection ID
			// id, err := redis.Int(redis.DoWithTimeout(c, 2*time.Millisecond, "CLIENT", "ID"))
			id, err := redis.Int(c.Do("CLIENT", "ID"))
			if err != nil {
				log.Printf("Error getting client ID: %v", err)
			} else {
				fmt.Printf("Connected with client ID: %d\n", id)
			}
			fmt.Printf("connection borrowed: %d, lastUsed: %v\n", id, lastUsed)
			return nil
		},
	}

	// Close the pool when the main function exits
	defer pool.Close()

	// simulate a web request
	for {

		go func() {
			// Get a connection from the pool
			conn := pool.Get()
			defer conn.Close()

			// Read the key
			// _, err := redis.String(redis.DoWithTimeout(conn, 2*time.Millisecond, "GET", "RANDOM_KEY"))
			_, err := redis.String(conn.Do("GET", "RANDOM_KEY"))
			if err != nil && err != redis.ErrNil {
				log.Printf("Error reading key: %v", err)
			}
		}()

		// Wait for the next interval
		time.Sleep(NewRequestInterval)
	}
}
