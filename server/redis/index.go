package rds

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

var Redis redis.Conn

func InitConnection() (c redis.Conn) {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
			fmt.Println("conn redis failed,", err)
			return
	}

	Redis = c

	return 
}