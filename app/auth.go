package app

import (
	"context"
	"fmt"
	"github.com/devtaofeek/ContactApp.Api/Utils"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var JwtAuth = func(next http.Handler) http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		noAuth := []string{"/api/user/login","/api/user/register"}
		secret :=  []byte(os.Getenv("atsecret"))
		incoming := request.URL.Path
		if incoming == "/api/user/refreshtoken"{
            secret = []byte(os.Getenv("rtsecret"))
		}
		for _,value:= range noAuth {
			 if value == incoming {
			 	next.ServeHTTP(writer,request)
				 return
			 }
		}
// inspect the incoming token
response := make(map[string]interface{})
tokenheader := request.Header.Get("Authorization")
		var tokenstringarray = strings.Split(tokenheader, " ")
		if tokenheader == ""{
	response = Utils.Message(false,"Request UnAuthorized")
	writer.WriteHeader(http.StatusForbidden)
	Utils.Respond(writer,response)
	return
}else if len(tokenstringarray) != 2 {
			response = Utils.Message(false,"Request UnAthorized")
			writer.WriteHeader(http.StatusForbidden)
			Utils.Respond(writer,response)
			return
		}
var tokenpart = tokenstringarray[1]
 token , err := jwt.Parse(tokenpart, func(token *jwt.Token) (interface{}, error) {
 	if _,ok :=  token.Method.(*jwt.SigningMethodHMAC); !ok{
 		response = Utils.Message(false,"Request UnAuthorised")
 		writer.WriteHeader(http.StatusForbidden)
 		Utils.Respond(writer,response)
	}
	return  secret,nil
 })
 if err != nil{
			response = Utils.Message(false,"Request Unauthorized")
			writer.WriteHeader(http.StatusForbidden)
			Utils.Respond(writer,response)
			return
 }else if !token.Valid{
			response = Utils.Message(false,"Request Unauthorized")
			writer.WriteHeader(http.StatusSeeOther)
			Utils.Respond(writer,response)
	        return
		}
		 claims, _ := token.Claims.(jwt.MapClaims)
		userid, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["userid"]), 10, 64)
		ctx:= context.WithValue(request.Context(),"user", userid)

		next.ServeHTTP(writer,request.WithContext(ctx))
	})
}
