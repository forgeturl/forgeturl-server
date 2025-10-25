package api

import (
	"context"
	"forgeturl-server/api/common"
	"forgeturl-server/api/space"
	"forgeturl-server/conf"
	"forgeturl-server/dal"
	"forgeturl-server/dal/model"
	"forgeturl-server/pkg/maths"

	"github.com/bytedance/sonic"
	"github.com/sunmi-OS/gocore/v2/glog"
)

func isNeedLoginPageId(pageId string) bool {
	switch pageId[0] {
	case conf.OwnerPrefix, conf.EditPrefix:
		return true
	case conf.ReadonlyPrefix, conf.TempPrefix:
		return false
	}
	return true
}

func toPage(ctx context.Context, uid int64, pageId string, page *model.Page) *space.Page {
	co := make([]*space.Collections, 0)
	err0 := sonic.UnmarshalString(page.Content, &co)
	if err0 != nil {
		glog.WarnC(ctx, "unmarshal content failed: err0 %v pageId:%v", err0, page.Pid)
	}
	isSelf := page.UID == uid
	pageResp := &space.Page{
		// PageId:      page.Pid,
		PageId:      pageId,
		Title:       page.Title,
		Brief:       page.Brief,
		Collections: co, // Collections are not used in this context
		CreateTime:  page.CreatedAt.Unix(),
		UpdateTime:  page.UpdatedAt.Unix(),
		IsSelf:      isSelf,
		PageConf:    &space.PageConf{},
		Version:     page.Version,
	}

	if page.ReadonlyPid == pageId {
		pageResp.PageConf.ReadOnly = true

		pageResp.ReadonlyPageId = page.ReadonlyPid
	} else if page.EditPid == pageId {
		pageResp.PageConf.CanEdit = true

		pageResp.ReadonlyPageId = page.ReadonlyPid
		pageResp.EditPageId = page.EditPid
	} else if page.AdminPid == pageId {
		pageResp.PageConf.CanEdit = true
		pageResp.PageConf.CanDelete = true

		pageResp.ReadonlyPageId = page.ReadonlyPid
		pageResp.EditPageId = page.EditPid
		pageResp.AdminPageId = page.AdminPid
	} else if page.Pid == pageId || isSelf {
		// 如果是自己的页面，则展示一下信息，并且标最高权限
		pageResp.PageConf.CanEdit = true
		pageResp.PageConf.CanDelete = true

		pageResp.ReadonlyPageId = page.ReadonlyPid
		pageResp.EditPageId = page.EditPid
		pageResp.AdminPageId = page.AdminPid
	}
	return pageResp
}

func toPageBrief(uid int64, pageId string, page *model.Page) *space.PageBrief {
	isSelf := page.UID == uid
	pageResp := &space.PageBrief{
		// PageId:      page.Pid,
		PageId:     pageId,
		Title:      page.Title,
		Brief:      page.Brief,
		CreateTime: page.CreatedAt.Unix(),
		UpdateTime: page.UpdatedAt.Unix(),
		IsSelf:     isSelf,
		PageConf:   &space.PageConf{},
	}

	if page.ReadonlyPid == pageId {
		pageResp.PageConf.ReadOnly = true
	} else if page.EditPid == pageId {
		pageResp.PageConf.CanEdit = true
	} else if page.AdminPid == pageId {
		pageResp.PageConf.CanEdit = true
		pageResp.PageConf.CanDelete = true
	}

	// 如果是自己的页面，则展示一下信息，并且标最高权限
	if isSelf {
		pageResp.PageConf.CanEdit = true
		pageResp.PageConf.CanDelete = true

	}
	return pageResp
}

func genOwnerPageId() string {
	return maths.GenPageID(string(conf.OwnerPrefix))
}

func genReadOnlyPageId() string {
	return maths.GenPageID(string(conf.ReadonlyPrefix))
}

func genEditPageId() string {
	return maths.GenPageID(string(conf.EditPrefix))
}

func genAdminPageId() string {
	return maths.GenPageID(string(conf.AdminPrefix))
}

func canEditPage(ctx context.Context, userInfo *model.User, pageId string) (string, error) {
	uid := userInfo.ID
	pageType := conf.ParseIdType(pageId)
	if pageType.IsReadOnlyPage() {
		return "", common.ErrNotSupport("cannot edit readonly page")
	}
	// 其他情况都可以编辑

	page, err := dal.Page.GetPageBrief(ctx, uid, pageId)
	if err != nil {
		return "", err
	}
	// 如果是自己的页面，特别校验下归属
	if pageType.IsOwnPage() {
		if page.UID != uid {
			return "", common.ErrNotYourPageOrPageNotExist()
		}
	}
	// 如果能编辑，返回原始pid，给后续内部接口更新使用
	return page.Pid, nil
}

func canDeletePage(ctx context.Context, userInfo *model.User, pageId string) error {
	uid := userInfo.ID
	pageType := conf.ParseIdType(pageId)
	if !pageType.IsOwnPage() {
		// 只读和可编辑类型的页面不支持删除
		return common.ErrBadRequest("This page type does not support deletion")
	}
	// 如果是自己的页面，需要验证所有权
	err0 := dal.Page.MustBeYourPage(ctx, uid, []string{pageId})
	if err0 != nil {
		return common.ErrBadRequest("You don't have permission to edit this page")
	}
	return nil
}
