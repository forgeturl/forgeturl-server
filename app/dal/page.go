package dal

import (
	"context"

	"forgeturl-server/dal/model"
)

type pageImpl struct{
}

var Page = &pageImpl{}

func (*pageImpl) Create(ctx context.Context, page *model.Page) error {
	u := Q.User
	return .Create(ctx, page)
}