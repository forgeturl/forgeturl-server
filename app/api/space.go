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

func NewSpaceService() space.SpaceServiceHTTPServer {
	return &spaceServiceImpl{}
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
		if conf.ParseIdType(pageId).IsOwnPage() {
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

// GetMySpace 只需要给个大概
func (s spaceServiceImpl) GetMySpace(context *api.Context, req *space.GetMySpaceReq) (*space.GetMySpaceResp, error) {
	ctx := context.Request.Context()
	uid := req.Uid
	userInfo, err := dal.User.Get(ctx, uid)
	if err != nil {
		return nil, err
	}
	resp := &space.GetMySpaceResp{
		SpaceName:  userInfo.DisplayName,
		PageBriefs: make([]*space.PageBrief, 0),
	}

	pageIds := make([]string, 0)
	_ = sonic.UnmarshalString(userInfo.PageIds, &pageIds)

	for _, pageId := range pageIds {
		page, err := dal.Page.GetPageBrief(ctx, uid, pageId)
		if err != nil {
			return nil, err
		}
		pageResp := toPageBrief(uid, pageId, page)
		resp.PageBriefs = append(resp.PageBriefs, pageResp)
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

	// 拉取某个页面数据

	page, err := dal.Page.GetPage(ctx, uid, pageId)
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
	return &space.RemovePageLinkResp{}, nil
}

func (s spaceServiceImpl) CreatePageLink(context *api.Context, req *space.CreatePageLinkReq) (*space.CreatePageLinkResp, error) {
	// 当前页面是你的，则你可以创建 readonly edit admin链接
	// 当前如果你有该页面的adminId，则可以创建 readonly edit链接
	// 其他情况会被拒绝

	// 如果同样的链接已存在，则需要让用户RemoveLink后，再创建新的链接。避免用户以为，同一个页面可以存在多个链接。
	ctx := context.Request.Context()
	uid := middleware.GetLoginUid(context)
	if uid == 0 {
		return nil, common.ErrNeedLogin("")
	}

	data, err := dal.Page.GetPageBrief(ctx, uid, req.PageId)
	if err != nil {
		return nil, err
	}
	inputPageType := conf.ParseIdType(req.PageId)
	canCreateAdminPid := data.UID == uid
	pageTypeStr := req.PageType

	newPageId := ""
	switch pageTypeStr {
	case conf.ReadOnlyStr:
		// todo 如果已经存在直接返回？

		if data.ReadonlyPid != "" {
			return nil, common.ErrBadRequest("readonly link already exists")
		}
		newPageId = genReadOnlyPageId()
		err = dal.Page.UpdateReadonlyPid(ctx, data.Pid, newPageId)

	case conf.EditStr:
		// 如果inputPageType是readonly的，则不允许创建
		if inputPageType == conf.ReadOnlyPage {
			return nil, common.ErrBadRequest("cannot create edit link from readonly page")
		}
		if data.EditPid != "" {
			return nil, common.ErrBadRequest("edit link already exists")
		}
		newPageId = genEditPageId()
		err = dal.Page.UpdateEditPid(ctx, data.Pid, newPageId)

	case conf.AdminStr:
		if data.AdminPid != "" {
			return nil, common.ErrBadRequest("admin link already exists")
		}
		if !canCreateAdminPid {
			return nil, common.ErrBadRequest("you are not the owner of this page, cannot create admin link")
		}
		newPageId = genAdminPageId()
		err = dal.Page.UpdateAdminPid(ctx, data.Pid, newPageId)

	default:
		return nil, common.ErrBadRequest("invalid page type")
	}

	if err != nil {
		return nil, err
	}

	return &space.CreatePageLinkResp{
		PageType:  pageTypeStr,
		NewPageId: newPageId,
	}, nil
}
