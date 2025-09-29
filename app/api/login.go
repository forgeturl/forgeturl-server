package api

import (
	"forgeturl-server/api/common"
	"forgeturl-server/api/login"
	"forgeturl-server/dal"
	"forgeturl-server/dal/model"
	"forgeturl-server/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/sunmi-OS/gocore/v2/api"
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
		resp, err := Connector(&apiCtx, req)
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
		resp, err := ConnectorCallback(&apiCtx, req)
		apiCtx.RetJSON(resp, err)
	}
}

// Connector 连接器登录，跳转鉴权的url
func Connector(context *api.Context, req *login.ConnectorReq) (*login.ConnectorResp, error) {
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

func ConnectorCallback(apiCtx *api.Context, req *login.ConnectorCallbackReq) (*login.ConnectorCallbackResp, error) {
	ctx := apiCtx.Request.Context()
	user, err := gothic.CompleteUserAuth(apiCtx.Writer, apiCtx.Request)
	if err != nil {
		return nil, common.ErrNotAuthenticated(err.Error())
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

	uuid := middleware.NewUUID()
	err = dal.C.SetXToken(ctx, middleware.NewUUID(), userInfo.ID)
	if err != nil {
		return nil, common.ErrInternalServerError("set x-token failed")
	}
	apiCtx.Writer.Header().Set("X-Token", uuid)

	return &login.ConnectorCallbackResp{
		Uid:         userInfo.ID,
		DisplayName: userInfo.DisplayName,
		Username:    userInfo.Username,
		Avatar:      userInfo.Avatar,
		Email:       userInfo.Email,
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
