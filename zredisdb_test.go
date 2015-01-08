package gimg

import (
	"testing"
)

func TestConnect(t *testing.T) {
	s := "127.0.0.1"
	p := 8888

	db, err := NewRedisDB(s, p)
	if err != nil {
		t.Fail()
	} else {
		t.Log("connect redis db")

		if db.Exist("name") {
			t.Log("name is exists!!!")
		} else {
			t.Log("name isnot exists!!!")
		}

	}
}

func TestGetRedis(t *testing.T) {
	s := "127.0.0.1"
	p := 8888
	key := "e351d9d8c9f409ec0ff4d518d6f7551f"

	db, err := NewRedisDB(s, p)
	if err != nil {
		t.Fail()
	} else {
		t.Log("connect redis db")
		data, err := db.Get(key)
		if err != nil {
			t.Fail()
		} else {
			t.Log(data)
		}

	}
}
