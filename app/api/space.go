package api

import (
	"forgeturl-server/api/common"
	"forgeturl-server/api/space"
	"forgeturl-server/conf"
	"forgeturl-server/dal"
	"forgeturl-server/pkg/middleware"

	"github.com/bytedance/sonic"
	"github.com/sunmi-OS/gocore/v2/api"
)

type spaceServiceImpl struct {
}

func (s spaceServiceImpl) SavePageIds(context *api.Context, req *space.SavePageIdsReq) (*space.SavePageIdsResp, error) {
	ctx := context.Request.Context()
	uid := middleware.GetLoginUid(context)
	if uid == 0 {
		return nil, common.ErrNeedLogin("")
	}

	pageIds := req.PageIds

	// 如果是owner page则必须是他自己的才能保存
	ownerIds := make([]string, 0)
	for _, pageId := range pageIds {
		if conf.ParseIdType(pageId) == conf.OwnerPage {
			ownerIds = append(ownerIds, pageId)
		}
	}
	err := dal.Page.CheckIsYourPage(ctx, uid, ownerIds)
	if err != nil {
		return nil, err
	}
	
	pageIdsStr, err := sonic.MarshalString(pageIds)
	if err != nil {
		return nil, common.ErrBadRequest(err.Error())
	}
	err = dal.User.UpdatePageIds(ctx, uid, pageIdsStr)
	if err != nil {
		return nil, err
	}
	return &space.SavePageIdsResp{}, nil
}

func NewSpaceService() space.SpaceServiceHTTPServer {
	return &spaceServiceImpl{}
}

// GetMySpace 只需要给个大概
func (s spaceServiceImpl) GetMySpace(context *api.Context, req *space.GetMySpaceReq) (*space.GetMySpaceResp, error) {
	ctx := context.Request.Context()
	uid := req.Uid
	userInfo, err := dal.User.Get(ctx, uid)
	if err != nil {
		return nil, err
	}
	resp := &space.GetMySpaceResp{
		SpaceName: userInfo.DisplayName,
		Pages:     make([]*space.Page, 0),
	}

	pageIds := make([]string, 0)
	_ = sonic.UnmarshalString(userInfo.PageIds, &pageIds)

	for _, pageId := range pageIds {
		pageType := conf.ParseIdType(pageId)
		page, err := dal.Page.GetPage(ctx, uid, pageId, pageType)
		if err != nil {
			return nil, err
		}
		pageResp := toPage(uid, pageId, page)
		resp.Pages = append(resp.Pages, pageResp)
	}

	return resp, nil
}

func (s spaceServiceImpl) GetPage(context *api.Context, req *space.GetPageReq) (*space.GetPageResp, error) {
	ctx := context.Request.Context()
	// 获取某个页面数据
	uid := middleware.GetLoginUid(context)
	pageId := req.PageId
	if isNeedLoginPageId(pageId) && uid == 0 {
		return nil, common.ErrNeedLogin("")
	}

	pageType := conf.ParseIdType(req.PageId)
	// 拉取某个页面数据

	page, err := dal.Page.GetPage(ctx, uid, pageId, pageType)
	if err != nil {
		return nil, err
	}

	pageResp := toPage(uid, pageId, page)
	return &space.GetPageResp{Page: pageResp}, nil
}

func (s spaceServiceImpl) CreateTmpPage(context *api.Context, req *space.CreateTmpPageReq) (*space.CreateTmpPageResp, error) {
	ctx := context.Request.Context()
	_ = ctx
	//TODO implement me
	panic("implement me")
}

func (s spaceServiceImpl) UpdatePage(context *api.Context, req *space.UpdatePageReq) (*space.UpdatePageResp, error) {
	ctx := context.Request.Context()
	uid := middleware.GetLoginUid(context)
	if uid == 0 {
		return nil, common.ErrNeedLogin("")
	}

	err := dal.Page.UpdatePage(ctx, uid, req.Mask, req.Version, req.PageId, req.Title, req.Brief, req.Content)
	if err != nil {
		return nil, err
	}
	return &space.UpdatePageResp{}, nil
}

func (s spaceServiceImpl) DeletePage(context *api.Context, req *space.DeletePageReq) (*space.DeletePageResp, error) {
	ctx := context.Request.Context()
	uid := middleware.GetLoginUid(context)
	if uid == 0 {
		return nil, common.ErrNeedLogin("")
	}

	err := dal.Page.DeleteByPid(ctx, uid, req.PageId)
	if err != nil {
		return nil, err
	}
	return &space.DeletePageResp{}, nil
}

func (s spaceServiceImpl) RemovePageLink(context *api.Context, req *space.RemovePageLinkReq) (*space.RemovePageLinkResp, error) {
	ctx := context.Request.Context()
	uid := middleware.GetLoginUid(context)
	if uid == 0 {
		return nil, common.ErrNeedLogin("")
	}

	err := dal.Page.UnlinkPage(ctx, uid, req.PageId)
	if err != nil {
		return nil, err
	}
	return &space.UnlinkPageResp{}, nil
}

func (s spaceServiceImpl) CreatePageLink(context *api.Context, req *space.CreateNewPageLinkReq) (*space.CreateNewPageLinkResp, error) {
	ctx := context.Request.Context()
	_ = ctx
	//TODO implement me
	panic("implement me")
}
