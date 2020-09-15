package model

import "time"

type Auction struct {
	ID                int64     `json:"id,omitempty"`
	Name              string    `json:"name,omitempty"`
	StartTime         time.Time `json:"start_time,omitempty"`
	EndTime           time.Time `json:"end_time,omitempty"`
	StartPriceDisplay string    `json:"start_amount_display,omitempty"`
	StartAmount       float32   `json:"start_amount,omitempty"`
	EndAmount         float32   `json:"end_amount,omitempty"`
	EndPriceDisplay   string    `json:"end_amount_display,omitempty"`
	Currency          string    `json:"currency,omitempty"`
	Status            string    `json:"status,omitempty"`
	WonByUser         string    `json:"won_by_user,omitempty"`
}

type Bid struct {
	UserID    int64   `json:"user_id,omitempty"`
	AuctionID int64   `json:"auction_id,omitempty"`
	BidAmount float32 `json:"bid_amount,omitempty"`
}
