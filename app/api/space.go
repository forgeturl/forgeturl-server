package api

import (
	"forgeturl-server/api/common"
	"forgeturl-server/api/space"
	"forgeturl-server/conf"
	"forgeturl-server/dal"
	"forgeturl-server/dal/model"
	"forgeturl-server/dal/query"
	"forgeturl-server/pkg/middleware"

	"github.com/bytedance/sonic"
	"github.com/sunmi-OS/gocore/v2/api"
)

type spaceServiceImpl struct {
}

func NewSpaceService() space.SpaceServiceHTTPServer {
	return &spaceServiceImpl{}
}

func (s spaceServiceImpl) GetUserInfo(context *api.Context, req *space.GetUserInfoReq) (*space.GetUserInfoResp, error) {
	// 已登录才能获取到详情，否则拉拉取不到
	ctx := context.Request.Context()
	uid := req.Uid
	loginUid := middleware.GetLoginUid(context)
	userInfo, err := dal.User.Get(ctx, uid)
	if err != nil {
		return nil, err
	}

	// 如果是自己，则返回内容多些
	info := &space.GetUserInfoResp{
		Uid:         uid,
		DisplayName: userInfo.DisplayName,
		Avatar:      userInfo.Avatar,
		Status:      userInfo.Status,

		Username:      userInfo.Username,
		Email:         userInfo.Email,
		LastLoginTime: userInfo.LastLoginDate.Unix(),
		IsAdmin:       userInfo.IsAdmin,
		Provider:      userInfo.Provider,
		CreateTime:    userInfo.CreatedAt.Unix(),
		UpdateTime:    userInfo.UpdatedAt.Unix(),
	}

	// 如果不是自己，则返回内容少些
	if loginUid != uid {
		info.Username = ""
		info.Email = ""
		info.LastLoginTime = 0
		info.IsAdmin = 0
		info.Provider = ""
		info.CreateTime = 0
		info.UpdateTime = 0
	}

	return info, nil
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

	newPageIds := make([]string, 0, len(pageIds))
	err = dal.Q.Transaction(func(tx *query.Query) error {
		err0 := dal.UserPage.SaveUserPageIds(ctx, uid, pageIds, tx)
		if err0 != nil {
			return err0
		}

		newPageIds, err0 = dal.UserPage.GetUserPageIds(ctx, uid, tx)
		if err0 != nil {
			return err0
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &space.SavePageIdsResp{PageIds: newPageIds}, nil
}

// GetMySpace 只需要给个大概
func (s spaceServiceImpl) GetMySpace(context *api.Context, req *space.GetMySpaceReq) (*space.GetMySpaceResp, error) {
	ctx := context.Request.Context()
	// 获取某个页面数据
	uid := middleware.GetLoginUid(context)
	userInfo, err := dal.User.Get(ctx, uid)
	if err != nil {
		return nil, err
	}
	resp := &space.GetMySpaceResp{
		SpaceName:  userInfo.DisplayName,
		PageBriefs: make([]*space.PageBrief, 0),
	}

	pageIds, err := dal.UserPage.GetUserPageIds(ctx, uid)
	if err != nil {
		return nil, err
	}

	for _, pageId := range pageIds {
		page, err0 := dal.Page.GetPageBrief(ctx, uid, pageId)
		if err0 != nil {
			return nil, err0
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

func (s spaceServiceImpl) CreatePage(context *api.Context, req *space.CreatePageReq) (*space.CreatePageResp, error) {
	// 首先搜下，他有几个页面
	ctx := context.Request.Context()
	// 获取某个页面数据
	uid := middleware.GetLoginUid(context)
	if uid == 0 {
		return nil, common.ErrNeedLogin("")
	}

	content, err := sonic.MarshalString(req.Collections)
	if err != nil {
		return nil, err
	}

	var pageId string
	startVersion := int64(0)
	var pageIds []string
	err = dal.Q.Transaction(func(tx *query.Query) error {
		page, err0 := dal.Page.GetSelfPage(ctx, uid, tx)
		if err0 == nil && page != nil && page.ID > 0 {
			return common.ErrBadRequest("You already have a self page, cannot create more")
		}
		if !common.IsErrNotFound(err0) {
			if err0 != nil {
				return err0
			}
		}

		pageId = genOwnerPageId()
		err0 = dal.UniquePid.Create(ctx, uid, pageId, tx)
		if err0 != nil {
			return err0
		}

		err0 = dal.Page.Create(ctx, &model.Page{
			UID:         uid,
			Pid:         pageId,
			ReadonlyPid: "",
			EditPid:     "",
			AdminPid:    "",
			Title:       req.Title,
			Brief:       req.Brief,
			Content:     content,
			Version:     startVersion,
		}, tx)
		if err0 != nil {
			return err0
		}

		pageIds, err0 = dal.UserPage.GetUserPageIds(ctx, uid, tx)
		if err0 != nil {
			return err0
		}

		pageIds = append([]string{pageId}, pageIds...)
		err0 = dal.UserPage.SaveUserPageIds(ctx, uid, pageIds, tx)
		if err0 != nil {
			return err0
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &space.CreatePageResp{
		PageId:  pageId,
		Version: startVersion,
		PageIds: pageIds,
	}, nil
}

func (s spaceServiceImpl) UpdatePage(context *api.Context, req *space.UpdatePageReq) (*space.UpdatePageResp, error) {
	ctx := context.Request.Context()
	uid := middleware.GetLoginUid(context)
	if uid == 0 {
		return nil, common.ErrNeedLogin("")
	}

	content, err := sonic.MarshalString(req.Collections)
	if err != nil {
		return nil, err
	}
	// 先看下页面是否存在，且是自己的
	err = dal.Page.CheckIsYourPage(ctx, uid, []string{req.PageId})
	if err != nil {
		return nil, err
	}

	err = dal.Page.UpdatePage(ctx, uid, req.Mask, req.Version, req.PageId, req.Title, req.Brief, content)
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
	pid := req.PageId

	err := dal.Q.Transaction(func(tx *query.Query) error {
		err0 := dal.Page.DeleteByPid(ctx, uid, pid, tx)
		if err0 != nil {
			return err0
		}

		_, err0 = dal.UserPage.DeleteUserPageId(ctx, uid, pid, tx)
		if err0 != nil {
			return err0
		}

		// 还有一个unique_pid先保留，不删除

		return nil
	})
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

	err := dal.Q.Transaction(func(tx *query.Query) error {
		// 只有创建者可以移除该页面，其他人移除会报错
		err0 := dal.Page.UnlinkPage(ctx, uid, req.PageId)
		if err0 != nil {
			return err0
		}

		_, err0 = dal.UserPage.BatchRemovePageLink(ctx, req.PageId, tx)
		if err0 != nil {
			return err0
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &space.RemovePageLinkResp{}, nil
}

func (s spaceServiceImpl) AddPageLink(context *api.Context, req *space.AddPageLinkReq) (*space.AddPageLinkResp, error) {
	// 1.当前页面若是你的，则你可以创建 readonly edit admin链接
	// 2.如果你有该页面的adminId，则可以创建 readonly edit链接
	// 3. 其他情况会被拒绝
	// 4. 如果同样的链接已存在，则需要让用户RemoveLink后，再创建新的链接。避免用户以为，同一个页面可以存在多个链接。
	ctx := context.Request.Context()
	uid := middleware.GetLoginUid(context)
	if uid == 0 {
		return nil, common.ErrNeedLogin("")
	}
	pageTypeStr := req.PageType
	newPageId := ""

	err := dal.Q.Transaction(func(tx *query.Query) error {
		data, err0 := dal.Page.GetPageBrief(ctx, uid, req.PageId, tx)
		if err0 != nil {
			return err0
		}
		inputPageType := conf.ParseIdType(req.PageId)
		canCreateAdminPid := data.UID == uid

		switch pageTypeStr {
		case conf.ReadOnlyStr:
			// todo 如果已经存在直接返回？

			if data.ReadonlyPid != "" {
				return common.ErrBadRequest("readonly link already exists")
			}
			newPageId = genReadOnlyPageId()

			err0 = dal.UniquePid.Create(ctx, uid, newPageId, tx)
			if err0 != nil {
				return err0
			}

			// 如果之前有存在其他的readonlyPid，需要提示该用户该页面被删除了，让他自己决定是否要删除
			err0 = dal.Page.UpdateReadonlyPid(ctx, data.Pid, newPageId, tx)

		case conf.EditStr:
			// 如果inputPageType是readonly的，则不允许创建
			if inputPageType == conf.ReadOnlyPage {
				return common.ErrBadRequest("cannot create edit link from readonly page")
			}
			if data.EditPid != "" {
				return common.ErrBadRequest("edit link already exists")
			}
			newPageId = genEditPageId()

			err0 = dal.UniquePid.Create(ctx, uid, newPageId, tx)
			if err0 != nil {
				return err0
			}

			err0 = dal.Page.UpdateEditPid(ctx, data.Pid, newPageId, tx)

		case conf.AdminStr:
			if data.AdminPid != "" {
				return common.ErrBadRequest("admin link already exists")
			}
			if !canCreateAdminPid {
				return common.ErrBadRequest("you are not the owner of this page, cannot create admin link")
			}
			newPageId = genAdminPageId()

			err0 = dal.UniquePid.Create(ctx, uid, newPageId, tx)
			if err0 != nil {
				return err0
			}

			err0 = dal.Page.UpdateAdminPid(ctx, data.Pid, newPageId, tx)

		default:
			return common.ErrBadRequest("invalid page type")
		}
		return err0
	})
	if err != nil {
		return nil, err
	}

	return &space.AddPageLinkResp{
		PageType:  pageTypeStr,
		NewPageId: newPageId,
	}, nil
}
