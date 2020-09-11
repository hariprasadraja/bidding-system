package auction

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-micro/util/log"
)

func AddRoutes(router *httprouter.Router) {
	router.GET("/auctions", GetAuction)
	router.GET("/auctions/rise_bid", RiseBid)
	router.PUT("/auctions", GetAuction)
	router.DELETE("/auctions", DeleteAuction)
}

func GetAuction(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// Get Auction By ID
	// Get All Auction in Progress  -> only for normal users
	// Get All Auction  -> only for admins
	log.Info("Get Auction called.")
}

func RiseBid(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	log.Info("Bid Auction called.")
}

func CreateAuction(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	log.Info("Create Auction called.")
}

func UpdateAuction(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	log.Info("Update Auction called.")
}

func DeleteAuction(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	log.Info("Delete Auction called.")
}
