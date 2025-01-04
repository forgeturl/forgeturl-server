package api

import (
	"forgeturl-server/api/common"
	"forgeturl-server/api/space"
	"forgeturl-server/dal/model"
	"forgeturl-server/pkg/maths"

	"github.com/bytedance/sonic"
)

type PageType int

const (
	OwnerPage    PageType = 0 // start with O
	ReadOnlyPage PageType = 1 // start With R
	EditPage     PageType = 2 // start With E

	OwnerPrefix    = uint8('O')
	ReadonlyPrefix = uint8('R')
	EditPrefix     = uint8('E')
)

func parsePageId(pageId string) PageType {
	switch pageId[0] {
	case OwnerPrefix:
		return OwnerPage
	case ReadonlyPrefix:
		return ReadOnlyPage
	case EditPrefix:
		return EditPage
	}
	return OwnerPage
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
		case OwnerPage:
			ownerIds = append(ownerIds, pageId)
		case ReadOnlyPage:
			readonlyIds = append(readonlyIds, pageId)
		case EditPage:
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
	return maths.GenPageID(string(OwnerPrefix))
}

func genReadOnlyPageId() string {
	return maths.GenPageID(string(ReadonlyPrefix))
}

func genEditPageId() string {
	return maths.GenPageID(string(EditPrefix))
}
