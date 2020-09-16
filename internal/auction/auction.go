package auction

import (
	"context"
	"io"
	"sellerapp-bidding-system/internal/backend"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/leekchan/accounting"
	"github.com/micro/go-micro/errors"
	log "github.com/micro/go-micro/v2/logger"
)

type Handler struct{}

// TimeFormat to save aution time in the database.
const TimeFormat = "2006-01-02 15:04:05"

/*
Create inserts a new auction record into the `Auction` table

It does no validation. you have to send a valid data.
*/
func (h *Handler) Create(ctx context.Context, in *AuctionRequest, out *Response) error {
	db := backend.GetConnection(ctx)
	defer backend.PutConnection(db)

	result, err := sq.Insert("Auction").
		Columns("name", "start_time", "end_time", "start_amount", "currency").
		Values(in.GetName(), in.GetStartTime(), in.GetEndTime(), in.GetStartAmount(), in.GetCurrency()).RunWith(db).ExecContext(ctx)

	if err != nil {
		return err
	}

	out.Id, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

/*
Update updateds an already existing auction record into the `Auction` table

It does no validation. you have to send a valid data.
similar to Create, but `Id` of an existing record is needed.
*/
func (h *Handler) Update(ctx context.Context, in *AuctionRequest, out *Response) error {
	db := backend.GetConnection(ctx)
	defer backend.PutConnection(db)

	updateQuery := sq.Update("Auction").Where(sq.Eq{"id": in.GetId()}).RunWith(db)

	if in.GetName() != "" {
		updateQuery = updateQuery.Set("name", in.GetName())
	}

	if in.GetStartTime() != "" {
		updateQuery = updateQuery.Set("start_time", in.GetStartTime())
	}

	if in.GetEndTime() != "" {
		updateQuery = updateQuery.Set("end_time", in.GetEndTime())
	}

	if in.GetStartAmount() != 0 {
		updateQuery = updateQuery.Set("start_amount", in.GetStartAmount())
	}

	result, err := updateQuery.ExecContext(ctx)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.BadRequest("auction.update", "invalid auction to update")
	}

	return err
}

/*
Delete deletes an already existing auction record in the `Auction` table
*/
func (h *Handler) Delete(ctx context.Context, in *DeleteRequest, out *Response) error {
	db := backend.GetConnection(ctx)
	defer backend.PutConnection(db)

	_, err := sq.Delete("Auction").Where(sq.Eq{
		"id": in.GetId(),
	}).RunWith(db).ExecContext(ctx)

	return err
}

/*
Get returns a single auction record from the `Auction` table

Id of an auction is required in the input to get the results
*/
func (h *Handler) Get(ctx context.Context, in *GetRequest, out *GetResponse) error {
	db := backend.GetConnection(ctx)
	defer backend.PutConnection(db)

	row := sq.Select("id", "name", "start_time", "end_time", "start_amount", "currency").From("Auction").Where(sq.Eq{
		"id": in.GetId(),
	}).RunWith(db).QueryRowContext(ctx)

	out.Auction = &AuctionRequest{}
	err := row.Scan(
		&out.Auction.Id,
		&out.Auction.Name,
		&out.Auction.StartTime,
		&out.Auction.EndTime,
		&out.Auction.StartAmount,
		&out.Auction.Currency,
	)
	if err != nil {
		return err
	}

	out.Bids, err = h.GetBids(ctx, db, out.Auction.Id)
	if err != nil {
		return err
	}

	ac := accounting.Accounting{
		Symbol:         out.Auction.Currency,
		Precision:      2,
		Format:         "%s %v",
		FormatNegative: "%s (%v)",
	}

	for i := range out.Bids {
		out.Bids[i].BidAmountDisplay = ac.FormatMoney(out.Bids[i].GetBidAmount())
	}

	return nil
}

/*
GetLive will  return all the auction record which are in live state from the `Auction` table

live state refers here is the start_time <= current_time AND end_time >= current_time"

*/
func (h *Handler) GetLive(ctx context.Context, in *NoRequest, out *All) error {
	db := backend.GetConnection(ctx)
	defer backend.PutConnection(db)

	now := time.Now().UTC().Format(TimeFormat)
	rows, err := sq.Select("id,name,start_time,end_time,start_amount,currency").From("Auction").
		Where("start_time <= ? AND end_time >= ?", now, now).
		RunWith(db).QueryContext(ctx)

	if err != nil {
		return err
	}
	defer Close(rows)
	out.Auctions = make([]*GetResponse, 0)
	var response GetResponse
	var nextResultSet = true

	for nextResultSet {
		for rows.Next() {
			var auction AuctionRequest
			err := rows.Scan(&auction.Id, &auction.Name, &auction.StartTime, &auction.EndTime, &auction.StartAmount, &auction.Currency)
			if err != nil {
				return err
			}

			response.Bids, err = h.GetBids(ctx, db, auction.Id)
			if err != nil {
				return err
			}

			log.Info("bids", response.Bids)
			ac := accounting.Accounting{
				Symbol:         auction.Currency,
				Precision:      2,
				Format:         "%s %v",
				FormatNegative: "%s (%v)",
			}

			for i := range response.Bids {
				response.Bids[i].BidAmountDisplay = ac.FormatMoney(response.Bids[i].GetBidAmount())
			}

			response.Auction = &auction
			out.Auctions = append(out.Auctions, &response)
		}

		if err := rows.Err(); err != nil {
			return err
		}

		nextResultSet = rows.NextResultSet()
	}

	return nil
}

