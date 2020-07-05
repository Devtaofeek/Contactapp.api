package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"os"
)

var db *gorm.DB

func init() {
env_load_err := godotenv.Load()
if env_load_err!= nil{
	fmt.Println(env_load_err)
}

	username:=os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)

	conn , err := gorm.Open("postgres",dbUri)

	if err!=nil{
		fmt.Println(err)
	}

	db = conn

}


func GetDB() *gorm.DB {
	return db

}
