package frontend

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"sellerapp-bidding-system/internal/model"
	"sellerapp-bidding-system/internal/user"
	"strconv"
	"time"

	"github.com/gorilla/securecookie"

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

	router.POST("/user", user.Create)
	router.GET("/user", user.Get)
	router.PUT("/user", user.Update)
	router.DELETE("/user", user.Delete)
	router.POST("/authenticate", user.Authenticate)
}

type Response struct {
	Status  string      `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Token   string      `json:"token,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (u User) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userModel := model.User{}
	err := json.NewDecoder(r.Body).Decode(&userModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("user_details", userModel)

	enc := sha256.New()
	_, err = enc.Write([]byte(userModel.Password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userModel.Password = base64.StdEncoding.EncodeToString(enc.Sum(nil))
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
		Status:  http.StatusText(http.StatusOK),
		Message: resp.GetMsg(),
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
		Status:  http.StatusText(http.StatusOK),
		Message: "user retrived successfully.",
		Data:    user,
	})

}
func (u User) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userModel := model.User{}
	err := json.NewDecoder(r.Body).Decode(&userModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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
		Status:  http.StatusText(http.StatusOK),
		Message: resp.GetMsg(),
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
		Status:  http.StatusText(http.StatusOK),
		Message: resp.GetMsg(),
	})
}
func (u User) Authenticate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userModel := model.User{}
	err := json.NewDecoder(r.Body).Decode(&userModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := u.userService.Exist(r.Context(), &user.ExistRequest{
		Email:    userModel.Email,
		Password: userModel.Password,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !result.GetExist() {
		http.Error(w, "invalid email/password", http.StatusNotFound)
		return
	}

	// TODO:  Use redis cache and save `Id` as key and `Audience` as value
	token := EncJWT(AppJWTClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  userModel.Email,
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Id:        string(securecookie.GenerateRandomKey(10)),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "sellerapp.bidding.backend",
			Subject:   "login token",
		},
	})

	RenderJSON(w, Response{
		Status:  http.StatusText(http.StatusOK),
		Message: "loggedIn successfull",
		Token:   token,
	})
}

type AppJWTClaims struct {
	Auctions []int64 `json:"auctions,omitempty"`
	jwt.StandardClaims
}

var secretKey = "yww@y9eNApn4Nsb@Hm4Z3&Uee9zjKJwtVn^%eXdW$Q#igUGvVHNeYh5iEDK!VfhwLSuLVhpHb9vo4uFuuWZm6B4jnTgU6cmefyveF$!2T7PM8^mEnjM9eJ#mAk2amkCK"

// EncryptJwt with a secretKey with claims
func EncJWT(claims AppJWTClaims) (JWTToken string) {
	hmac512Algo := jwt.SigningMethodHS512.Alg()
	signingMethod := jwt.GetSigningMethod(hmac512Algo)
	token := jwt.New(signingMethod)

	token.Claims = claims

	// Signing and Serialization.
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Error("Error: ", err)
		return
	}

	return tokenString
}

// DecJWT decrypts the JWT token and returns it's claims
func DecJWT(tokenString string, secretKey []byte) (claims jwt.Claims, err error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, errors.Forbidden("dec.jwt.signing", "Unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if parsedToken.Valid {
		if claims, ok := parsedToken.Claims.(AppJWTClaims); ok {
			return claims, nil
		}
	}

	return nil, errors.Forbidden("dec.jwt.invalid", "invalid token in request")
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
