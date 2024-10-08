package api

import (
	"forgeturl-server/api/login"
	"github.com/sunmi-OS/gocore/v2/api"
)

type loginSerivceImpl struct {
}

func NewLoginService() login.LoginServiceHTTPServer {
	return &loginSerivceImpl{}
}

func (l loginSerivceImpl) Connector(context *api.Context, req *login.ConnectorReq) (*login.ConnectorResp, error) {
	ctx := context.Request.Context()
	_ = ctx
	//TODO implement me
	panic("implement me")
}
