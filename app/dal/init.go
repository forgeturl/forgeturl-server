package dal

import (
	"fmt"
	"forgeturl-server/dal/query"
	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var Q *query.Query

const nacosDBKey = "appstoredb"

func Init() error {
	var err error
	db, err = gorm.Open(sqlite.Open(viper.C.GetString("sqliteDB.Path")), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("open db fail: %w", err))
	}

	db.Config.CreateBatchSize = 1000
	// enable the TranslateError flag when opening a db connection.
	db.Config.TranslateError = true

	Q = query.Use(db)

	return nil
}
