package redisclipool

import(
	"time"
	"fmt"
	
	"github.com/gomodule/redigo/redis"
	
)

var (
	Clipool       *redis.Pool
	redisServer   string
	redisPassword string
	maxIdle       int
	maxActive     int

	//idleTimeout   int
)

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: (60 * time.Second),
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func InitRedisPool(host, password string, port, idle, active int) {
	redisServer = fmt.Sprintf("%s:%d", host, port)
	redisPassword = password
	maxIdle = idle
	maxActive = active
	
	Clipool = newPool(redisServer, redisPassword)

	return
}

func String(reply interface{}, err1 error) (value string, err2 error) {
	value, err2 = redis.String(reply, err1)
	return
}
