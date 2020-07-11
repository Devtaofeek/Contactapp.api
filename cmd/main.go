package main

import (
	"fmt"
	Controllers "github.com/devtaofeek/ContactApp.Api/api"
	"github.com/devtaofeek/ContactApp.Api/app"
	"github.com/devtaofeek/ContactApp.Api/database"
	"github.com/devtaofeek/ContactApp.Api/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)
func main() {
	var r = mux.NewRouter()
	r.Use(app.JwtAuth)
	port := os.Getenv("PORT")
	r.HandleFunc("/api/user/login",Controllers.LoginController).Methods("POST")
	r.HandleFunc("/api/user/register", Controllers.RegisterController).Methods("POST")
	r.HandleFunc("/api/user/signout",Controllers.SignOutController).Methods("POST")
	r.HandleFunc("/api/user/refreshtoken", Controllers.RefreshTokenController).Methods("POST")
	r.HandleFunc("/api/user/logout",Controllers.LogoutController).Methods("POST")
	r.HandleFunc("/api/user/GetContact/{id}",Controllers.GetContactController).Methods("GET")
	r.HandleFunc("/api/user/GetContact",Controllers.GetContactsControllers).Methods("GET")
	r.HandleFunc("/api/user/DeleteContact/{id}",Controllers.DeleteContactController).Methods("DELETE")
	r.HandleFunc("/api/user/UpdateContact",Controllers.UpdateContactController).Methods("PUT")

	r.HandleFunc("/api/user/createcontact",Controllers.CreateContactController).Methods("POST")

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
	
}




