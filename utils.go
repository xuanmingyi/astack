package main

import (
	"crypto/sha1"
	"os"
	"io"
	"fmt"
	"errors"
)


func find_first_config_file() (string, error){
	config_files := []string {"/etc/astack/config.ini", "/etc/astack/astack.ini", "astack.ini"}
	for _, value := range config_files {
		_, err := os.Stat(value)
		if err == nil {
			return value, nil
		}
	}
	return "", errors.New("config not found")
}

func sha1_sum(filename string) (string, error) {
	file, err1 := os.Open(filename)
	if err1 != nil {
		return "", err1
	}
	defer file.Close()

	sum := sha1.New()

	_, err2 := io.Copy(sum, file)
	if err2 != nil {
		return "", err2
	}
	return fmt.Sprintf("%x", sum.Sum(nil)), nil
}

func copy_file(src_file, dst_file string) (written int64, err error) {
	src, err := os.Open(src_file)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.OpenFile(dst_file, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}
