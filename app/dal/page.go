package dal

import (
	"context"
	"forgeturl-server/api/common"
	"forgeturl-server/conf"
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

func (*pageImpl) GetSelfPage(ctx context.Context, uid int64) (*model.Page, error) {
	u := Q.Page
	data, err := u.WithContext(ctx).Where(u.UID.Eq(uid)).First()
	if err != nil {
		return nil, transGormErr(err)
	}
	return data, nil
}

// GetPage 可以通过多种id获取页面
func (*pageImpl) GetPage(ctx context.Context, uid int64, pageId string) (*model.Page, error) {
	// 这里缺少一些权限校验，比如是否是该页面的owner，是否是该页面的readonly，是否是该页面的edit，是否是该页面的admin
	u := Q.Page
	do := u.WithContext(ctx)
	pageIdType := conf.ParseIdType(pageId)
	switch pageIdType {
	case conf.OwnerPage:
		do = do.Where(u.Pid.Eq(pageId), u.UID.Eq(uid)) // 默认页面，需要是该页面的owner。需要登录。
	case conf.ReadOnlyPage:
		do = do.Where(u.ReadonlyPid.Eq(pageId)) // 只读页面，谁都能访问。不需要登录。
	case conf.EditPage:
		do = do.Where(u.EditPid.Eq(pageId)) // 编辑页面，谁有链接就能编辑。需要登录。
	case conf.AdminPage:
		do = do.Where(u.AdminPid.Eq(pageId)) // 管理页面，谁有链接就能管理。需要登录。
	default:
		return nil, common.ErrBadRequest("invalid page id type")
	}

	page, err := do.First()
	if err != nil {
		return nil, transGormErr(err)
	}
	return page, nil
}

// GetPageBrief 通过pid获取页面，不获取content内容
func (*pageImpl) GetPageBrief(ctx context.Context, uid int64, pageId string) (*model.Page, error) {
	// 这里缺少一些权限校验，比如是否是该页面的owner，是否是该页面的readonly，是否是该页面的edit，是否是该页面的admin
	pageIdType := conf.ParseIdType(pageId)
	u := Q.Page
	do := u.WithContext(ctx)
	switch pageIdType {
	case conf.OwnerPage:
		do = do.Where(u.Pid.Eq(pageId), u.UID.Eq(uid)) // 默认页面，需要是该页面的owner。需要登录。
	case conf.ReadOnlyPage:
		do = do.Where(u.ReadonlyPid.Eq(pageId)) // 只读页面，谁都能访问。不需要登录。
	case conf.EditPage:
		do = do.Where(u.EditPid.Eq(pageId)) // 编辑页面，谁有链接就能编辑。需要登录。
	case conf.AdminPage:
		do = do.Where(u.AdminPid.Eq(pageId)) // 管理页面，谁有链接就能管理。需要登录。
	default:
		return nil, common.ErrBadRequest("invalid page id type")
	}

	do = do.Omit(u.Content)

	page, err := do.First()
	if err != nil {
		return nil, transGormErr(err)
	}
	return page, nil
}

func (*pageImpl) CheckIsYourPage(ctx context.Context, uid int64, pageIds []string) error {
	u := Q.Page
	do := u.WithContext(ctx)
	do = do.Select(u.ID).Where(u.UID.Eq(uid), u.Pid.In(pageIds...))
	infos, err := do.Find()
	if err != nil {
		return transGormErr(err)
	}
	if len(infos) != len(pageIds) {
		return common.ErrNotYourPageOrPageNotExist()
	}
	return nil
}

const (
	MaskTitle = 1 << iota // 1
	MaskBrief             // 2
	MaskContent
)

func (*pageImpl) UpdatePage(ctx context.Context, uid, mask, version int64, pageId, title, brief, content string) error {
	u := Q.Page
	do := u.WithContext(ctx).Where(u.UID.Eq(uid), u.Pid.Eq(pageId), u.Version.Eq(version))
	do = do.Select(u.Version)
	if mask&MaskTitle != 0 {
		do = do.Select(u.Title)
	}
	if mask&MaskBrief != 0 {
		do = do.Select(u.Brief)
	}
	if mask&MaskContent != 0 {
		do = do.Select(u.Content)
	}
	info, err := do.Updates(&model.Page{
		Title:   title,
		Brief:   brief,
		Content: content,
		Version: version + 1, // 更新版本号
	})
	if err != nil {
		return transGormErr(err)
	}

	if info.RowsAffected == 0 {
		return common.ErrUpdateMissNeedRefreshPage()
	}
	return nil
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

// CleanPage 把某个页面所有内容清除，仅该页面归属者才能清楚
func (*pageImpl) CleanPage(ctx context.Context, uid int64, pid string) error {
	u := Q.Page
	do := u.WithContext(ctx)
	_, err := do.Where(u.UID.Eq(uid), u.Pid.Eq(pid)).UpdateSimple(u.Content.Value(""))
	return transGormErr(err)
}

// UnlinkPage 解除页面链接
// 不真删，只是把该页面链接移除
// 需要该页面创建者，才有权限移除
// pid只能删除不能unlink
func (*pageImpl) UnlinkPage(ctx context.Context, uid int64, pid string) error {
	u := Q.Page
	do := u.WithContext(ctx)
	do = do.Where(u.UID.Eq(uid))
	var err error

	pT := conf.ParseIdType(pid)
	switch pT {
	case conf.OwnerPage:
		// 不允许unlink
		return common.ErrNotSupport()
	case conf.ReadOnlyPage:
		_, err = do.Where(u.ReadonlyPid.Eq(pid)).UpdateSimple(u.ReadonlyPid.Value(""))
	case conf.EditPage:
		_, err = do.Where(u.EditPid.Eq(pid)).UpdateSimple(u.EditPid.Value(""))
	case conf.TempPage:
		// 临时页面不需要unlink
		return common.ErrNotSupport()
	case conf.AdminPage:
		_, err = do.Where(u.AdminPid.Eq(pid)).UpdateSimple(u.AdminPid.Value(""))

	}
	return transGormErr(err)
}

func (*pageImpl) UpdateReadonlyPid(ctx context.Context, pid, readonlyPid string) error {
	u := Q.Page
	do := u.WithContext(ctx)
	_, err := do.Where(u.Pid.Eq(pid)).UpdateSimple(u.ReadonlyPid.Value(readonlyPid))
	return transGormErr(err)
}

func (*pageImpl) UpdateEditPid(ctx context.Context, pid, editPid string) error {
	u := Q.Page
	do := u.WithContext(ctx)
	_, err := do.Where(u.Pid.Eq(pid)).UpdateSimple(u.EditPid.Value(editPid))
	return transGormErr(err)
}

func (*pageImpl) UpdateAdminPid(ctx context.Context, pid, adminPid string) error {
	u := Q.Page
	do := u.WithContext(ctx)
	_, err := do.Where(u.Pid.Eq(pid)).UpdateSimple(u.AdminPid.Value(adminPid))
	return transGormErr(err)
}
