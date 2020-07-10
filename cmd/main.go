package main

import (
	"fmt"
	Controllers "github.com/devtaofeek/ContactApp.Api/api"
	"github.com/devtaofeek/ContactApp.Api/app"
	"github.com/devtaofeek/ContactApp.Api/database"
	"github.com/devtaofeek/ContactApp.Api/models"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)
var client redis.Client
func main() {
	var r = mux.NewRouter()
	r.Use(app.JwtAuth)
	port := os.Getenv("PORT")
	r.HandleFunc("/api/user/login",Controllers.LoginController).Methods("POST")
	r.HandleFunc("/api/user/register", Controllers.RegisterController).Methods("POST")
	r.HandleFunc("/api/user/signout",Controllers.SignOutController).Methods("POST")
	r.HandleFunc("/api/user/refreshtoken", Controllers.RefreshTokenController).Methods("POST")
	r.HandleFunc("/api/user/logout",Controllers.LogoutController).Methods("POST")
	r.HandleFunc("/api/user/GetContacts",Controllers.GetContactsControllers).Methods("POST")

	if port == "" {
		port = "8000" //localhost
	}
	//CreateDb()
	fmt.Println(port)
	log.Fatal(http.ListenAndServe(":"+port, r))


}

func CreateDb()  {
	err:= database.GetDB().DropTableIfExists(models.User{},models.Contact{})
	database.GetDB().AutoMigrate(models.User{},models.Contact{})
	database.GetDB().Model(models.Contact{}).AddForeignKey("userid","users(id)","CASCADE","CASCADE")
	fmt.Println(err)
	
// the only reason why am using redis is to forcefully invalidate the token when a user logs in before expiry
}




