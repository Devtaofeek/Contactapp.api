package app

import (
	"context"
	"fmt"
	"github.com/devtaofeek/ContactApp.Api/Utils"
	"github.com/devtaofeek/ContactApp.Api/models"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var JwtAuth = func(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		noAuth := []string{"/api/user/login","/api/user/register"}
		secret :=  []byte(os.Getenv("atsecret"))
		incoming := r.URL.Path
		if incoming == "/api/user/refreshtoken"{
            secret = []byte(os.Getenv("rtsecret"))
		}
		for _,value:= range noAuth {
			 if value == incoming {
			 	next.ServeHTTP(w, r)
				 return
			 }
		}
// inspect the incoming token
response := make(map[string]interface{})
tokenheader := r.Header.Get("Authorization")
		var tokenstringarray = strings.Split(tokenheader, " ")
		if tokenheader == ""{
	response = Utils.Message(false,"Request UnAuthorized")
	w.WriteHeader(http.StatusForbidden)
	Utils.Respond(w,response)
	return
}else if len(tokenstringarray) != 2 {
			response = Utils.Message(false,"Request UnAthorized")
			w.WriteHeader(http.StatusForbidden)
			Utils.Respond(w,response)
			return
		}
var tokenpart = tokenstringarray[1]
 token , err := jwt.Parse(tokenpart, func(token *jwt.Token) (interface{}, error) {
 	if _,ok :=  token.Method.(*jwt.SigningMethodHMAC); !ok{
 		response = Utils.Message(false,"Request UnAuthorised")
 		w.WriteHeader(http.StatusForbidden)
 		Utils.Respond(w,response)
	}
	return  secret,nil
 })
 if err != nil{
			response = Utils.Message(false,"Request Unauthorized")
			w.WriteHeader(http.StatusForbidden)
			Utils.Respond(w,response)
			return
 }else if !token.Valid{
			response = Utils.Message(false,"Request Unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			Utils.Respond(w,response)
	        return
		}


		 claims, _ := token.Claims.(jwt.MapClaims)
		userid, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["userid"]), 10, 64)
		accessuuid := claims["access_uuid"].(string)

		err = models.Fetchredis(accessuuid)
		if err!=nil{
			response = Utils.Message(false,"Request Unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			Utils.Respond(w,response)
			return
		}
		if r.URL.Path == "/api/user/logout"{
			ctx := context.WithValue(r.Context(), "accessUuid", accessuuid,)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}
		var ctx = context.WithValue(r.Context(), "user", userid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


