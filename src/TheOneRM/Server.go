package main

import (
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

var (
	cacheKey      = "onlineClient"
	redisServer   = flag.String("redisServer", "127.0.0.1:6379", "")
	redisPassword = flag.String("redisPassword", "", "")
	userId        = flag.String("userId", "", "")
	kickCmd       = flag.String("kick", "", "")
	clearCmd      = flag.String("clear", "", "")
)

type redisClient struct {
	pool *redis.Pool
}

//关闭redis连接池
func (c *redisClient) Close() {
	if c.pool != nil {
		c.pool.Close()
	}
}

//获取redis链接
func newClient(server, password string) *redisClient {
	return &redisClient{
		pool: newPool(server, password),
	}
}

//创建redis connection pool
func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}

			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func main() {
	flag.Parse()

	client := newClient(*redisServer, *redisPassword)
	defer client.Close()

	fmt.Println("...")
}
