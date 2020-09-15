package auction

import (
	"context"
	"io"
	"sellerapp-bidding-system/internal/backend"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/micro/go-micro/errors"
	log "github.com/micro/go-micro/v2/logger"
)

type Handler struct{}

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

func (h *Handler) Delete(ctx context.Context, in *DeleteRequest, out *Response) error {
	db := backend.GetConnection(ctx)
	defer backend.PutConnection(db)

	_, err := sq.Delete("Auction").Where(sq.Eq{
		"id": in.GetId(),
	}).RunWith(db).ExecContext(ctx)

	return err
}

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

	return nil
}

func (h *Handler) GetLive(ctx context.Context, in *NoRequest, out *All) error {
	db := backend.GetConnection(ctx)
	defer backend.PutConnection(db)

	now := time.Now()
	rows, err := sq.Select("id,name,start_time,end_time,start_amount").From("Auction").
		Where("start_time >= ? AND end_time <= ?", now, now).
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
			err := rows.Scan(&auction.Id, &auction.Name, &auction.StartTime, &auction.EndTime, &auction.StartAmount)
			if err != nil {
				return err
			}

			response.Bids, err = h.GetBids(ctx, db, auction.Id)
			if err != nil {
				return err
			}

			response.Auction = &auction
		}

		if err := rows.Err(); err != nil {
			return err
		}

		nextResultSet = rows.NextResultSet()
	}

	return nil
}

func (h *Handler) GetAll(ctx context.Context, in *NoRequest, out *All) error {
	db := backend.GetConnection(ctx)
	defer backend.PutConnection(db)

	rows, err := sq.
		Select("*").From("Auction").
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
			err := rows.Scan(&auction.Id, &auction.Name, &auction.StartTime, &auction.EndTime, &auction.StartAmount)
			if err != nil {
				return err
			}

			response.Bids, err = h.GetBids(ctx, db, auction.Id)
			if err != nil {
				return err
			}

			response.Auction = &auction
		}

		if err := rows.Err(); err != nil {
			return err
		}

		nextResultSet = rows.NextResultSet()
	}

	return nil
}

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

	_, err = sq.Insert("Bid").
		Columns("auction_id", "user_id", "bid_amount").
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

func (h *Handler) IsLive(ctx context.Context, db *sqlx.DB, auctionID int64) (bool, error) {
	now := time.Now()
	row := db.QueryRow("SELECT EXISTS (SELECT * FROM `Auction` WHERE id= ? AND start_time >= ? AND end_time <= ?) AS 'count'", auctionID, now, now)
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

func (h *Handler) GetBids(ctx context.Context, db *sqlx.DB, auctionID int64) ([]*Bid, error) {
	rows, err := sq.Select("*").From("Bid").Where(sq.Eq{
		"auction_id": auctionID,
	}).OrderBy("bid_amount").RunWith(db).QueryContext(ctx)

	if err != nil {
		return nil, err
	}

	Bids := make([]*Bid, 0)
	for rows.NextResultSet() && rows.Next() {
		bid := Bid{}
		err := rows.Scan(&bid.AuctionId, &bid.UserId, &bid.BidAmount)
		if err != nil {
			return nil, err
		}

		Bids = append(Bids, &bid)
	}

	return Bids, nil
}
