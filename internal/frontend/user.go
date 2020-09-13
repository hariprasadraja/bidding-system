package frontend

import (
	"encoding/json"
	"net/http"
	"sellerapp-bidding-system/internal/model"
	"sellerapp-bidding-system/internal/user"

	log "github.com/micro/go-micro/v2/logger"

	"github.com/julienschmidt/httprouter"
)

type User struct {
	user user.UserService
}

func RegisterUserRoutes(router *httprouter.Router, service user.UserService) {
	user := User{
		user: service,
	}

	router.POST("/user", user.Create)
	router.GET("/user", user.Get)
	router.PUT("/user", user.Update)
	router.DELETE("/user", user.Delete)
	router.POST("/login", user.Login)
}

func (u User) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("user_details", user)

}

// Get Returns all the users
func (u User) Get(w http.ResponseWriter, r *http.Request, params httprouter.Params)    {}
func (u User) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {}
func (u User) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {}
func (u User) Login(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user := model.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
