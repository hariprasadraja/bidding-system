package user

import (
	"bidding-system/internal/backend"
	context "context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/micro/go-micro/errors"

	log "github.com/micro/go-micro/v2/logger"
)

type Handler struct{}

// Create creates a new user record into the `User` table
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

/*
Update updates an already existing user record into the `User` table
*/
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

// Get will get a user record finding either by using email or id
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

// Delete will delete an already existing user record
func (h *Handler) Delete(ctx context.Context, req *DeleteRequest, resp *DeleteResponse) error {
	db := backend.GetConnection(ctx)
	defer backend.PutConnection(db)

	row := db.QueryRow("SELECT EXISTS (SELECT * FROM `User` WHERE id= ?) AS 'count'", req.GetId())
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.BadRequest("user.delete", "user does not exist")
	}

	_, err = sq.ExecContextWith(ctx, db, sq.Delete("User").Where(sq.Eq{
		"id": req.GetId(),
	}))

	if err != nil {
		return err
	}

	resp.Msg = "user deleted successfully."
	return nil
}

// Exist will return all the existing reocrd.
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
