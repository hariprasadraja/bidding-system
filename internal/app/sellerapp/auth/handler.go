package auth

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func AddRoutes(router *httprouter.Router) {
	router.POST("/login", Login)
}

func Login(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	log.Print("Login called.")
}
