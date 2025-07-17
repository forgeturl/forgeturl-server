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

func toPages(uid int64, pageIds []string, owner, readonly, edit []*model.Page) []*space.Page {
	pages := make([]*space.Page, 0, len(owner)+len(readonly)+len(edit))
	pageMap := map[string]*space.Page{}

	for _, pg := range owner {
		var co []*space.Collections
		_ = sonic.UnmarshalString(pg.Content, co)
		isSelf := pg.UID == uid
		if pg.Pid == "" {
			continue
		}

		pageMap[pg.Pid] = toPage(uid, pg.Pid, pg)


		pageMap[pg.Pid] = &space.Page{
			PageId:      pg.Pid,
			Title:       pg.Title,
			Brief:       pg.Brief,
			Collections: co,
			CreateTime:  pg.CreatedAt.Unix(),
			UpdateTime:  pg.UpdatedAt.Unix(),
			IsSelf:      isSelf,
			PageConf: &space.PageConf{
				ReadOnly:  false,
				CanEdit:   true,
				CanDelete: true,
			},
		}
	}

	for _, pg := range readonly {
		var co []*space.Collections
		_ = sonic.UnmarshalString(pg.Content, co)
		isSelf := pg.UID == uid
		if pg.ReadonlyPid == "" {
			continue
		}
		pageMap[pg.ReadonlyPid] = &space.Page{
			PageId:      pg.ReadonlyPid,
			Title:       pg.Title,
			Brief:       pg.Brief,
			Collections: co,
			CreateTime:  pg.CreatedAt.Unix(),
			UpdateTime:  pg.UpdatedAt.Unix(),
			IsSelf:      isSelf,
			PageConf: &space.PageConf{
				ReadOnly:  true,
				CanEdit:   false,
				CanDelete: false,
			},
		}
	}

	for _, pg := range edit {
		var co []*space.Collections
		_ = sonic.UnmarshalString(pg.Content, co)
		isSelf := pg.UID == uid
		if pg.EditPid == "" {
			continue
		}
		pageMap[pg.EditPid] = &space.Page{
			PageId:      pg.ReadonlyPid,
			Title:       pg.Title,
			Brief:       pg.Brief,
			Collections: co,
			CreateTime:  pg.CreatedAt.Unix(),
			UpdateTime:  pg.UpdatedAt.Unix(),
			IsSelf:      isSelf,
			PageConf: &space.PageConf{
				ReadOnly:  false,
				CanEdit:   true,
				CanDelete: false,
			},
		}
	}

	for _, pageId := range pageIds {
		if page, ok := pageMap[pageId]; ok {
			pages = append(pages, page)
		}
	}

	// 最后以pageIds的顺序展示
	return pages
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
