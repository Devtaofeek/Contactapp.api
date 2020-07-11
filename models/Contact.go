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

func GetContacts(userid uint) map[string]interface{}  {
	contacts:= make([]*Contact,0)
	var err = GetDB().Table("contacts").Where("userid=?",userid).Find(&contacts).Error
	if err!=nil{
		Message(false,"Please try again later")
	}
	resp := Message(true,"All contacts retrieved")

	resp["contacts"] = contacts
	return resp
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
	resp["contact"] = contact
	return resp
}

func DeleteContact(userid int,contactid int)  map[string]interface{} {
	contact:= &Contact{}
var err = GetDB().Table("contacts").Where("userid=? AND id=?",userid,contactid).Delete(contact)
if err!= nil{
	Message(false,"Please try agiain")
}
return Message(true,"Contact deleted successfully")
}

func (contact *Contact) UpdateContact() map[string]interface{}  {

	var errr = GetDB().Save(&contact).Error
	if errr != nil && errr==gorm.ErrRecordNotFound{
		Message(false,"record not found")
	}

resp:=	Message(true,"Contact updated  successfully")
resp["contact"] = contact
return resp

}

