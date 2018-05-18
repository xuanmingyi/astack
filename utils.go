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
