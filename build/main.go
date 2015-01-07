package main

import (
	"flag"
	"fmt"
	_ "github.com/bradfitz/gomemcache/memcache"
	zimg "github.com/leon2012/gimg"
	"net/http"
	"os"
	// "os/signal"
	// "syscall"
)

var cfgFile string
var zContext *zimg.ZContext

func main() {
	configPtr := flag.String("config", "", "config file")
	flag.Usage = usage
	flag.Parse()

	if *configPtr == "" {
		*configPtr = "./conf/config.ini"
	}

	cfgFile = *configPtr

	isExist, _ := exists(cfgFile)
	if !isExist {
		fmt.Println("config file not exist!")
		os.Exit(-1)
	}

	/*===================加载配置文件========================*/
	cfg, err := zimg.LoadConfig(cfgFile)
	checkError(err)

	/*===================加载log========================*/
	//log, err := logger.NewFileLogger("zimg", 0, cfg.System.LogName)
	log, err := zimg.NewLogger("zimg", 0)
	checkError(err)
	defer log.Close()

	/*===================加载memcache cache========================*/
	var cache *zimg.ZCache
	if cfg.Cache.Cache == 1 {
		// var client *memcache.Client
		// cacheAddr := fmt.Sprintf("%s:%d", cfg.Cache.MemcacheHost, cfg.Cache.MemcachePort)
		// client = memcache.New(cacheAddr)
		cache = zimg.NewCache(cfg.Cache.MemcacheHost, cfg.Cache.MemcachePort)
	} else {
		cache = nil
	}

	/*===================加载image handler ========================*/
	img := zimg.NewImage()

	zContext = zimg.NewContext(cfg, log, cache, img)
	//zContext.Logger.Info("load config.ini success!")

	addr := fmt.Sprintf("%s:%d", zContext.Config.System.Host, zContext.Config.System.Port)
	zContext.Logger.Info("server start run :  %s", addr)

	zHttpd := zimg.NewHttpd(zContext)
	err = http.ListenAndServe(addr, zHttpd)
	if err != nil {
		zContext.Logger.Error("error : %s", err.Error())
	}

	//signalHandle()
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:--config=/etc/config.ini \n")
	flag.PrintDefaults()
	os.Exit(-2)
}

func checkError(err error) {
	if err != nil {
		panic(err)
		os.Exit(-2)
	}
}

// func signalHandle() {
// 	// Handle SIGINT and SIGTERM.
// 	ch := make(chan os.Signal, 1)
// 	over := make(chan bool, 1)
// 	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

// 	go func() {
// 		sig := <-ch
// 		//zContext.Logger.Info(sig)
// 		over <- true
// 	}()

// 	//zContext.Logger.Info(<-over)

// 	zContext.Logger.Info("server stop!!!")
// }
