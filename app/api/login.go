package api

import (
	"forgeturl-server/api/common"
	"forgeturl-server/api/login"
	"forgeturl-server/dal"
	connector_google "forgeturl-server/pkg/connector-google"
	"github.com/sunmi-OS/gocore/v2/api"
)

type loginSerivceImpl struct {
}

func NewLoginService() login.LoginServiceHTTPServer {
	return &loginSerivceImpl{}
}

func (l loginSerivceImpl) Connector(context *api.Context, req *login.ConnectorReq) (*login.ConnectorResp, error) {
	ctx := context.Request.Context()

	name := ""
	externalId := ""
	switch req.Name {
	case "google":
		userInfo, err := connector_google.Connector.ConnectorReceiver(ctx, req.Code)
		if err != nil {
			return nil, err
		}
		name = userInfo.Name
		externalId = userInfo.Sub
	case "facebook":
	case "weixin":
	default:
		return nil, common.ErrNotAuthenticated("invalid name")

	}

	userInfo, err := dal.User.GetByExternalID(ctx, externalId)
	if err != nil {
		// 如果是不存在，则找不到
		if common.IsErrNotFound(err) {

		}
		return nil, err
	}

	_ = name

	//TODO implement me
	panic("implement me")
}
