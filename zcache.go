package gimg

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
)

type ZCache struct {
	host   string
	port   int
	client *memcache.Client
}

// func NewCache(c *memcache.Client) *ZCache {
// 	return &ZCache{client: c}
// }

func NewCache(h string, p int) *ZCache {
	var c *memcache.Client
	cacheAddr := fmt.Sprintf("%s:%d", h, p)
	c = memcache.New(cacheAddr)
	return &ZCache{client: c, host: h, port: p}
}

func (z *ZCache) ReTry() {
	cacheAddr := fmt.Sprintf("%s:%d", z.host, z.port)
	c := memcache.New(cacheAddr)
	z.client = c
}

func (z *ZCache) FindCache(key string) (string, error) {
	it, err := z.client.Get(key)
	if err != nil {
		return "", err
	}
	return string(it.Value), nil
}

func (z *ZCache) SetCache(k string, v string) error {
	it := &memcache.Item{Key: k, Value: []byte(v)}
	return z.client.Set(it)
}

func (z *ZCache) Exist(key string) bool {
	_, err := z.FindCache(key)
	if err != nil {
		return false
	}
	return true
}

func (z *ZCache) FindCacheBin(key string) ([]byte, error) {
	it, err := z.client.Get(key)
	if err != nil {
		return nil, err
	} else {
		return it.Value, nil
	}

}

func (z *ZCache) SetCacheBin(k string, v []byte) error {
	it := &memcache.Item{Key: k, Value: v}
	return z.client.Set(it)
}

func (z *ZCache) DelCache(key string) error {
	return z.client.Delete(key)
}
