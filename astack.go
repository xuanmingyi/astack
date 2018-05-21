package main

import (
	"os"

	"github.com/xuanmingyi/astack/models"
)

type command_func func([]string)

type Option struct {
	db        string
	image_dir string
	volume_vg string
}

var global_option Option


var command_mappings = map[string]command_func{
	"image-create": image_create,
	"image-list":   image_list,
	"image-delete": image_delete,
	"volume-create": volume_create,
	"volume-list": volume_list,
	"volume-delete": volume_delete,
	"volume-info": volume_info,
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
	db.AutoMigrate(&models.Volume{})

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
