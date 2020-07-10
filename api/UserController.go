package Controllers

import (
	"encoding/json"
	"fmt"
	. "github.com/devtaofeek/ContactApp.Api/Utils"
	"github.com/devtaofeek/ContactApp.Api/models"
	"net/http"
)

var LoginController = func(w http.ResponseWriter, r *http.Request) {
  var authmodel = &models.Authmodel{}
  err := json.NewDecoder(r.Body).Decode(authmodel)
  if err!=nil{
  	w.WriteHeader(http.StatusBadRequest)
  	Respond(w, Message(false,"invalid request"))
  }
var resp = models.Login(authmodel.Email,authmodel.Password)
Respond(w,resp)
}


var RegisterController = func(w http.ResponseWriter, r *http.Request) {
   user := &models.Authmodel{}
   err := json.NewDecoder(r.Body).Decode(user)
   if err!=nil{
   	w.WriteHeader(http.StatusBadRequest)
   	Respond(w, Message(false,"Invalid request"))
	   return
   }
 var resp = user.CreateAccount()
   Respond(w,resp)
}

var SignOutController = func(w http.ResponseWriter, r *http.Request) {
fmt.Println("sign out hit")
}

var RefreshTokenController = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(int64)
   resp, err := models.RefreshToken(uint(user))
   if err!=nil{
   	w.WriteHeader(http.StatusInternalServerError)
	   Respond(w, Message(false,"Invalid request"))
	   return
   }
   Respond(w,resp)
}

var LogoutController = func(w http.ResponseWriter, r *http.Request) {
	accessuuid := r.Context().Value("accessUuid").(string)
	err  := models.DeleteFromredis(accessuuid)
	if err!=nil{
		Respond(w,Message(false,"Invalid request"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Respond(w,Message(true,"Log out successful"))
}