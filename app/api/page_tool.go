package api

import (
	"forgeturl-server/api/common"
	"forgeturl-server/api/space"
	"forgeturl-server/dal/model"
	"forgeturl-server/pkg/maths"
	"forgeturl-server/conf"

	"github.com/bytedance/sonic"
)


func parsePageId(pageId string) conf.PageType {
	switch pageId[0] {
	case conf.OwnerPrefix:
		return conf.OwnerPage
	case conf.ReadonlyPrefix:
		return conf.ReadOnlyPage
	case conf.EditPrefix:
		return conf.EditPage
	case conf.TempPrefix:
		return conf.TempPage
	}
	return conf.OwnerPage
}

func parsePageIds(pageIdStr string) (pageIds, ownerIds []string, readonlyIds []string, editIds []string, err error) {
	err = sonic.UnmarshalString(pageIdStr, &pageIds)
	if err != nil {
		err = common.ErrInternalServerError(err.Error())
		return
	}

	for _, pageId := range pageIds {
		pType := parsePageId(pageId)
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

func toPages(pageIds []string, owner, readonly, edit []*model.Page) []*space.Page {
	pages := make([]*space.Page, 0, len(owner)+len(readonly)+len(edit))
	pageMap := map[string]*space.Page{}

	for _, pg := range owner {
		pageMap[pg.Pid] = &space.Page{
			PageId:      pg.Pid,
			Title:       pg.Title,
			Content:     "",
			Collections: nil,
			CreateTime:  pg.CreatedAt.Unix(),
			UpdateTime:  pg.UpdatedAt.Unix(),
			IsSelf:      true,
			PageConf: &space.PageConf{
				ReadOnly:  false,
				CanEdit:   true,
				CanDelete: true,
			},
		}
	}

	for _, pg := range readonly {
		pageMap[pg.Pid] = &space.Page{
			PageId:      pg.ReadonlyPid,
			Title:       pg.Title,
			Content:     "",
			Collections: nil,
			CreateTime:  pg.CreatedAt.Unix(),
			UpdateTime:  pg.UpdatedAt.Unix(),
			IsSelf:      true,
			PageConf: &space.PageConf{
				ReadOnly:  true,
				CanEdit:   false,
				CanDelete: false,
			},
		}
	}

	for _, pg := range edit {
		pageMap[pg.Pid] = &space.Page{
			PageId:      pg.ReadonlyPid,
			Title:       pg.Title,
			Content:     "",
			Collections: nil,
			CreateTime:  pg.CreatedAt.Unix(),
			UpdateTime:  pg.UpdatedAt.Unix(),
			IsSelf:      true,
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
