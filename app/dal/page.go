package dal

import (
	"context"
	"forgeturl-server/api/common"
	"forgeturl-server/conf"
	"forgeturl-server/dal/model"
	"forgeturl-server/dal/query"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

type pageImpl struct {
}

var Page = &pageImpl{}

func (*pageImpl) Create(ctx context.Context, page *model.Page, tx ...*query.Query) error {
	u := Q.Page
	if len(tx) > 0 {
		u = tx[0].Page
	}
	err := u.WithContext(ctx).Create(page)
	return transGormErr(err)
}

// GetSelfPage 获取用户的默认页面
// 这个页面是用户的默认页面，通常是用户的个人主页
func (*pageImpl) GetSelfPage(ctx context.Context, uid int64, tx ...*query.Query) (*model.Page, error) {
	u := Q.Page
	if len(tx) > 0 {
		u = tx[0].Page
	}
	data, err := u.WithContext(ctx).Where(u.UID.Eq(uid)).First()
	if err != nil {
		return nil, transGormErr(err)
	}
	return data, nil
}

// GetAllSelfPages 获取用户所有页面
func (*pageImpl) GetAllSelfPages(ctx context.Context, uid int64, tx ...*query.Query) ([]*model.Page, error) {
	u := Q.Page
	if len(tx) > 0 {
		u = tx[0].Page
	}
	data, err := u.WithContext(ctx).Where(u.UID.Eq(uid)).Order(u.ID.Asc()).Find()
	if err != nil {
		return nil, transGormErr(err)
	}
	return data, nil
}

// GetPage 可以通过多种id获取页面
func (*pageImpl) GetPage(ctx context.Context, uid int64, pageId string, tx ...*query.Query) (*model.Page, error) {
	// 这里缺少一些权限校验，比如是否是该页面的owner，是否是该页面的readonly，是否是该页面的edit，是否是该页面的admin
	u := Q.Page
	if len(tx) > 0 {
		u = tx[0].Page
	}
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
func (*pageImpl) GetPageBrief(ctx context.Context, uid int64, pageId string, tx ...*query.Query) (*model.Page, error) {
	// 这里缺少一些权限校验，比如是否是该页面的owner，是否是该页面的readonly，是否是该页面的edit，是否是该页面的admin
	pageIdType := conf.ParseIdType(pageId)
	u := Q.Page
	if len(tx) > 0 {
		u = tx[0].Page
	}
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

func (*pageImpl) MustBeYourPage(ctx context.Context, uid int64, pageIds []string, tx ...*query.Query) error {
	u := Q.Page
	if len(tx) > 0 {
		u = tx[0].Page
	}
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

func (*pageImpl) UpdatePage(ctx context.Context, mask, version int64, ownerPageId, title, brief, content string, tx ...*query.Query) error {
	u := Q.Page
	if len(tx) > 0 {
		u = tx[0].Page
	}

	do := u.WithContext(ctx).Where(u.Pid.Eq(ownerPageId), u.Version.Eq(version))

	// 收集需要更新的字段
	fields := []field.Expr{u.Version}
	if mask&MaskTitle != 0 {
		fields = append(fields, u.Title)
	}
	if mask&MaskBrief != 0 {
		fields = append(fields, u.Brief)
	}
	if mask&MaskContent != 0 {
		fields = append(fields, u.Content)
	}

	// 一次性选择所有需要更新的字段
	do = do.Select(fields...)

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

// DeleteByPid 可以通过adminPid、Pid进行页面删除
func (*pageImpl) DeleteByPid(ctx context.Context, uid int64, pid string, tx ...*query.Query) error {
	u := Q.Page
	if len(tx) > 0 {
		u = tx[0].Page
	}
	do := u.WithContext(ctx)
	// 删除页面
	pT := conf.ParseIdType(pid)
	switch pT {
	case conf.OwnerPage:
		// 只能通过owner页面删除
		do = do.Where(u.UID.Eq(uid), u.Pid.Eq(pid))
	case conf.AdminPage:
		// 通过admin页面删除
		do = do.Where(u.UID.Eq(uid), u.AdminPid.Eq(pid))
	default:
		return common.ErrNotSupport()
	}
	_, err := do.Delete()
	if err != nil {
		return transGormErr(err)
	}
	return nil
}

// CleanPage 把某个页面所有内容清除，仅该页面归属者才能清楚
func (*pageImpl) CleanPage(ctx context.Context, uid int64, pid string, tx ...*query.Query) error {
	u := Q.Page
	if len(tx) > 0 {
		u = tx[0].Page
	}
	do := u.WithContext(ctx)
	_, err := do.Where(u.UID.Eq(uid), u.Pid.Eq(pid)).UpdateSimple(u.Content.Value(""))
	return transGormErr(err)
}

// UnlinkPage 解除页面链接
// 不真删，只是把该页面链接移除
// 需要该页面创建者，才有权限移除
// pid只能删除不能unlink
func (*pageImpl) UnlinkPage(ctx context.Context, uid int64, pid string, tx ...*query.Query) error {
	u := Q.Page
	if len(tx) > 0 {
		u = tx[0].Page
	}
	do := u.WithContext(ctx)
	do = do.Where(u.UID.Eq(uid))
	var err error

	pT := conf.ParseIdType(pid)
	var result gen.ResultInfo
	switch pT {
	case conf.OwnerPage:
		// 不允许unlink
		return common.ErrNotSupport()
	case conf.ReadOnlyPage:
		result, err = do.Where(u.ReadonlyPid.Eq(pid)).UpdateSimple(u.ReadonlyPid.Value(""))
	case conf.EditPage:
		result, err = do.Where(u.EditPid.Eq(pid)).UpdateSimple(u.EditPid.Value(""))
	case conf.TempPage:
		// 临时页面不需要unlink
		return common.ErrNotSupport()
	case conf.AdminPage:
		result, err = do.Where(u.AdminPid.Eq(pid)).UpdateSimple(u.AdminPid.Value(""))
	}
	if err != nil {
		return transGormErr(err)
	}
	if result.RowsAffected == 0 {
		return common.ErrNotYourPageOrLinkNotExist()
	}

	return nil
}

func (*pageImpl) UpdateReadonlyPid(ctx context.Context, pid, readonlyPid string, tx ...*query.Query) error {
	u := Q.Page
	if len(tx) > 0 {
		u = tx[0].Page
	}
	do := u.WithContext(ctx)
	_, err := do.Where(u.Pid.Eq(pid)).UpdateSimple(u.ReadonlyPid.Value(readonlyPid))
	return transGormErr(err)
}

func (*pageImpl) UpdateEditPid(ctx context.Context, pid, editPid string, tx ...*query.Query) error {
	u := Q.Page
	if len(tx) > 0 {
		u = tx[0].Page
	}
	do := u.WithContext(ctx)
	_, err := do.Where(u.Pid.Eq(pid)).UpdateSimple(u.EditPid.Value(editPid))
	return transGormErr(err)
}

func (*pageImpl) UpdateAdminPid(ctx context.Context, pid, adminPid string, tx ...*query.Query) error {
	u := Q.Page
	if len(tx) > 0 {
		u = tx[0].Page
	}
	do := u.WithContext(ctx)
	_, err := do.Where(u.Pid.Eq(pid)).UpdateSimple(u.AdminPid.Value(adminPid))
	return transGormErr(err)
}

// UnlinkPageByType 根据page_type解除页面链接
// 需要该页面创建者才有权限移除
// 返回被移除的链接pid，用于后续清理user_page表
func (*pageImpl) UnlinkPageByType(ctx context.Context, uid int64, ownerPid string, pageType string, tx ...*query.Query) (string, error) {
	u := Q.Page
	if len(tx) > 0 {
		u = tx[0].Page
	}

	// 首先验证所有权并获取页面
	page, err := u.WithContext(ctx).Where(u.UID.Eq(uid), u.Pid.Eq(ownerPid)).First()
	if err != nil {
		return "", transGormErr(err)
	}
	if page == nil {
		return "", common.ErrNotYourPageOrPageNotExist()
	}

	var linkPid string

	switch pageType {
	case conf.ReadOnlyStr:
		linkPid = page.ReadonlyPid
		if linkPid == "" {
			return "", nil // 已经为空，直接返回成功
		}
		_, err = u.WithContext(ctx).Where(u.UID.Eq(uid), u.Pid.Eq(ownerPid)).UpdateSimple(u.ReadonlyPid.Value(""))
	case conf.EditStr:
		linkPid = page.EditPid
		if linkPid == "" {
			return "", nil
		}
		_, err = u.WithContext(ctx).Where(u.UID.Eq(uid), u.Pid.Eq(ownerPid)).UpdateSimple(u.EditPid.Value(""))
	case conf.AdminStr:
		linkPid = page.AdminPid
		if linkPid == "" {
			return "", nil
		}
		_, err = u.WithContext(ctx).Where(u.UID.Eq(uid), u.Pid.Eq(ownerPid)).UpdateSimple(u.AdminPid.Value(""))
	default:
		return "", common.ErrBadRequest("invalid page type")
	}

	if err != nil {
		return "", transGormErr(err)
	}

	return linkPid, nil
}
