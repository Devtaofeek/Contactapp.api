package models

import (
	. "github.com/devtaofeek/ContactApp.Api/Utils"
	"github.com/devtaofeek/ContactApp.Api/database"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type User struct {
	gorm.Model
	Email string `gorm:"type:varchar(100);unique_index"`
	Password string
	Contacts []Contact `gorm:"foreignkey:Userid"`
}

type Authmodel struct {
	Email    string
	Password string
}

func (authmodel *Authmodel) Validate() (map[string]interface{},bool) {
	emailregex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	var response = make(map[string]interface{})
	var ok bool
	var user = &User{}
	var err = database.GetDB().Table("users").Where("email=?", authmodel.Email).Find(user).Error
	if err!= nil && err != gorm.ErrRecordNotFound{
		return Message(false,"could not create account please try again"),false
	}
	if user.ID != 0{
		return Message(false,"Email already exist"), false
	}
	if !emailregex.MatchString(authmodel.Email) {
		response = Message(false, "Email Address is incorrect")
		ok = false
		return response , ok
	}

	if len(authmodel.Password) < 6 {
		response = Message(false, "Password must be more than six characters")
		ok = false
		return response,ok
	}
	response = Message(true, "")
	ok = true
	return response,ok
}
func (authmodel *Authmodel) CreateAccount() map[string]interface{}  {
	if resp, ok := authmodel.Validate(); !ok{
		return resp
	}

	hashedpassword,_ := bcrypt.GenerateFromPassword([]byte(authmodel.Password),bcrypt.DefaultCost)
	authmodel.Password = string(hashedpassword)

	user := &User{Email: authmodel.Email, Password: authmodel.Password }
	 returnobj := database.GetDB().Create(user)

    if returnobj.Error!= nil{
    resp := Message(false, "could not create account please try again later")
    return resp
	}else if(user.ID<=0){
		resp := Message(false,"could not create account please try again")
		return resp
	}else {
		resp := Message(true, "User created successfully")
		return resp
	}


}

func Login(Email string, Password string) map[string]interface{} {
   var user = &User{}
   err := database.GetDB().Table("users").Where("email=?", Email).Find(user).Error

   if err!=nil {

   	return  Message(false,"please try again")

   }else if err == gorm.ErrRecordNotFound {

   	return Message(false,"Email or password is incorrect")

   }else {
   	  err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Password))
   	  if err != nil && err == bcrypt.ErrMismatchedHashAndPassword{
   	  	return  Message(false,"user name or password is incorrect")
	  }
   }
 return Message(true,"")
}
