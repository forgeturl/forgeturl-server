package api

import (
	"forgeturl-server/api/common"
	"forgeturl-server/api/login"
	"forgeturl-server/dal"
	"forgeturl-server/dal/model"
	"forgeturl-server/pkg/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/sunmi-OS/gocore/v2/api"
	"github.com/sunmi-OS/gocore/v2/glog"
)

type loginServiceImpl struct {
}

func NewLoginService() login.LoginServiceHTTPServer {
	return &loginServiceImpl{}
}

func LoginAuth(l login.LoginServiceHTTPServer) gin.HandlerFunc {
	return func(g *gin.Context) {
		req := &login.ConnectorReq{}
		apiCtx := api.NewContext(g)
		req.Provider = apiCtx.Query("provider")
		req.Code = apiCtx.Query("code")
		resp, err := connector(&apiCtx, req)
		apiCtx.RetJSON(resp, err)
	}
}

func LoginCallback(l login.LoginServiceHTTPServer) gin.HandlerFunc {
	return func(g *gin.Context) {
		req := &login.ConnectorCallbackReq{}
		apiCtx := api.NewContext(g)
		// Get provider name from path parameter
		req.Provider = apiCtx.Param("provider")
		apiCtx.Request.SetPathValue("provider", req.Provider)
		resp, err := connectorCallback(&apiCtx, req)
		apiCtx.RetJSON(resp, err)
	}
}

// connector 连接器登录，跳转鉴权的url
func connector(context *api.Context, req *login.ConnectorReq) (*login.ConnectorResp, error) {
	providerName := req.Provider
	if providerName == "" {
		return nil, common.ErrNotAuthenticated("you must select a provider")
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

func connectorCallback(apiCtx *api.Context, req *login.ConnectorCallbackReq) (*login.ConnectorCallbackResp, error) {
	// 未来这里可以增加，上一次登录时间的记录
	ctx := apiCtx.Request.Context()
	user, err := gothic.CompleteUserAuth(apiCtx.Writer, apiCtx.Request)
	if err != nil {
		return nil, common.ErrNotAuthenticated(err.Error())
	}

	// 处理用户信息
	userInfo, err := dal.User.GetByExternalID(ctx, user.UserID)
	isNewUser := false
	if err != nil {
		if common.IsErrNotFound(err) {
			// 创建新用户
			now := time.Now()
			newUser := &model.User{
				DisplayName:   user.Name,
				Username:      user.NickName,
				Email:         user.Email,
				ExternalID:    user.UserID,
				Avatar:        user.AvatarURL,
				LastLoginDate: now,
			}
			err = dal.User.Create(ctx, newUser)
			if err != nil {
				return nil, err
			}
			userInfo = newUser
			isNewUser = true
		} else {
			return nil, err
		}
	}
	uid := userInfo.ID

	uuid := middleware.NewUUID()
	err = dal.C.SetXToken(ctx, uuid, userInfo.ID)
	if err != nil {
		return nil, common.ErrInternalServerError("set x-token failed")
	}
	lastLoginTime := time.Now()
	apiCtx.Writer.Header().Set("X-Token", uuid)
	// 更新登录时间
	err = dal.User.UpdateLastLoginTime(ctx, uid, lastLoginTime)
	if err != nil {
		glog.ErrorC(ctx, "update last login time failed, err:%s", err.Error())
		err = nil
	}

	return &login.ConnectorCallbackResp{
		Uid:         userInfo.ID,
		DisplayName: userInfo.DisplayName,
		Username:    userInfo.Username,
		Avatar:      userInfo.Avatar,
		Email:       userInfo.Email,
		IsNewUser:   isNewUser,
	}, nil
}

func (l loginServiceImpl) Logout(context *api.Context, req *login.LogoutReq) (*login.LogoutResp, error) {
	ctx := context.Request.Context()
	token := context.GetHeader("X-Token")
	err := dal.C.DelXToken(ctx, token)
	if err != nil {
		return nil, common.ErrInternalServerError("logout failed")
	}
	return &login.LogoutResp{}, nil
}
