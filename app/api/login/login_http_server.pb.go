// Code generated by protoc-gen-go-gin. DO NOT EDIT.
// versions:
// - protoc-gen-go-gin v1.0.3
// - protoc            v4.24.2
// source: api/proto/login.proto

package login

import (
	gin "github.com/gin-gonic/gin"
	api "github.com/sunmi-OS/gocore/v2/api"
)

// LoginServiceHTTPServer is the server API for LoginService service.
type LoginServiceHTTPServer interface {
	// 连接器登录，跳转鉴权的url
	// https://github.com/googleapis/googleapis/blob/master/google/api/http.proto
	Connector(*api.Context, *ConnectorReq) (*ConnectorResp, error)
	// 第三方登录回调
	ConnectorCallback(*api.Context, *ConnectorCallbackReq) (*ConnectorCallbackResp, error)
	GetUserInfo(*api.Context, *GetUserInfoReq) (*GetUserInfoResp, error)
}

func RegisterLoginServiceHTTPServer(s *gin.Engine, srv LoginServiceHTTPServer) {
	r := s.Group("/")
	r.GET("/login/connector/auth/:name", _LoginService_Connector_HTTP_Handler(srv))             // 连接器登录，跳转鉴权的url
	r.GET("/login/connector/callback/:name", _LoginService_ConnectorCallback_HTTP_Handler(srv)) // 第三方登录回调
	r.POST("/page/getUserInfo", _LoginService_GetUserInfo_HTTP_Handler(srv))
}

func _LoginService_Connector_HTTP_Handler(srv LoginServiceHTTPServer) func(g *gin.Context) {
	return func(g *gin.Context) {
		req := &ConnectorReq{}
		var err error
		ctx := api.NewContext(g)
		err = ctx.ShouldBindUri(req)
		err = checkValidate(err)
		if err != nil {
			setRetJSON(&ctx, nil, err)
			return
		}
		err = parseReq(&ctx, req)
		err = checkValidate(err)
		if err != nil {
			setRetJSON(&ctx, nil, err)
			return
		}
		resp, err := srv.Connector(&ctx, req)
		setRetJSON(&ctx, resp, err)
	}
}

func _LoginService_ConnectorCallback_HTTP_Handler(srv LoginServiceHTTPServer) func(g *gin.Context) {
	return func(g *gin.Context) {
		req := &ConnectorCallbackReq{}
		var err error
		ctx := api.NewContext(g)
		err = ctx.ShouldBindUri(req)
		err = checkValidate(err)
		if err != nil {
			setRetJSON(&ctx, nil, err)
			return
		}
		err = parseReq(&ctx, req)
		err = checkValidate(err)
		if err != nil {
			setRetJSON(&ctx, nil, err)
			return
		}
		resp, err := srv.ConnectorCallback(&ctx, req)
		setRetJSON(&ctx, resp, err)
	}
}

func _LoginService_GetUserInfo_HTTP_Handler(srv LoginServiceHTTPServer) func(g *gin.Context) {
	return func(g *gin.Context) {
		req := &GetUserInfoReq{}
		var err error
		ctx := api.NewContext(g)
		err = parseReq(&ctx, req)
		err = checkValidate(err)
		if err != nil {
			setRetJSON(&ctx, nil, err)
			return
		}
		resp, err := srv.GetUserInfo(&ctx, req)
		setRetJSON(&ctx, resp, err)
	}
}
