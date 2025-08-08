package dal

import (
	"context"
	"forgeturl-server/dal/model"
	"forgeturl-server/dal/query"
)

type uniquePidImpl struct {
}

var UniquePid = &uniquePidImpl{}

func (*uniquePidImpl) Create(ctx context.Context, uid int64, pid string, tx ...*query.Query) error {
	u := Q.UniquePid
	if len(tx) > 0 {
		u = tx[0].UniquePid
	}
	err := u.WithContext(ctx).Create(&model.UniquePid{
		Pid: pid,
		UID: uid,
	})
	if err != nil {
		return transGormErr(err)
	}
	return nil
}
