package frontend

import (
	"encoding/json"
	"net/http"
	"sellerapp-bidding-system/internal/auction"
	"sellerapp-bidding-system/internal/model"
	"sellerapp-bidding-system/internal/user"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-micro/v2/util/log"
)

type Auction struct {
	AuctionService auction.AuctionService
	UserService    user.UserService
}

func RegisterAuctionRoutes(router *httprouter.Router, userService user.UserService, auctionService auction.AuctionService) {
	auction := Auction{
		AuctionService: auctionService,
		UserService:    userService,
	}

	user := User{
		userService: userService,
	}

	router.POST("/auction", user.AllowAdmin(auction.Create))
	router.PUT("/auction", user.AllowAdmin(auction.Update))
	router.DELETE("/auction", user.AllowAdmin(auction.Delete))
	router.GET("/auction/status", user.AllowUser(auction.Status))
	router.POST("/auction/raise_bid", user.AllowAuction(user.AllowUser(auction.RaiseBid)))
	router.GET("/auction/all", user.AllowAdmin(auction.GetAll))
	router.GET("/auction/live_only", user.AllowUser(auction.GetLiveOnly))
	router.POST("/auction/join/:id", user.AllowUser(auction.RequestToJoin))

}

const TimeFormat = "2006-01-02 15:04:05"

func (u Auction) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	auctionModel := model.Auction{}
	err := json.NewDecoder(r.Body).Decode(&auctionModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = auctionModel.CreateValidate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(TimeFormat+" MST", auctionModel.StartTime)
	if err != nil {
		http.Error(w, "start_time: "+err.Error(), http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(TimeFormat+" MST", auctionModel.EndTime)
	if err != nil {
		http.Error(w, "end_time: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := u.AuctionService.Create(r.Context(), &auction.AuctionRequest{
		Name:        auctionModel.Name,
		StartTime:   startTime.UTC().Format(TimeFormat),
		EndTime:     endTime.UTC().Format(TimeFormat),
		StartAmount: auctionModel.StartAmount,
		Currency:    auctionModel.Currency,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RenderJSON(w, Response{
		ID:         resp.Id,
		StatusText: http.StatusText(http.StatusCreated),
		Status:     http.StatusCreated,
		Message:    "auction created successfully.",
	})

}

func (u Auction) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	auctionModel := model.Auction{}
	err := json.NewDecoder(r.Body).Decode(&auctionModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = auctionModel.CreateValidate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(TimeFormat+" MST", auctionModel.StartTime)
	if err != nil {
		http.Error(w, "start_time: "+err.Error(), http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(TimeFormat+" MST", auctionModel.EndTime)
	if err != nil {
		http.Error(w, "start_time: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := u.AuctionService.Update(r.Context(), &auction.AuctionRequest{
		Id:          auctionModel.ID,
		Name:        auctionModel.Name,
		StartTime:   startTime.UTC().Format(TimeFormat),
		EndTime:     endTime.UTC().Format(TimeFormat),
		StartAmount: auctionModel.StartAmount,
		Currency:    auctionModel.Currency,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RenderJSON(w, Response{
		ID:         resp.Id,
		StatusText: http.StatusText(http.StatusCreated),
		Status:     http.StatusCreated,
		Message:    "auction updated successfully.",
	})
}

func (u Auction) GetAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	result, err := u.AuctionService.GetAll(r.Context(), new(auction.NoRequest))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RenderJSON(w, Response{
		StatusText: http.StatusText(http.StatusOK),
		Status:     http.StatusOK,
		Data:       result,
	})

}
func (u Auction) GetLiveOnly(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	result, err := u.AuctionService.GetLive(r.Context(), new(auction.NoRequest))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RenderJSON(w, Response{
		StatusText: http.StatusText(http.StatusOK),
		Status:     http.StatusOK,
		Data:       result,
	})
}

// Status will return status of the current auctions for the user
func (u Auction) Status(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	auctionID := r.Header.Get(HeaderAuctionID)
	id, err := strconv.ParseInt(auctionID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := u.AuctionService.Get(r.Context(), &auction.GetRequest{
		Id: id,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RenderJSON(w, Response{
		StatusText: http.StatusText(http.StatusOK),
		Status:     http.StatusOK,
		Data:       result,
	})
}
func (u Auction) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	auctionID := params.ByName("id")
	id, err := strconv.ParseInt(auctionID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = u.AuctionService.Delete(r.Context(), &auction.DeleteRequest{
		Id: id,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RenderJSON(w, Response{
		StatusText: http.StatusText(http.StatusOK),
		Status:     http.StatusOK,
		Message:    "auction deleted successfully.",
	})
}

func (u Auction) RequestToJoin(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	auctionIDStr := params.ByName("id")
	auctionID, err := strconv.ParseInt(auctionIDStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := u.AuctionService.Get(r.Context(), &auction.GetRequest{
		Id: auctionID,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(TimeFormat, result.Auction.StartTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(TimeFormat, result.Auction.EndTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	if now.Before(startTime) {
		http.Error(w, "auction not yet started.", http.StatusBadRequest)
		return
	}

	if now.Before(startTime) {
		http.Error(w, "auction already done. item sold.", http.StatusBadRequest)
		return
	}

	claims, err := GetClaims(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims.Auctions = jwt.StandardClaims{
		//NOTE: set to local time, since claims.Valid() works on local time.
		ExpiresAt: endTime.Unix(),
		Id:        auctionIDStr,
		IssuedAt:  time.Now().Unix(),
		Issuer:    r.URL.RequestURI(),
		NotBefore: startTime.Unix(),
		Subject:   "Auction Token",
	}

	token, err := EncJWT(claims)
	if err != nil {
		log.Error("request_to_join.jwt", err)
		http.Error(w, "sorry,something went wrong", http.StatusInternalServerError)
		return
	}

	RenderJSON(w, Response{
		Status:     http.StatusOK,
		StatusText: http.StatusText(http.StatusOK),
		Message:    "Access Granted",
		Token:      token,
	})
}

func (u Auction) RaiseBid(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	auctionIDStr := r.Header.Get(HeaderAuctionID)
	auctionID, err := strconv.ParseInt(auctionIDStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bidAmountStr := r.URL.Query().Get("bid_amount")
	bidAmount, err := strconv.ParseFloat(bidAmountStr, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	adminID, err := strconv.ParseInt(r.Header.Get(HeaderAdminID), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = u.AuctionService.IncreaseBid(r.Context(), &auction.Bid{
		AuctionId: auctionID,
		UserId:    adminID,
		BidAmount: float32(bidAmount),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
