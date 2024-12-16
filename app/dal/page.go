package dal

import (
	"context"
	"forgeturl-server/dal/model"
)

type PageIdType int

const (
	DefaultPageIdType  PageIdType = 0
	ReadonlyPageIdType PageIdType = 1
	EditPageIdType     PageIdType = 2
	AdminPageIdType    PageIdType = 3
)

type pageImpl struct {
}

var Page = &pageImpl{}

func (*pageImpl) Create(ctx context.Context, page *model.Page) error {
	u := Q.Page
	err := u.WithContext(ctx).Create(page)
	return transGormErr(err)
}

// Get 可以通过多种id获取页面
func (*pageImpl) Get(ctx context.Context, pageId string, pageIdType PageIdType) (*model.Page, error) {
	u := Q.Page
	do := u.WithContext(ctx)
	switch pageIdType {
	case DefaultPageIdType:
		do = do.Where(u.Pid.Eq(pageId))
	case ReadonlyPageIdType:
		do = do.Where(u.ReadonlyPid.Eq(pageId))
	case EditPageIdType:
		do = do.Where(u.EditPid.Eq(pageId))
	case AdminPageIdType:
		do = do.Where(u.AdminPid.Eq(pageId))
	}

	page, err := do.First()
	if err != nil {
		return nil, transGormErr(err)
	}
	return page, nil
}

// 首页拉取大致详情使用
func (*pageImpl) GetPages(ctx context.Context, uid int64, ownerIds, readonlyIds, editIds []string) (owner, readonly, edit []*model.Page, err error) {
	u := Q.Page
	do := u.WithContext(ctx)
	if len(ownerIds) > 0 {
		owner, err = do.Where(u.UID.Eq(uid), u.Pid.In(ownerIds...)).Find()
		if err != nil {
			err = transGormErr(err)
			return
		}
	}

	if len(readonlyIds) > 0 {
		readonly, err = do.Where(u.UID.Eq(uid), u.ReadonlyPid.In(readonlyIds...)).Find()
		if err != nil {
			err = transGormErr(err)
			return
		}
	}

	if len(readonlyIds) > 0 {
		edit, err = do.Where(u.UID.Eq(uid), u.EditPid.In(editIds...)).Find()
		if err != nil {
			err = transGormErr(err)
			return
		}
	}

	return
}
