package auction

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func AddRoutes(router *httprouter.Router) {
	router.GET("/auctions", GetAuction)
}

func GetAuction(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	log.Print("Get Auction Called")

}
