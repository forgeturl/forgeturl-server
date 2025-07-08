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
func (*pageImpl) Get(ctx context.Context, uid int64, pageId string, pageIdType PageIdType) (*model.Page, error) {
	// 这里缺少一些权限校验，比如是否是该页面的owner，是否是该页面的readonly，是否是该页面的edit，是否是该页面的admin
	u := Q.Page
	do := u.WithContext(ctx)
	switch pageIdType {
	case DefaultPageIdType:
		do = do.Where(u.Pid.Eq(pageId), u.UID.Eq(uid)) // 默认页面，需要是该页面的owner
	case ReadonlyPageIdType:
		do = do.Where(u.ReadonlyPid.Eq(pageId)) // 只读页面，谁都能访问
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

// GetPages 首页拉取大致详情使用
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

	if len(editIds) > 0 {
		edit, err = do.Where(u.UID.Eq(uid), u.EditPid.In(editIds...)).Find()
		if err != nil {
			err = transGormErr(err)
			return
		}
	}

	return
}

func (*pageImpl) DeleteByPid(ctx context.Context, uid int64, pid string) error {
	u := Q.Page
	do := u.WithContext(ctx)
	// 删除页面
	_, err := do.Where(u.UID.Eq(uid), u.Pid.Eq(pid)).Delete()
	if err != nil {
		return transGormErr(err)
	}
	return nil
}

// UnlinkPage 解除页面链接
// 不真删，只是把该页面链接移除
// 需要时该页面的人，才有权限移除
func (*pageImpl) UnlinkPage(ctx context.Context, uid int64, readOnlyPid, editPid, adminPid string) error {
	u := Q.Page
	do := u.WithContext(ctx)
	do = do.Where(u.UID.Eq(uid))
	var err error
	if readOnlyPid != "" {
		_, err = do.Where(u.ReadonlyPid.Eq(readOnlyPid)).UpdateSimple(u.ReadonlyPid.Value(""))
	} else if editPid != "" {
		_, err = do.Where(u.EditPid.Eq(editPid)).UpdateSimple(u.EditPid.Value(""))
	} else if adminPid != "" {
		_, err = do.Where(u.AdminPid.Eq(adminPid)).UpdateSimple(u.AdminPid.Value(""))
	}
	return transGormErr(err)
}
