package api

import (
	"forgeturl-server/api/common"
	"forgeturl-server/api/space"
	"forgeturl-server/conf"
	"forgeturl-server/dal/model"
	"forgeturl-server/pkg/maths"

	"github.com/bytedance/sonic"
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

func parsePageIds(pageIdStr string) (pageIds, ownerIds []string, readonlyIds []string, editIds []string, err error) {
	err = sonic.UnmarshalString(pageIdStr, &pageIds)
	if err != nil {
		err = common.ErrInternalServerError(err.Error())
		return
	}

	for _, pageId := range pageIds {
		pType := conf.ParseIdType(pageId)
		switch pType {
		case conf.OwnerPage:
			ownerIds = append(ownerIds, pageId)
		case conf.ReadOnlyPage:
			readonlyIds = append(readonlyIds, pageId)
		case conf.EditPage:
			editIds = append(editIds, pageId)
		}
	}
	return
}

func toPage(uid int64, pageId string, page *model.Page) *space.Page {
	var co []*space.Collections
	_ = sonic.UnmarshalString(page.Content, co)
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
		pageResp.ReadonlyPageId = page.ReadonlyPid
		pageResp.EditPageId = page.EditPid
		pageResp.AdminPageId = page.AdminPid

		pageResp.PageConf.CanEdit = true
		pageResp.PageConf.CanDelete = true

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
