package db

import (
	"database/sql"
	"github.com/spf13/viper"

	"github.com/ferdoran/go-sro/framework/config"

	_ "github.com/go-sql-driver/mysql"
)

func OpenConnAccount() (db *sql.DB) {
	db, errDb := sql.Open(viper.GetString(config.DatabaseAccountDriver), config.ConnStringAccount())
	if errDb != nil {
		panic(errDb.Error())
	}
	return db
}

func OpenConnShard() (db *sql.DB) {
	db, errDb := sql.Open(viper.GetString(config.DatabaseShardDriver), config.ConnStringShard())
	if errDb != nil {
		panic(errDb.Error())
	}
	return db
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
