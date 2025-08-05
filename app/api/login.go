package api

import (
	"forgeturl-server/api/common"
	"forgeturl-server/api/login"
	"forgeturl-server/dal"
	"forgeturl-server/dal/model"
	"forgeturl-server/pkg/middleware"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/sunmi-OS/gocore/v2/api"
)

type loginServiceImpl struct {
}

func NewLoginService() login.LoginServiceHTTPServer {
	return &loginServiceImpl{}
}

// Connector 连接器登录，跳转鉴权的url
func (l loginServiceImpl) Connector(context *api.Context, req *login.ConnectorReq) (*login.ConnectorResp, error) {
	providerName, err := gothic.GetProviderName(context.Request)
	if err != nil {
		return nil, common.ErrNotAuthenticated(err.Error())
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return nil, common.ErrNotAuthenticated("invalid provider name")
	}

	sess, err := provider.BeginAuth(gothic.SetState(context.Request))
	if err != nil {
		return nil, err
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		return nil, err
	}

	err = gothic.StoreInSession(providerName, sess.Marshal(), context.Request, context.Writer)
	if err != nil {
		return nil, err
	}

	return &login.ConnectorResp{
		AuthUrl: url,
	}, nil
}

func (l loginServiceImpl) ConnectorCallback(context *api.Context, req *login.ConnectorCallbackReq) (*login.ConnectorCallbackResp, error) {
	ctx := context.Request.Context()
	providerName, err := gothic.GetProviderName(context.Request)
	if err != nil {
		return nil, common.ErrNotAuthenticated(err.Error())
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return nil, common.ErrNotAuthenticated("invalid provider name")
	}

	value, err := gothic.GetFromSession(providerName, context.Request)
	if err != nil {
		return nil, err
	}
	defer gothic.Logout(context.Writer, context.Request)

	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		return nil, err
	}

	params := context.Request.URL.Query()
	if params.Encode() == "" && context.Request.Method == "POST" {
		context.Request.ParseForm()
		params = context.Request.Form
	}

	_, err = sess.Authorize(provider, params)
	if err != nil {
		return nil, err
	}

	// 这已经偷偷塞到cookie里中去了
	err = gothic.StoreInSession(providerName, sess.Marshal(), context.Request, context.Writer)
	if err != nil {
		return nil, err
	}

	user, err := provider.FetchUser(sess)
	if err != nil {
		return nil, err
	}

	// 处理用户信息
	userInfo, err := dal.User.GetByExternalID(ctx, user.UserID)
	if err != nil {
		if common.IsErrNotFound(err) {
			// 创建新用户
			newUser := &model.User{
				DisplayName: user.Name,
				Username:    user.NickName,
				Email:       user.Email,
				ExternalID:  user.UserID,
				Avatar:      user.AvatarURL,
			}
			err = dal.User.Create(ctx, newUser)
			if err != nil {
				return nil, err
			}
			userInfo = newUser
		} else {
			return nil, err
		}
	}

	return &login.ConnectorCallbackResp{
		Uid:         userInfo.ID,
		DisplayName: userInfo.DisplayName,
		Username:    userInfo.Username,
		Avatar:      userInfo.Avatar,
		Email:       userInfo.Email,
	}, nil
}

func (l loginServiceImpl) GetUserInfo(context *api.Context, req *login.GetUserInfoReq) (*login.GetUserInfoResp, error) {
	// 已登录才能获取到详情，否则拉拉取不到
	ctx := context.Request.Context()
	uid := req.Uid
	loginUid := middleware.GetLoginUid(context)
	userInfo, err := dal.User.Get(ctx, uid)
	if err != nil {
		return nil, err
	}

	// 如果是自己，则返回内容多些
	info := &login.GetUserInfoResp{
		Uid:           uid,
		DisplayName:   userInfo.DisplayName,
		Username:      userInfo.Username,
		Avatar:        userInfo.Avatar,
		Email:         userInfo.Email,
		Status:        userInfo.Status,
		LastLoginTime: userInfo.LastLoginDate.Unix(),
		IsAdmin:       userInfo.IsAdmin,
		Provider:      userInfo.Provider,
		CreateTime:    userInfo.CreatedAt.Unix(),
		UpdateTime:    userInfo.UpdatedAt.Unix(),
	}

	// 如果不是自己，则返回内容少些
	if loginUid != uid {
		info.Email = ""
		info.Provider = ""
		info.Username = ""
		info.LastLoginTime = 0
		info.IsAdmin = 0
		info.CreateTime = 0
		info.UpdateTime = 0
	}

	return info, nil
}
