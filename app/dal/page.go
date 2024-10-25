package dal

import (
	"context"

	"forgeturl-server/dal/model"
)

type pageImpl struct {
}

var Page = &pageImpl{}

func (*pageImpl) Create(ctx context.Context, page *model.Page) error {
	u := Q.Page
	err := u.WithContext(ctx).Create(page)
	return transGormErr(err)
}

func (*pageImpl) Get(ctx context.Context, pageId string) (*model.Page, error) {
	u := Q.Page
	page, err := u.WithContext(ctx).Where(u.Pid.Eq(pageId)).First()
	if err != nil {
		return nil, transGormErr(err)
	}
	return page, nil
}
