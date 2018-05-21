package main

import (
  "fmt"
  "github.com/xuanmingyi/astack/models"
)

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
	image := models.Image{ Name: name, Format: format, Checksum: sum, Status: "creating" }
	db.Create(&image)

	dst_file := fmt.Sprintf("%s/%s.%s", global_option.image_dir, sum, format)
	copy_file(file, dst_file)

	image.Status = "available"
	db.Save(&image)
}

func image_list(args []string) {
	db := get_db()
	defer db.Close()

	var images []models.Image
	db.Find(&images)
	headers := []string{"ID", "Name", "Format", "Checksum", "Status"}
	var data [][]string
	for _, v := range images {
		data = append(data, []string{fmt.Sprintf("%d", v.ID), v.Name, v.Format, v.Checksum, v.Status})
	}
	print_table(headers, data)
}

func image_delete(args []string) {
	db := get_db()
	defer db.Close()

	var image models.Image

	db.First(&image, args[0])
	if image != (models.Image{}) {
		db.Delete(&image)
	}else{
		fmt.Println("no image id", args[0])
	}
}
