package model

import "time"

type Auction struct {
	Name       string    `json:"name,omitempty"`
	StartTime  time.Time `json:"start_time,omitempty"`
	EndTime    time.Time `json:"end_time,omitempty"`
	StartPrice string    `json:"start_price,omitempty"`
	EndPrice   string    `json:"end_price,omitempty"`
	WonByUser  string    `json:"won_by_user,omitempty"`
}
