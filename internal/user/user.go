package user

import (
	context "context"
	"sellerapp-bidding-system/internal/backend"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/micro/go-micro/errors"

	log "github.com/micro/go-micro/v2/logger"
)

type Handler struct{}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (h *Handler) Create(ctx context.Context, req *CreateRequest, resp *CreateResponse) error {
	db := backend.GetConnection(ctx)
	defer backend.PutConnection(db)
	row := db.QueryRow("SELECT EXISTS (SELECT * FROM `User` WHERE email= ?) AS 'count'", req.GetEmail())
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	if count != 0 {
		return errors.BadRequest("user.create.email", "email already exist.")
	}

	result, err := db.Exec("INSERT INTO `User`(name, `role`, date_created, date_modified, email,password) VALUES(?,?,?,?,?,?);", req.GetName(), req.GetRole(), time.Now(), time.Now(), req.GetEmail(), req.GetPassword())
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	resp.Id = id
	log.Info("User Created - id: ", id)
	resp.Msg = "user created successfully."
	return nil
}

func (h *Handler) Update(ctx context.Context, req *UpdateRequest, resp *UpdateResponse) error {
	db := backend.GetConnection(ctx)
	defer backend.PutConnection(db)

	updateQuery := sq.Update("User")
	if req.GetName() != "" {
		updateQuery = updateQuery.Set("name", req.GetName())
	}

	if req.GetEmail() != "" {
		updateQuery = updateQuery.Set("email", req.GetEmail())
	}

	updateQuery = updateQuery.Set("date_modified", time.Now())
	_, err := sq.ExecContextWith(ctx, db, updateQuery)
	if err != nil {
		return err
	}

	log.Infof("updated user - %s", req.GetId)
	resp.Msg = "user updated successfully."
	return nil
}

func (h *Handler) Get(ctx context.Context, req *GetRequest, resp *GetResponse) error {
	log.Info(req.String())
	db := backend.GetConnection(ctx)
	defer backend.PutConnection(db)

	selectQuery := sq.Select("id,name,role,email").From("User")
	equal := make(sq.Eq)
	if req.GetId() != 0 {
		equal["id"] = req.GetId()
	}

	if req.GetEmail() != "" {
		equal["email"] = req.GetEmail()
	}

	row := selectQuery.Where(equal).RunWith(db).QueryRowContext(ctx)
	err := row.Scan(&resp.Id, &resp.Name, &resp.Role, &resp.Email)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) Delete(ctx context.Context, req *DeleteRequest, resp *DeleteResponse) error {
	db := backend.GetConnection(ctx)
	defer backend.PutConnection(db)

	_, err := sq.ExecContextWith(ctx, db, sq.Delete("User").Where(sq.Eq{
		"id": req.GetId(),
	}))

	if err != nil {
		return err
	}

	log.Infof("updated user - %s", req.GetId)
	resp.Msg = "user deleted successfully."
	return nil
}

func (h *Handler) Exist(ctx context.Context, req *ExistRequest, resp *ExistResponse) error {
	db := backend.GetConnection(ctx)
	defer backend.PutConnection(db)
	row := db.QueryRow("SELECT EXISTS (SELECT * FROM `User` WHERE email= ? AND password = ?) AS 'count'", req.GetEmail(), req.GetPassword())
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.BadRequest("user.login", "email/password must be wrong.")
	}

	return nil
}
