package models

import "github.com/jinzhu/gorm"

type Contact struct {
	gorm.Model
	Name string
	PhoneticName string
	NickName string
	JobTitle string
	Company string
	Phone string
	Email string
	Address string
	Relationship string
	Userid uint
}
