package models

type User struct {
	ID int
	Email string `gorm:"type:varchar(100);unique_index"`
	Password string
	Contacts []Contact
}