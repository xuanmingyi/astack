package models

import (
	"github.com/jinzhu/gorm"
)

type Image struct {
	gorm.Model
	Name     string
	Format   string
	Checksum string
	Status   string
}

type Volume struct {
	gorm.Model
	Name string
	Size int64
}
