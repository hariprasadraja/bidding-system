package user

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-micro/util/log"
)

func AddRoutes(router *httprouter.Router) {
	router.POST("/login", Login)
}

func Login(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	log.Info("Login Called")
}

func CreateUser(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	log.Info("Create User Called")
}

func UpdateUser(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	log.Info("Update User Called")
}

func Delete(w http.ResponseWriter, r *http.Request, param httprouter.Param) {
	log.Info("Delete User called.")
}

func GetUser(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	log.Info("Get User called.")
}
