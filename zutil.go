package gimg

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	_ "reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

/*
#include <stdlib.h>
#include <stdio.h>
char **EMPTY = NULL;
*/
import "C" //此行和上面的注释之间不能有空行，否则会报错
import "unsafe"

/**
 * @brief str_hash Hash algorithm of processing a md5 string.
 *
 * @param str The md5 string.
 *
 * @return The number less than 1024.
 */
func str_hash(str string) int {
	b := []byte(str)
	c := b[0:3]

	cc := C.CString(string(c))
	defer C.free(unsafe.Pointer(cc))

	d := C.strtol(cc, C.EMPTY, 16)
	d = d / 4
	return int(d)
}

/**
 * @brief is_md5 Check the string is a md5 style.
 *
 * @param s The string.
 *
 * @return 1 for yes and -1 for no.
 */
func is_md5(str string) bool {
	regular := `^([0-9a-zA-Z]){32}$`
	regx := regexp.MustCompile(regular)
	return regx.MatchString(str)
}

/**
 * @brief delete_file delete a file
 *
 * @param path the path of the file
 *
 * @return 1 for OK or -1 for fail
 */
func delete_file(path string) bool {
	os.RemoveAll(path)
	return true
}

/**
 * @brief is_file Check a filename is a file.
 *
 * @param filename The filename input.
 *
 * @return 1 for yes and -1 for no.
 */
func is_file(str string) bool {
	fmt.Println(str)
	f, err := os.Open(str)
	if err != nil {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	if fi.Mode().IsDir() {
		return false
	}
	return true
}

/**
 * @brief is_dir Check a path is a directory.
 *
 * @param path The path input.
 *
 * @return 1 for yes and -1 for no.
 */
func is_dir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	if fi.IsDir() {
		return true
	}

	return false
}

/**
 * @brief is_img Check a file is a image we support(jpg, png, gif).
 *
 * @param filename The name of the file.
 *
 * @return  1 for success and 0 for fail.
 */
func is_img(file_name string) bool {
	return true
}

/**
 * @brief get_type It tell you the type of a file.
 *
 * @param filename The name of the file.
 * @param type Save the type string.
 *
 * @return 1 for success and -1 for fail.
 */
func get_type(file_name string) (string, error) {
	i := strings.LastIndex(file_name, ".")
	if i == -1 {
		return "", fmt.Errorf("FileName [%s] Has No '.' in It.", file_name)
	}
	//fmt.Printf("ext index : %d", i)

	ext := file_name[(i + 1):len(file_name)]
	return ext, nil
}

func is_exist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

/**
 * @brief mk_dir It create a new directory with the path input.
 *
 * @param path The path you want to create.
 *
 * @return  1 for success and -1 for fail.
 */
func mk_dir(path string) bool {
	if is_exist(path) {
		fmt.Printf("Path[%s] is Existed!", path)
		return false
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		if os.IsPermission(err) {
			fmt.Printf("permission is error")
		}
		return false
	}

	return true
}

/**
 * @brief is_special_dir check if the path is a special path
 *
 * @param path the path want to check
 *
 * @return 1 for yes and -1 for not
 */
func is_special_dir(path string) bool {
	if path == "." || path == ".." {
		return true
	}
	return false
}

/**
 * @brief get_cpu_cores get the cpu cores of a server
 *
 * @return the cpu core number
 */
func get_cpu_cores() int {
	return runtime.NumCPU()
}

/**
 * @brief gettid get pid
 *
 * @return pid
 */
func gettid() int {
	return os.Getpid()
}

/**
 * @brief get_file_path get the file's path
 *
 * @param path the file's parent path
 * @param file_name the file name
 * @param file_path the full path of the file
 */
func get_file_path(path string, file_name string) string {
	s := string(path[(len(path) - 1):len(path)])
	if s != "/" {
		return path + "/" + file_name
	} else {
		return path + file_name
	}
}

/**
 * @brief gen_key Generate storage key from md5 and other args.
 *
 * @param key The key string.
 * @param md5 The md5 string.
 * @param argc Count of args.
 * @param ... Args.
 *
 * @return Generate result.
 */
func gen_key(md5 string, args ...interface{}) string {
	s := []string{}
	s = append(s, md5)
	for _, argv := range args {
		switch v := argv.(type) {
		case string:
			s = append(s, v)
		case int:
			s = append(s, strconv.Itoa(v))
		}
	}
	return strings.Join(s, ":")
}

func gen_md5_str(data []byte) string {
	fmt.Println("Begin to Caculate MD5...")
	m := md5.New()
	m.Write(data)
	return hex.EncodeToString(m.Sum(nil))
}

func round(val float64) float64 {
	if val > 0.0 {
		return math.Floor(val + 0.5)
	} else {
		return math.Ceil(val - 0.5)
	}
}
