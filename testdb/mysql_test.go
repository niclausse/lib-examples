package testdb

import (
	"fmt"
	"gorm.io/gorm"
	"testing"
)

var (
	UnitClient, BaseClient *gorm.DB
)

func initDB() {
	cs := make(map[string]MysqlConf)
	cs["unit"] = MysqlConf{
		Name:     "unit",
		Addr:     fmt.Sprintf("%s:%d", dbHost, dbPort),
		Database: "hxx_unit",
		User:     dbDefaultUser,
		Password: dbDefaultPassword,
		Charset:  dbDefaultCharset,
	}
	cs["base"] = MysqlConf{
		Name:     "base",
		Addr:     fmt.Sprintf("%s:%d", dbHost, dbPort),
		Database: "hxx_mis",
		User:     dbDefaultUser,
		Password: dbDefaultPassword,
		Charset:  dbDefaultCharset,
	}

	for s, db := range InitServerAndClients(cs) {
		switch s {
		case "base":
			BaseClient = db
			if err := InitData("./sql/hxx_mis.sql", db); err != nil {
				panic(err)
			}
		case "unit":
			UnitClient = db
			if err := InitData("./sql/hxx_unit.sql", db); err != nil {
				panic(err)
			}
		case "mark":
			if err := InitData("./sql/mark.sql", db); err != nil {
				panic(err)
			}
		}
	}
}

func TestMySQL(t *testing.T) {
	initDB()
}
