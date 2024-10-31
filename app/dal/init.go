package dal

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"

	"forgeturl-server/api/common"
	"forgeturl-server/dal/query"
	"forgeturl-server/pkg/core"

	"github.com/sunmi-OS/gocore/v2/db/orm"
	"gorm.io/gorm"
)

var db *gorm.DB
var Q *query.Query

const nacosDBKey = "mysql"

func Init() error {
	name := nacosDBKey
	if err0 := orm.NewOrUpdateDB(nacosDBKey); err0 != nil {
		return fmt.Errorf("init db(%v) failed, err: %v", name, err0)
	}

	db = orm.GetORM(nacosDBKey)
	if db == nil {
		return fmt.Errorf("get db(%v) failed", name)
	}

	db.Config.CreateBatchSize = 1000
	// enable the TranslateError flag when opening a db connection.
	db.Config.TranslateError = true

	Q = query.Use(db)

	return nil
}

func transGormErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return core.WrapError(common.ErrNotFound(), funcName(3))
	} else if errors.Is(err, gorm.ErrDuplicatedKey) {
		return core.WrapError(common.ErrConflict(err.Error()), funcName(3))
	} else {
		return core.WrapError(common.ErrInternalServerError(err.Error()), funcName(3))
	}
}

// funcName get func name.
func funcName(skip int) (name string) {
	if _, file, lineNo, ok := runtime.Caller(skip); ok {
		return file + ":" + strconv.Itoa(lineNo)
	}
	return "unknown:0"
}
