package main

import (
	"crypto/sha1"
	"os"
	"io"
	"fmt"
	"errors"

	"github.com/widuu/goini"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/olekukonko/tablewriter"
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


func panic_exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func args_parse(args []string, name string, require_value bool) string {
	for index, value := range args {
		if name == value {
			if require_value {
				if index+1 <= len(args) {
					return args[index+1]
				} else {
					panic_exit(fmt.Sprintf("%s must need a value", name))
				}
			} else {
				return "true"
			}
		}
	}
	return "false"
}

func print_table(headers []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}


func config_parse(config_file string) {
	conf := goini.SetConfig(config_file)
	global_option.db = conf.GetValue("default", "db")
	global_option.image_dir = conf.GetValue("default", "image_dir")
	global_option.volume_vg = conf.GetValue("default", "volume_vg")
}


func get_db() *gorm.DB{
        db, err := gorm.Open("sqlite3", global_option.db)
        if err != nil {
                panic(fmt.Sprintf("%s: %s", err, "failed to connect database"))
        }
	return db
}
