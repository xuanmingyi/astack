package main

import (
	"fmt"
	"github.com/widuu/goini"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/xuanmingyi/astack/models"

//	"github.com/olekukonko/tablewriter"

)

var base_url = "http://localhost:9090"

type command_func func([]string)

type Option struct {
	db        string
	image_dir string
}

var global_option Option

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

func config_parse(config_file string) {
	conf := goini.SetConfig(config_file)
	global_option.db = conf.GetValue("default", "db")
	global_option.image_dir = conf.GetValue("default", "image_dir")
}


func get_db() *gorm.DB{
        db, err := gorm.Open("sqlite3", global_option.db)
        if err != nil {
                panic(fmt.Sprintf("%s: %s", err, "failed to connect database"))
        }
	return db
}


func image_create(args []string) {
	//astack image-create --name cirros --file cirros.img --format qcow2
	file := args_parse(args, "--file", true)
	name := args_parse(args, "--name", true)
	format := args_parse(args, "--format", true)

	db := get_db()
	defer db.Close()

	// sha1
	sum, err1 := sha1_sum(file)
	if err1 != nil {
		panic_exit(err1.Error())
	}
	db.Create(&models.Image{ Name: name, Format: format, Checksum: sum, Status: "creating" })
}


func image_list(args []string) {
	db := get_db()
	defer db.Close()

	images := db.Where("status <> ?", "deleted").Find(&models.Image)
}

var command_mappings = map[string]command_func{
	"image-create": image_create,
	"image-list":   image_list,
}

func _init() {
	config_file, err1 := find_first_config_file()
	if err1 != nil {
		panic_exit(err1.Error())
	}
	config_parse(config_file)

	db := get_db()
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&models.Image{})

	err2 := os.MkdirAll(global_option.image_dir, os.ModePerm)
	if err2 != nil {
		panic_exit(err2.Error())
	}


}

func main() {
	_init()

	if len(os.Args) == 1 {
		panic_exit("Error os.Args = 1")
	}
	if _, ok := command_mappings[os.Args[1]]; ok {
		command := os.Args[1]
		command_mappings[command](os.Args[2:])
	} else {
		panic_exit("Error command not found")
	}
}
