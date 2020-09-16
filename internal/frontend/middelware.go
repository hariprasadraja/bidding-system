package frontend

import (
	"net/http"
	"sellerapp-bidding-system/internal/model"
	"sellerapp-bidding-system/internal/user"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	log "github.com/micro/go-micro/v2/logger"
)

const HeaderAdminID = "X-ADMIN"
const HeaderAuctionID = "X-AUCTION"

// AllowAdmin middleware allows only the Admin role users to perform operations
func (u *User) AllowAdmin(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		claims, err := GetClaims(r)
		if err != nil {
			log.Error("middleware.admin.auth ", err)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		if claims.Role != model.AdminUser {
			log.Info("middleware.admin.invalid")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		err = claims.Login.Valid()
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		userID, err := strconv.ParseInt(claims.Login.Id, 10, 64)
		if err != nil {
			log.Error("invalid.user.token ", err, "login id", claims.Login.Id)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		// Verify user exist or not
		_, err = u.userService.Get(r.Context(), &user.GetRequest{
			Id: userID,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		r.Header.Set(HeaderAdminID, claims.Login.Id)
		next(w, r, params)
	}
}

// AllowUser will allow model.NormalUser role users to perform API actions.
func (u *User) AllowUser(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		claims, err := GetClaims(r)
		if err != nil {
			log.Error("middleware.admin.auth ", err)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		if claims.Role != model.NormalUser {
			log.Info("middleware.admin.invalid")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		err = claims.Login.Valid()
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		log.Info("login claims", claims.Login)
		userID, err := strconv.ParseInt(claims.Login.Id, 10, 64)
		if err != nil {
			log.Error("invalid.user.token ", err, "login id", claims.Login.Id)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		// Verify user exist or not
		_, err = u.userService.Get(r.Context(), &user.GetRequest{
			Id: userID,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		r.Header.Set(HeaderAdminID, claims.Login.Id)
		next(w, r, params)
	}
}

// AllowAuction will allow only the live auction request to perform
func (u *User) AllowAuction(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		claims, err := GetClaims(r)
		if err != nil {
			log.Error("middleware.admin.auth ", err)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		err = claims.Auctions.Valid()
		if err != nil {
			log.Error("error ", err)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		}

		r.Header.Set(HeaderAuctionID, claims.Auctions.Id)
		next(w, r, params)
	}
}

func GetClaims(r *http.Request) (AppJWTClaims, error) {
	tokenStr := r.Header.Get("Authorization")
	jwtToken := strings.TrimPrefix(tokenStr, "Bearer ")
	claims, err := DecJWT(jwtToken)
	if err != nil {
		return AppJWTClaims{}, err
	}

	return claims, nil
}