/*
GetAll will  return all the auction records that are in the `Auction` table

It also do query to the 'Bid' table and get bids record for the auction.

live state refers here is the start_time <= current_time AND end_time >= current_time"
*/
func (h *Handler) GetAll(ctx context.Context, in *NoRequest, out *All) error {
	db := backend.GetConnection(ctx)
	defer backend.PutConnection(db)

	rows, err := sq.
		Select("id,name,start_time,end_time,start_amount,currency").From("Auction").
		RunWith(db).QueryContext(ctx)

	if err != nil {
		return err
	}

	defer Close(rows)
	out.Auctions = make([]*GetResponse, 0)
	var response GetResponse
	var nextResultSet = true
	for nextResultSet {
		for rows.Next() {
			var auction AuctionRequest
			err := rows.Scan(&auction.Id, &auction.Name, &auction.StartTime, &auction.EndTime, &auction.StartAmount, &auction.Currency)
			if err != nil {
				return err
			}

			response.Bids, err = h.GetBids(ctx, db, auction.Id)
			if err != nil {
				return err
			}

			ac := accounting.Accounting{
				Symbol:         auction.Currency,
				Precision:      2,
				Format:         "%s %v",
				FormatNegative: "%s (%v)",
				FormatZero:     "%s --",
			}

			for i := range response.Bids {
				response.Bids[i].BidAmountDisplay = ac.FormatMoney(response.Bids[i].GetBidAmount())
			}

			response.Auction = &auction
			out.Auctions = append(out.Auctions, &response)
		}

		if err := rows.Err(); err != nil {
			return err
		}

		nextResultSet = rows.NextResultSet()
	}

	return nil
}

/*
IncreaseBid update the `amount` coloumn in the `Bid` table.

A model.User can have only one record for an model.Auction in the `Bid` table
*/
func (h *Handler) IncreaseBid(ctx context.Context, in *Bid, out *NoResponse) error {
	db := backend.GetConnection(ctx)
	defer backend.PutConnection(db)

	live, err := h.IsLive(ctx, db, in.GetAuctionId())
	if err != nil {
		return err
	}

	if !live {
		return errors.Timeout("auction.increase_bid.timeout", "auction is not in live.")
	}

	// delete the previous entry for the user in the auction
	_, err = sq.Delete("Bid").Where("auction_id = ? AND user_id = ?", in.GetAuctionId(), in.GetUserId()).RunWith(db).ExecContext(ctx)
	if err != nil {
		return err
	}

	// insert the new update amount
	_, err = sq.Insert("Bid").
		Columns("auction_id", "user_id", "amount").
		Values(in.GetAuctionId(), in.GetUserId(), in.GetBidAmount()).RunWith(db).ExecContext(ctx)

	if err != nil {
		return err
	}

	return nil
}

func Close(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		log.Errorf("failed to close the object: %+v \n error: %s", closer, err.Error())
	}
}

// IsLive checks the given Auction is in live or not
func (h *Handler) IsLive(ctx context.Context, db *sqlx.DB, auctionID int64) (bool, error) {
	now := time.Now().UTC().Format(TimeFormat)

	log.Info("now", now)
	row := db.QueryRow("SELECT EXISTS (SELECT * FROM `Auction` WHERE id= ? AND start_time <= ? AND end_time >= ?) AS 'count'", auctionID, now, now)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

// GetBids will return the Bids associated with an Auction
func (h *Handler) GetBids(ctx context.Context, db *sqlx.DB, auctionID int64) ([]*Bid, error) {
	rows, err := sq.Select("*").From("Bid").Where(sq.Eq{
		"auction_id": auctionID,
	}).OrderBy("amount").RunWith(db).QueryContext(ctx)

	if err != nil {
		return nil, err
	}

	Bids := make([]*Bid, 0)
	nextResultSet := true
	for nextResultSet {
		for rows.Next() {
			bid := Bid{}
			err := rows.Scan(&bid.AuctionId, &bid.UserId, &bid.BidAmount)
			if err != nil {
				return nil, err
			}

			Bids = append(Bids, &bid)
		}

		nextResultSet = rows.NextResultSet()
	}
	return Bids, nil
}
