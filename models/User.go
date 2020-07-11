package models

import (
	. "github.com/devtaofeek/ContactApp.Api/Utils"
	"github.com/devtaofeek/ContactApp.Api/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"os"
	"regexp"
	"time"
)

type User struct {
	gorm.Model
	Email string `gorm:"type:varchar(100);unique_index"`
	Password string
	Contacts []Contact `gorm:"foreignkey:Userid"`
}
var client *redis.Client
type TokenDetails struct {
	ID    uint64
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}
type Authmodel struct {
	Email    string
	Password string
}

func init()  {
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
}
func (authmodel *Authmodel) Validate() (map[string]interface{},bool) {
	emailregex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	var response = make(map[string]interface{})
	var ok bool
	var user = &User{}
	var err = database.GetDB().Table("users").Where("email=?", authmodel.Email).Find(user).Error
	CreateToken(uint64(user.ID))
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


		token,err := CreateToken(uint64(user.ID))
		if err!= nil{
			return Message(false,"could not complete the request please try again")
		}
		resp := Message(true, "User created successfully")
		resp["token"] = token
		return resp
	}


}

func Login(Email string, Password string) map[string]interface{} {

   var user = &User{}
   err := database.GetDB().Table("users").Where("email=?", Email).Find(user).Error

   if err!=nil && err== gorm.ErrRecordNotFound {

   	return  Message(false,"Email or password id incorrect")

   }else {
	   err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Password))
	   if err != nil && err == bcrypt.ErrMismatchedHashAndPassword{
		   return  Message(false,"user name or password is incorrect")
	   }
   }



token,err := CreateToken(uint64(user.ID))
if err!= nil{
	return Message(false,"could not complete the request please try again")
}
resp := Message(true,"Login successful")
resp["token"] = token
return resp
}

func CreateToken(id uint64) (*TokenDetails, error) {
	tokendetails := &TokenDetails{}
	tokendetails.ID  = id
	tokendetails.AtExpires = time.Now().Add(time.Hour * 2).Unix() // TODO reset access token to five
	tokendetails.AccessUuid  = uuid.New().String()
	tokendetails.RefreshUuid = uuid.New().String()
	tokendetails.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()

	atclaims := jwt.MapClaims{}
	atclaims["authorized"] = true
	atclaims["access_uuid"] = tokendetails.AccessUuid
	atclaims["userid"] = id
	atclaims["exp"] = tokendetails.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256,atclaims)
	var err error
	tokendetails.AccessToken, err = at.SignedString([]byte(os.Getenv("atsecret")))
	if err!= nil {
		return nil, err
	}

	rtclaims := jwt.MapClaims{}
	rtclaims["refreshuuid"] = tokendetails.RefreshUuid
	rtclaims["refreshexpiry"] = tokendetails.RtExpires
	rtclaims["access_uuid"] = tokendetails.AccessUuid
	rtclaims["userid"] = id
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256,rtclaims)
	tokendetails.RefreshToken , err = rt.SignedString([]byte(os.Getenv("rtsecret")))
	if err!= nil{
		return nil, err
	}

	err = PersistToken(*tokendetails)
if err!=nil{
	return nil,err
}
return tokendetails, nil
}


func RefreshToken(id int) (map[string]interface{},error) {
	token, err:= CreateToken(uint64(id))
	if err!=nil{
		return  Message(false,"Invalid Request"), err
	}
	resp := Message(true,"")
	resp["token"] = token
	return resp ,err
}


func PersistToken(details TokenDetails) error  {
at:= time.Unix(details.AtExpires,0)
rt:= time.Unix(details.RtExpires,0)
now := time.Now()
access := client.Set(details.AccessUuid, int(details.ID),at.Sub(now))
if access.Err() != nil{
	return access.Err()
}

refresh := client.Set(details.RefreshUuid, int(details.ID),rt.Sub(now))

if refresh.Err() != nil{
	return refresh.Err()
}
return nil
}

func Fetchredis(accessuid string) error {
	_,err := client.Get(accessuid).Result()
	if err!=nil{
		return err
	}
	return nil
}

func DeleteFromredis(accessuuid string) error  {
  _ , erraccess  := client.Del(accessuuid).Result()


  if erraccess!=nil{
  	return erraccess
  }
return nil
}
