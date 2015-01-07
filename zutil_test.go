package gimg

import (
	"fmt"
	"strings"
	"testing"
)

func TestStrHash(t *testing.T) {
	md5Sum := "814de890c0d588060a9390dcf331f4ed"
	t.Logf("result : %d", str_hash(md5Sum))
	t.Logf("result1 : %d", str_hash(string(md5Sum[3:])))
}

func TestIsMd5(t *testing.T) {
	md5Sum := "814de890c0d588060a9390dcf331f4ed"
	if is_md5(md5Sum) {
		t.Log("success")
	} else {
		t.Fail()
	}
}

func TestIsFile(t *testing.T) {
	f := "/Users/kentchen/Desktop/time_mshutdown.txt"
	//f := "/Users/kentchen/Desktop"
	//f := "/Users/kentchen/Desktop/Windows 8.1"

	if is_file(f) {
		t.Log("success")
	} else {
		t.Fail()
	}
}

func TestGetType(t *testing.T) {
	f := "/Users/kentchen/Desktop/time_mshutdown.txt.info"
	i := strings.LastIndex(f, ".")

	fmt.Printf("index : %d, len : %d, ext len: %d \n", i, len(f), (len(f) - i))

	//fmt.Printf("ext : %s", f[39:42])

	ext, err := get_type(f)
	if err != nil {
		t.Fatalf("%s", err.Error())
	} else {
		t.Logf("type : %s", ext)
	}
}

func TestGenKey(t *testing.T) {
	// s := []string{}
	// s = append(s, "814de890c0d588060a9390dcf331f4ed")
	// s = append(s, "12")

	// t.Logf("%s", s)
	t.Logf("%s", gen_key("814de890c0d588060a9390dcf331f4ed", 508, "300", "300"))
}

func TestRound(t *testing.T) {
	d := 12345678.6
	t.Logf("%ld", round(d))
}
