package api

import (
	"context"

	"forgeturl-server/api/common"
	"forgeturl-server/api/space"
	"forgeturl-server/conf"
	"forgeturl-server/dal"
	"forgeturl-server/dal/query"
	"forgeturl-server/pkg/middleware"
	gapi "github.com/sunmi-OS/gocore/v2/api"
)

func validateSharedPageID(id string) error {
	if !conf.ValidPageID(id) {
		return common.ErrBadRequest("invalid share id")
	}
	switch conf.ParseIdType(id) {
	case conf.ReadOnlyPage, conf.EditPage, conf.AdminPage:
		return nil
	default:
		return common.ErrBadRequest("expected a shared page id")
	}
}

func (s spaceServiceImpl) AddSharedPage(c *gapi.Context, req *space.AddSharedPageReq) (*space.AddSharedPageResp, error) {
	return addSharedPage(c.Request.Context(), middleware.GetLoginUid(c), req.PageId)
}

func addSharedPage(ctx context.Context, uid int64, id string) (*space.AddSharedPageResp, error) {
	if uid == 0 {
		return nil, common.ErrNeedLogin("")
	}
	if err := validateSharedPageID(id); err != nil {
		return nil, err
	}
	resp := &space.AddSharedPageResp{}
	err := dal.WithUserPageLock(ctx, uid, func(tx *query.Query) error {
		page, err := dal.Page.GetPageBrief(ctx, uid, id, tx)
		if err != nil {
			return err
		}
		if page.UID == uid {
			resp.PageBrief = toPageBrief(uid, page.Pid, page)
			resp.AlreadyAdded = true
			resp.PageIds, err = dal.UserPage.GetUserPageIds(ctx, uid, tx)
			return err
		}
		resp.PageBrief = toPageBrief(uid, id, page)
		resp.PageIds, resp.AlreadyAdded, err = dal.UserPage.AppendSharedPage(ctx, uid, id, tx)
		return err
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s spaceServiceImpl) RemoveSharedPage(c *gapi.Context, req *space.RemoveSharedPageReq) (*space.RemoveSharedPageResp, error) {
	return removeSharedPage(c.Request.Context(), middleware.GetLoginUid(c), req.PageId)
}

func removeSharedPage(ctx context.Context, uid int64, id string) (*space.RemoveSharedPageResp, error) {
	if uid == 0 {
		return nil, common.ErrNeedLogin("")
	}
	if err := validateSharedPageID(id); err != nil {
		return nil, err
	}
	resp := &space.RemoveSharedPageResp{}
	err := dal.WithUserPageLock(ctx, uid, func(tx *query.Query) error {
		page, err := dal.Page.GetPageBrief(ctx, uid, id, tx)
		if err != nil && !common.IsErrNotFound(err) {
			return err
		}
		if err == nil && page.UID == uid {
			return common.ErrBadRequest("owned pages cannot be removed as a shared entry")
		}
		if _, err = dal.UserPage.DeleteUserPageId(ctx, uid, id, tx); err != nil {
			return err
		}
		resp.PageIds, err = dal.UserPage.GetUserPageIds(ctx, uid, tx)
		return err
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
