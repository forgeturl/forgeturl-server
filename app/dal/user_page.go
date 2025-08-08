package dal

import (
	"context"

	"forgeturl-server/dal/model"
	"forgeturl-server/dal/query"
)

type userPageImpl struct{}

var UserPage = &userPageImpl{}

func (*userPageImpl) GetUserPageIds(ctx context.Context, uid int64, tx ...*query.Query) ([]string, error) {
	u := Q.UserPage
	if len(tx) > 0 {
		u = tx[0].UserPage
	}
	pageids := make([]string, 0)
	err := u.WithContext(ctx).Where(u.ID.Eq(uid)).Order(u.Sort.Asc()).Scan(&pageids)
	if err != nil {
		return nil, transGormErr(err)
	}
	return pageids, nil
}

func (*userPageImpl) SaveUserPageIds(ctx context.Context, uid int64, pageids []string, tx ...*query.Query) error {
	u := Q.UserPage
	if len(tx) > 0 {
		u = tx[0].UserPage
	}
	// 先删除旧的
	_, err := u.WithContext(ctx).Where(u.ID.Eq(uid)).Delete()
	if err != nil {
		return transGormErr(err)
	}

	datas := make([]*model.UserPage, 0, len(pageids))
	for idx, pageid := range pageids {
		datas = append(datas, &model.UserPage{
			ID:   uid,
			Pid:  pageid,
			Sort: int64(idx),
		})
	}
	// 再插入新的
	err = u.WithContext(ctx).Create(datas...)
	if err != nil {
		return transGormErr(err)
	}
	return nil
}
