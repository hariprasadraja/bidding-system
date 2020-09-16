package model

type Auction struct {
	ID                int64   `json:"id,omitempty"`
	Name              string  `json:"name,omitempty"`
	StartTime         string  `json:"start_time,omitempty"`
	EndTime           string  `json:"end_time,omitempty"`
	StartPriceDisplay string  `json:"start_amount_display,omitempty"`
	StartAmount       float32 `json:"start_amount,omitempty"`
	EndAmount         float32 `json:"end_amount,omitempty"`
	EndPriceDisplay   string  `json:"end_amount_display,omitempty"`
	Currency          string  `json:"currency,omitempty"`
	Status            string  `json:"status,omitempty"`
	WonByUser         string  `json:"won_by_user,omitempty"`
}

func (a Auction) CreateValidate() error {
	errors := make(ValidationError)
	if a.Name == "" {
		errors["name"] = "auction name is required."
	}

	if a.StartAmount == 0 {
		errors["start_amount"] = "start_amount is required."
	}

	if a.StartTime == "" {
		errors["start_time"] = "start_time of the auction is required."
	}

	if a.EndTime == "" {
		errors["end_time"] = "end_time of the auction is required."
	}

	if a.Currency == "" {
		errors["currency"] = "currency unit of the auction is required. eg: USD"
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (a Auction) UpdateValidate() error {
	errors := make(ValidationError)
	if a.ID == 0 {
		errors["id"] = "auction id is required."
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

type Bid struct {
	UserID    int64   `json:"user_id,omitempty"`
	AuctionID int64   `json:"auction_id,omitempty"`
	BidAmount float32 `json:"bid_amount,omitempty"`
}
