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
fmt.Println("refrsh token hit")
}