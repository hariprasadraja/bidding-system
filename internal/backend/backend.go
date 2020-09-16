package backend

import (
	"context"
	"fmt"
	"sellerapp-bidding-system/configs"
	"sync"

	log "github.com/micro/go-micro/v2/logger"
	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var connPool = sync.Pool{
	New: func() interface{} {

		userName := viper.GetString(configs.MYSQL_USER)
		password := viper.GetString(configs.MYSQL_PASSWORD)
		port := viper.GetString(configs.MYSQL_PORT)
		host := viper.GetString(configs.MYSQL_HOST)
		dbName := viper.GetString(configs.MYSQL_DB)

		url := fmt.Sprintf("%s:%s@(%s:%s)/%s", userName, password, host, port, dbName)

		db, err := sqlx.Connect("mysql", url)
		if err != nil {
			return err
		}

		return db
	},
}

var closeConn bool

func GetConnection(ctx context.Context) *sqlx.DB {
	if closeConn {
		return nil
	}

	temp := connPool.Get()

	err, ok := temp.(error)
	if ok {
		log.Errorf("unexpected error while establishing connection: %s", err.Error())
		return nil
	}

	db, ok := temp.(*sqlx.DB)
	if !ok {
		log.Errorf("unexpected error while establishing connection: %s", err.Error())
		return nil
	}

	err = db.PingContext(ctx)
	if err != nil {
		log.Errorf("ping failed. ", err)
	}

	log.Infof("mysql-status %+v", db.Stats())
	return db
}

func PutConnection(db *sqlx.DB) {
	connPool.Put(db)
}

func CloseConnection(ctx context.Context) {
	closeConn = true
	db := GetConnection(ctx)
	db.Close()
}
