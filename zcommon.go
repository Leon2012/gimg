package gimg

import (
	_ "encoding/json"
	_ "fmt"
)

const (
	PROJECT_VERSION = "1.0"
	MAX_LINE        = 1024
	CACHE_KEY_SIZE  = 128
	RETRY_TIME_WAIT = 1000
	CACHE_MAX_SIZE  = 1048576 //1024*1024
	PATH_MAX_SIZE   = 512
)

type ZRequest struct {
	Md5        string
	ImageType  string
	Width      int
	Height     int
	Proportion int
	Gary       int
	X          int
	Y          int
	Rotate     int
	Format     string
	Save       int
	Quality    int
}

type ZImageInfo struct {
	Size    int    `json:"size"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Quality int    `json:"quality"`
	Format  string `json:"format"`
}

type ZContext struct {
	Config AppConfig
	Logger *ZLogger
	Cache  *ZCache
	Image  *ZImage
}

func NewContext(conf AppConfig, log *ZLogger, c *ZCache, i *ZImage) *ZContext {
	return &ZContext{Config: conf,
		Logger: log,
		Cache:  c,
		Image:  i}
}
