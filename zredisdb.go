package gimg

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

type ZRedisDB struct {
	server    string
	port      int
	pool      *redis.Pool
	isConnect bool
}

func NewRedisDB(s string, p int) (*ZRedisDB, error) {
	addr := fmt.Sprintf("%s:%d", s, p)
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return &ZRedisDB{
		server:    s,
		port:      p,
		pool:      pool,
		isConnect: false,
	}, nil
}

func (z *ZRedisDB) getConnect() (redis.Conn, error) {
	if z.isConnect {
		conn := z.pool.Get()
		return conn, nil
	} else {
		return nil, errors.New("Can not connect db")
	}
}

func (z *ZRedisDB) Exist(key string) bool {
	conn, err := z.getConnect()
	defer conn.Close()
	if err != nil {
		return false
	}

	isExists, _ := redis.Bool(conn.Do("EXISTS", key))
	return isExists
}

func (z *ZRedisDB) Do(commandName string, args ...interface{}) (interface{}, error) {
	conn, err := z.getConnect()
	defer conn.Close()
	if err != nil {
		return nil, errors.New("Can not connect db!")
	}
	return conn.Do(commandName, args...)
}

func (z *ZRedisDB) Send(commandName string, args ...interface{}) error {
	conn, err := z.getConnect()
	defer conn.Close()
	if err != nil {
		return errors.New("Can not connect db!")
	}
	return conn.Send(commandName, args...)
}

func (z *ZRedisDB) Flush() {
	if z.isConnect {
		conn := z.pool.Get()
		defer conn.Close()
		conn.Flush()
	}
}

func (z *ZRedisDB) Close() {
	if z.isConnect {
		z.pool.Close()
	}
}
