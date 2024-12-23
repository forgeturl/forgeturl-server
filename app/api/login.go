package api

import (
	"forgeturl-server/api/common"
	"forgeturl-server/api/login"
	"forgeturl-server/dal"
	"forgeturl-server/dal/model"
	connector_google "forgeturl-server/pkg/connector-google"
	connector_weixin "forgeturl-server/pkg/connector-weixin"
	"github.com/sunmi-OS/gocore/v2/api"
)

type loginSerivceImpl struct {
}

func NewLoginService() login.LoginServiceHTTPServer {
	return &loginSerivceImpl{}
}

func (l loginSerivceImpl) Connector(context *api.Context, req *login.ConnectorReq) (*login.ConnectorResp, error) {
	ctx := context.Request.Context()

	externalId := ""
	username := ""
	email := ""
	displayName := ""
	avatar := ""
	switch req.Name {
	case "google":
		uInfo, err := connector_google.Connector.ConnectorReceiver(ctx, req.Code)
		if err != nil {
			return nil, err
		}
		username = uInfo.Name
		displayName = uInfo.Name
		externalId = uInfo.Sub
		email = uInfo.Email
		avatar = uInfo.Picture
	case "facebook":
	case "weixin":
		connector_weixin.Init()
	default:
		return nil, common.ErrNotAuthenticated("invalid name")
	}

	userInfo, err := dal.User.GetByExternalID(ctx, externalId)
	if err != nil {
		// 如果是不存在，则找不到
		if common.IsErrNotFound(err) {
			user := &model.User{
				DisplayName: displayName,
				Username:    username,
				Email:       email,
				ExternalID:  externalId,
				Avatar:      avatar,
			}
			err = dal.User.Create(ctx, user)
			if err != nil {
				return nil, err
			}
			userInfo = user
		} else {
			return nil, err
		}
		return nil, err
	}
	_ = userInfo
	return &login.ConnectorResp{}, nil
}
