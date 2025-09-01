package dal

import (
	"context"
	"forgeturl-server/api/common"
	"forgeturl-server/conf"

	"forgeturl-server/dal/model"
	"forgeturl-server/dal/query"

	"gorm.io/gen"
)

type userPageImpl struct{}

var UserPage = &userPageImpl{}

func (*userPageImpl) GetUserPageIds(ctx context.Context, uid int64, tx ...*query.Query) ([]string, error) {
	u := Q.UserPage
	if len(tx) > 0 {
		u = tx[0].UserPage
	}
	pageids := make([]string, 0)
	err := u.WithContext(ctx).Select(u.Pid).Where(u.UID.Eq(uid)).Order(u.Sort.Asc()).Scan(&pageids)
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
	_, err := u.WithContext(ctx).Where(u.UID.Eq(uid)).Delete()
	if err != nil {
		return transGormErr(err)
	}

	datas := make([]*model.UserPage, 0, len(pageids))
	for idx, pageid := range pageids {
		datas = append(datas, &model.UserPage{
			UID:  uid,
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

func (*userPageImpl) DeleteUserPageId(ctx context.Context, uid int64, pageid string, tx ...*query.Query) (int64, error) {
	u := Q.UserPage
	if len(tx) > 0 {
		u = tx[0].UserPage
	}
	result, err := u.WithContext(ctx).Where(u.UID.Eq(uid), u.Pid.Eq(pageid)).Delete()
	if err != nil {
		return 0, transGormErr(err)
	}
	return result.RowsAffected, nil
}

func (*userPageImpl) BatchRemovePageLink(ctx context.Context, pageid string, tx ...*query.Query) (int64, error) {
	// 批量删除pageid
	u := Q.UserPage
	if len(tx) > 0 {
		u = tx[0].UserPage
	}
	pT := conf.ParseIdType(pageid)
	var result gen.ResultInfo
	if pT == conf.OwnerPage {
		// 删除该页面，必须是通过uid删除
		return 0, common.ErrNotSupport()
	}

	result, err := u.WithContext(ctx).Where(u.Pid.Eq(pageid)).Delete()
	if err != nil {
		return 0, transGormErr(err)
	}
	return result.RowsAffected, nil
}
