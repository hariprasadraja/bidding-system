package backend

import (
	"os"
	"sync"

	log "github.com/micro/go-micro/v2/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var connPool = sync.Pool{
	New: func() interface{} {
		db, err := sqlx.Connect("mysql", "root:RYEfc53!t9t@(localhost:3306)/bidding_app")
		if err != nil {
			return err
		}

		return db
	},
}

var closeConn bool

func GetConnection() *sqlx.DB {
	if closeConn {
		return nil
	}

	temp := connPool.Get()

	err, ok := temp.(error)
	if !ok {
		log.Errorf("unexpected error while establishing connection: %s", err.Error())
		os.Exit(1)
	}

	db, ok := temp.(*sqlx.DB)
	if !ok {
		log.Errorf("unexpected error while establishing connection: %s", err.Error())
		os.Exit(1)
	}

	return db
}

func PutConnection(db *sqlx.DB) error {
	connPool.Put(db)
	return nil
}

func CloseConnection() {
	closeConn = true
	db := GetConnection()
	db.Close()
}
