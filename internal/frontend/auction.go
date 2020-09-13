package frontend

import (
	"net/http"
	"sellerapp-bidding-system/internal/auction"

	"github.com/julienschmidt/httprouter"
)

type Auction struct {
	auction.AuctionService
}

func RegisterAuctionRoutes(router *httprouter.Router, service auction.AuctionService) {
	auction := Auction{
		AuctionService: service,
	}

	router.POST("/auction", auction.Create)
	router.GET("/auction", auction.Get)
	router.PUT("/auction", auction.Delete)
	router.DELETE("/auction", auction.Delete)
	router.POST("/auction/raise_bid", auction.RaiseBid)
	router.POST("/auction/join", auction.RequestToJoin)
}

func (u Auction) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {}

// Get returns all the auctions
func (u Auction) Get(w http.ResponseWriter, r *http.Request, params httprouter.Params)           {}
func (u Auction) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params)        {}
func (u Auction) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params)        {}
func (u Auction) RequestToJoin(w http.ResponseWriter, r *http.Request, params httprouter.Params) {}
func (u Auction) RaiseBid(w http.ResponseWriter, r *http.Request, params httprouter.Params)      {}
