package models

import (
	. "github.com/devtaofeek/ContactApp.Api/Utils"
	. "github.com/devtaofeek/ContactApp.Api/database"
	"github.com/jinzhu/gorm"
	"regexp"
)

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

func (contact Contact) ValidateContact() (map[string]interface{},bool)  {
	emailregex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if contact.Name==""{
    return Message(false,"Error: Contact name is empty"),false
   }
   if contact.Email != "" && !emailregex.MatchString(contact.Email){
	return  Message(false,"Error: Contact Email is incorrect"),false
  }

  if contact.Phone == ""{
	return Message(false,"Error: Contact number is empty"),false
}

return Message(true,""),true
}

func (contact *Contact) Createcontact() map[string]interface{}  {
if  resp, status := contact.ValidateContact(); !status{
	return resp
	}
	var createresponse = GetDB().Create(contact)
	if createresponse.Error !=nil{
		return Message(false,"Error creting contact please try again")
	}
  resp:= Message(true, "contact saved successfully")
	resp["contact"] = contact
return resp
}

func GetContact(contactid int,userid int) map[string]interface{}  {
	if contactid<=0 || userid <= 0{
		return Message(false,"user id or contact id is incorrect")
	}
	contact := &Contact{}
	err := GetDB().Table("contacts").Where("id=? AND userid=?",contactid,userid).First(contact).Error
if err!=nil{
	return Message(false,"could not get record please try again")
}
	if err == gorm.ErrRecordNotFound{
          return Message(false,"record not found")
	}

	resp := Message(true,"")
	returncontact := struct {ID uint; Name string;Email string; Phone string;
	PhoneticName string;NickName string;JobTitle string;Company string;Address string;Relationship string;Userid uint}{}
	returncontact.Userid = contact.Userid
	returncontact.ID = contact.ID
	returncontact.Phone = contact.Phone
	returncontact.Email = contact.Email
	returncontact.Address = contact.Address
	returncontact.JobTitle = contact.JobTitle;
	returncontact.NickName = contact.NickName
	returncontact.PhoneticName = contact.PhoneticName
	returncontact.Relationship = contact.Relationship
	resp["contact"] = returncontact
	return resp
}