package gimg

import (
	"testing"
)

func TestLoadCfg(t *testing.T) {
	cfgFile := "/Users/kentchen/golang/src/github.com/xingskycn/go-zimg/config.ini"

	cfg, err := LoadConfig(cfgFile)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	} else {
		t.Log(cfg.System.Host)
		//t.Log(cfg.System.Types())
	}
}
