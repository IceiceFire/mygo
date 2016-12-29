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

//获取所有key
func (c *redisClient) getAllkey() {
	conn := c.pool.Get()
	keys, err := redis.Strings(conn.Do("KEYS", "*"))

	if err != nil {
		panic(err)
	}
	fmt.Printf("keys : %d \n", len(keys))
	for no, key := range keys {
		fmt.Println("No:", no, "--keys:", key)
	}
}

//获取所有key
func (c *redisClient) getinfo() {
	conn := c.pool.Get()

	r, err := redis.String(conn.Do("CLIENT", "LIST"))

	if err != nil {
		panic(err)
	}
	fmt.Println(r)

}

//获取所有在线用户
func (c *redisClient) GetAll() {
	conn := c.pool.Get()
	clients, err := redis.StringMap(conn.Do("HGETALL", cacheKey))
	if err != nil {
		panic(err)
	}
	fmt.Printf("online client: %d \n", len(clients))
	for uId, client := range clients {
		fmt.Printf("%s -- %s\n", uId, client)
	}
}

//根据用户ID获取单个用户
func (c *redisClient) GetOne(id string) {
	client, err := redis.String(c.pool.Get().Do("HGET", cacheKey, id))

	if err != nil {
		panic(err)
	}
	fmt.Println(client)
}

//踢出某个用户
func (c *redisClient) Kick(id string) {
	result, err := c.pool.Get().Do("HDEL", cacheKey, id)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

//清除所有在线用户信息
func (c *redisClient) ClearAll() {
	result, err := c.pool.Get().Do("DEL", cacheKey)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

//关闭redis连接池
func (c *redisClient) Close() {
	if c.pool != nil {
		c.pool.Close()
	}
}

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

	client.getinfo()

	if *kickCmd != "" {
		client.Kick(*userId)
	}

	if *clearCmd == "all" {
		client.ClearAll()
	}

	if *userId == "" {
		client.GetAll()
	} else {
		client.GetOne(*userId)
	}
}
