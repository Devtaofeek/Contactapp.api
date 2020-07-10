package Controllers

import (
	"encoding/json"
	. "github.com/devtaofeek/ContactApp.Api/Utils"
	"github.com/devtaofeek/ContactApp.Api/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var GetContactsControllers = func(w http.ResponseWriter, r *http.Request) {


}

var CreateContactController = func(w http.ResponseWriter, r *http.Request) {
var contact = &models.Contact{}
	user := r.Context().Value("user").(int64)
	err := json.NewDecoder(r.Body).Decode(contact)
	contact.Userid = uint(user)

	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		Respond(w, Message(false,"invalid request"))
	}
	var resp = contact.Createcontact()
	Respond(w,resp)
}

var GetContactController = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		Respond(w, Message(false, "There was an error in your request"))
		return
	}
	user := r.Context().Value("user").(int64)
   var resp = models.GetContact(id, int(user))
   Respond(w,resp)
}

var DeleteContactController = func(w http.ResponseWriter, r *http.Request) {

}

var DeleteMultipleContactsController = func(w http.ResponseWriter, r *http.Request) {

}

var UpdateContactController = func(w http.ResponseWriter, r *http.Request) {

}