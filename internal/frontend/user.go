package frontend

import (
	"encoding/json"
	"net/http"
	"sellerapp-bidding-system/internal/model"
	"sellerapp-bidding-system/internal/user"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/micro/go-micro/errors"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/julienschmidt/httprouter"
)

type User struct {
	userService user.UserService
}

func RegisterUserRoutes(router *httprouter.Router, service user.UserService) {
	user := User{
		userService: service,
	}

	router.POST("/user", user.AllowAdmin(user.Create))
	router.GET("/user", user.AllowAdmin(user.Get))
	router.PUT("/user", user.AllowAdmin(user.Update))
	router.DELETE("/user", user.AllowAdmin(user.Delete))
	router.POST("/authenticate", user.Authenticate)
}

type Response struct {
	ID         int64       `json:"id,omitempty"`
	StatusText string      `json:"status,omitempty"`
	Status     int         `json:"-"`
	Message    string      `json:"message,omitempty"`
	Token      string      `json:"token,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func (u User) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Info("create user called.")
	userModel := model.User{}
	err := json.NewDecoder(r.Body).Decode(&userModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = userModel.CreateValidate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	userModel.Password, err = OneTimeEnc(userModel.Password)
	if err != nil {
		log.Error("login.password.enc ", err.Error())
		return
	}
	resp, err := u.userService.Create(r.Context(), &user.CreateRequest{
		Name:     userModel.Name,
		Email:    userModel.Email,
		Role:     model.NormalUser,
		Password: userModel.Password,
	})

	if err != nil {
		log.Error("user.create.failed", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RenderJSON(w, Response{
		StatusText: http.StatusText(http.StatusOK),
		Status:     http.StatusOK,
		Message:    resp.GetMsg(),
		ID:         resp.GetId(),
	})
}

// Get Returns all the users
func (u User) Get(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, errors.BadRequest("/get", "`id` in query params is required.").Error(), http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := u.userService.Get(r.Context(), &user.GetRequest{
		Id: userID,
	})

	if err != nil {
		log.Error("user.create.failed", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := model.User{
		ID:           resp.GetId(),
		Name:         resp.GetName(),
		Role:         resp.GetRole(),
		Email:        resp.GetEmail(),
		DateCreated:  time.Time{},
		DateModified: time.Time{},
	}

	RenderJSON(w, Response{
		StatusText: http.StatusText(http.StatusFound),
		Status:     http.StatusFound,
		Message:    "user retrived successfully.",
		Data:       user,
	})

}
func (u User) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userModel := model.User{}
	err := json.NewDecoder(r.Body).Decode(&userModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = userModel.UpdateValidate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	log.Info("user_details", userModel)
	resp, err := u.userService.Update(r.Context(), &user.UpdateRequest{
		Id:    userModel.ID,
		Name:  userModel.Name,
		Email: userModel.Email,
		Role:  model.NormalUser,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RenderJSON(w, Response{
		StatusText: http.StatusText(http.StatusOK),
		Status:     http.StatusOK,
		Message:    resp.GetMsg(),
	})

}
func (u User) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, errors.BadRequest("/delete", "`id` in query params is required.").Error(), http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := u.userService.Delete(r.Context(), &user.DeleteRequest{
		Id: userID,
	})

	if err != nil {
		log.Error("user.create.failed", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RenderJSON(w, Response{
		StatusText: http.StatusText(http.StatusOK),
		Status:     http.StatusOK,
		Message:    resp.GetMsg(),
	})
}
func (u User) Authenticate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userModel := model.User{}
	err := json.NewDecoder(r.Body).Decode(&userModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userModel.Password, err = OneTimeEnc(userModel.Password)
	if err != nil {
		log.Error("Authenticate.password.enc ", err.Error())
		http.Error(w, "sorry, something went wrong.", http.StatusInternalServerError)
		return
	}

	log.Info("passoword ", userModel.Password)
	_, err = u.userService.Exist(r.Context(), &user.ExistRequest{
		Email:    userModel.Email,
		Password: userModel.Password,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := u.userService.Get(r.Context(), &user.GetRequest{
		Email: userModel.Email,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	userID := strconv.FormatInt(user.Id, 10)
	token, err := EncJWT(AppJWTClaims{
		Role: user.GetRole(),
		Login: jwt.StandardClaims{
			Audience:  userModel.Email,
			ExpiresAt: time.Now().Add((24 * 30) * time.Hour).Unix(),
			Id:        userID,
			IssuedAt:  time.Now().Unix(),
			Issuer:    "sellerapp.bidding.backend",
			Subject:   "login token",
		},
	})

	if err != nil {
		log.Error("authenticate.jwt.error ", err.Error())
		http.Error(w, "sorry, something went wrong.", http.StatusInternalServerError)
		return
	}

	RenderJSON(w, Response{
		StatusText: http.StatusText(http.StatusOK),
		Status:     http.StatusOK,
		Message:    "loggedIn successfull",
		Token:      token,
	})
}

/*
RenderJSON renders the data as a JSON file and sends it as a http response
it returns error when data unable to marshall or write into http.ResponseWriter,
Example
		data  := struct {
			Key string
			Value string
		}{
			"key",
			"value"
		}

	dataLen, err := RenderJSON(w, jsonData)
	if err != nil {
		log.Panicf("Failed to RenderJSON \n ERROR: %s \n Data: %#v \n Length: %v", err, data, dataLen)
	}
*/
func RenderJSON(w http.ResponseWriter, data interface{}) {
	js, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Error("frontend.user.renderjson", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		log.Error("frontend.user.renderjson", err)
		return
	}
}
