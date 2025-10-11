package dal

import (
	"context"
	"forgeturl-server/dal/model"
	"time"
)

type userImpl struct {
}

var User = &userImpl{}

func (*userImpl) Create(ctx context.Context, user *model.User) error {
	u := Q.User
	err := u.WithContext(ctx).Create(user)
	return transGormErr(err)
}

func (*userImpl) Get(ctx context.Context, uid int64) (*model.User, error) {
	u := Q.User
	page, err := u.WithContext(ctx).Where(u.ID.Eq(uid)).First()
	if err != nil {
		return nil, transGormErr(err)
	}
	return page, nil
}

func (*userImpl) GetByExternalID(ctx context.Context, externalID string) (*model.User, error) {
	u := Q.User
	page, err := u.WithContext(ctx).Where(u.ExternalID.Eq(externalID)).First()
	if err != nil {
		return nil, transGormErr(err)
	}
	return page, nil
}

func (*userImpl) UpdateDisplayName(ctx context.Context, uid int64, displayName string) error {
	u := Q.User
	_, err := u.WithContext(ctx).Where(u.ID.Eq(uid)).UpdateSimple(u.DisplayName.Value(displayName))
	return transGormErr(err)
}

func (*userImpl) UpdateLastLoginTime(ctx context.Context, uid int64, lastLoginTime time.Time) error {
	u := Q.User
	_, err := u.WithContext(ctx).Where(u.ID.Eq(uid)).UpdateSimple(u.LastLoginDate.Value(lastLoginTime))
	return transGormErr(err)
}
