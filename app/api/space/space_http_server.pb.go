// Code generated by protoc-gen-go-gin. DO NOT EDIT.
// versions:
// - protoc-gen-go-gin v1.0.3
// - protoc            v4.24.2
// source: api/proto/space.proto

package space

import (
	gin "github.com/gin-gonic/gin"
	api "github.com/sunmi-OS/gocore/v2/api"
)

// SpaceServiceHTTPServer is the server API for SpaceService service.
type SpaceServiceHTTPServer interface {
	// 拉取我的空间 || 空间
	// 登录状态才能拉到自己的空间
	// 部分页面如果消失或者没权限了，需要自动移除
	GetMySpace(*api.Context, *GetMySpaceReq) (*GetMySpaceResp, error)
	// 调整我的空间下面的页面顺序 || 空间
	ChangeSpacePageSequence(*api.Context, *ChangeSpacePageSequenceReq) (*ChangeSpacePageSequenceResp, error)
	// 创建临时页面 || 页面
	// 非登录状态可以创建临时页面，默认一个浏览器只能创建一个自己的临时页面
	// 创建完成后得到一个随机页面id(比如 240626-abcd)，不使用lo等字符串，只使用其他字母
	// 生成算法：当前时间转换的4个字母(时分秒)
	CreateTmpPage(*api.Context, *CreateTmpPageReq) (*CreateTmpPageResp, error)
	// 拉取某个页面数据 || 页面
	// 拉取某个页面
	// 临时页面，可以读到
	GetPages(*api.Context, *GetPagesReq) (*GetPagesResp, error)
	// 拉取某个页面数据 || 页面
	GetPage(*api.Context, *GetPageReq) (*GetPageResp, error)
	// 更新页面 || 页面
	UpdatePage(*api.Context, *UpdatePageReq) (*UpdatePageResp, error)
	// 把整个页面删除 || 页面
	// 自己的默认页面只能清空，无法删除
	DeletePage(*api.Context, *DeletePageReq) (*DeletePageResp, error)
	// 去除页面连接 || 页面
	// 把页面的只读链接删除
	RemovePageLink(*api.Context, *RemovePageLinkReq) (*RemovePageLinkResp, error)
	// 生成新页面链接 || 页面
	CreateNewPageLink(*api.Context, *CreateNewPageLinkReq) (*CreateNewPageLinkResp, error)
}

func RegisterSpaceServiceHTTPServer(s *gin.Engine, srv SpaceServiceHTTPServer) {
	r := s.Group("/")
	r.POST("/page/getMySpace", _SpaceService_GetMySpace_HTTP_Handler(srv))                           // 拉取我的空间 || 空间
	r.POST("/page/changeSpacePageSequence", _SpaceService_ChangeSpacePageSequence_HTTP_Handler(srv)) // 调整我的空间下面的页面顺序 || 空间
	r.POST("/page/createTmpPage", _SpaceService_CreateTmpPage_HTTP_Handler(srv))                     // 创建临时页面 || 页面
	r.POST("/page/getPages", _SpaceService_GetPages_HTTP_Handler(srv))                               // 拉取某个页面数据 || 页面
	r.POST("/page/getPage", _SpaceService_GetPage_HTTP_Handler(srv))                                 // 拉取某个页面数据 || 页面
	r.POST("/page/updatePage", _SpaceService_UpdatePage_HTTP_Handler(srv))                           // 更新页面 || 页面
	r.POST("/page/deletePage", _SpaceService_DeletePage_HTTP_Handler(srv))                           // 把整个页面删除 || 页面
	r.POST("/page/removePageLink", _SpaceService_RemovePageLink_HTTP_Handler(srv))                   // 去除页面连接 || 页面
	r.POST("/page/createNewPageLink", _SpaceService_CreateNewPageLink_HTTP_Handler(srv))             // 生成新页面链接 || 页面
}

func _SpaceService_GetMySpace_HTTP_Handler(srv SpaceServiceHTTPServer) func(g *gin.Context) {
	return func(g *gin.Context) {
		req := &GetMySpaceReq{}
		var err error
		ctx := api.NewContext(g)
		err = parseReq(&ctx, req)
		err = checkValidate(err)
		if err != nil {
			setRetJSON(&ctx, nil, err)
			return
		}
		resp, err := srv.GetMySpace(&ctx, req)
		setRetJSON(&ctx, resp, err)
	}
}

func _SpaceService_ChangeSpacePageSequence_HTTP_Handler(srv SpaceServiceHTTPServer) func(g *gin.Context) {
	return func(g *gin.Context) {
		req := &ChangeSpacePageSequenceReq{}
		var err error
		ctx := api.NewContext(g)
		err = parseReq(&ctx, req)
		err = checkValidate(err)
		if err != nil {
			setRetJSON(&ctx, nil, err)
			return
		}
		resp, err := srv.ChangeSpacePageSequence(&ctx, req)
		setRetJSON(&ctx, resp, err)
	}
}

func _SpaceService_CreateTmpPage_HTTP_Handler(srv SpaceServiceHTTPServer) func(g *gin.Context) {
	return func(g *gin.Context) {
		req := &CreateTmpPageReq{}
		ctx := api.NewContext(g)
		resp, err := srv.CreateTmpPage(&ctx, req)
		setRetJSON(&ctx, resp, err)
	}
}

func _SpaceService_GetPages_HTTP_Handler(srv SpaceServiceHTTPServer) func(g *gin.Context) {
	return func(g *gin.Context) {
		req := &GetPagesReq{}
		var err error
		ctx := api.NewContext(g)
		err = parseReq(&ctx, req)
		err = checkValidate(err)
		if err != nil {
			setRetJSON(&ctx, nil, err)
			return
		}
		resp, err := srv.GetPages(&ctx, req)
		setRetJSON(&ctx, resp, err)
	}
}

func _SpaceService_GetPage_HTTP_Handler(srv SpaceServiceHTTPServer) func(g *gin.Context) {
	return func(g *gin.Context) {
		req := &GetPageReq{}
		var err error
		ctx := api.NewContext(g)
		err = parseReq(&ctx, req)
		err = checkValidate(err)
		if err != nil {
			setRetJSON(&ctx, nil, err)
			return
		}
		resp, err := srv.GetPage(&ctx, req)
		setRetJSON(&ctx, resp, err)
	}
}

func _SpaceService_UpdatePage_HTTP_Handler(srv SpaceServiceHTTPServer) func(g *gin.Context) {
	return func(g *gin.Context) {
		req := &UpdatePageReq{}
		var err error
		ctx := api.NewContext(g)
		err = parseReq(&ctx, req)
		err = checkValidate(err)
		if err != nil {
			setRetJSON(&ctx, nil, err)
			return
		}
		resp, err := srv.UpdatePage(&ctx, req)
		setRetJSON(&ctx, resp, err)
	}
}

func _SpaceService_DeletePage_HTTP_Handler(srv SpaceServiceHTTPServer) func(g *gin.Context) {
	return func(g *gin.Context) {
		req := &DeletePageReq{}
		var err error
		ctx := api.NewContext(g)
		err = parseReq(&ctx, req)
		err = checkValidate(err)
		if err != nil {
			setRetJSON(&ctx, nil, err)
			return
		}
		resp, err := srv.DeletePage(&ctx, req)
		setRetJSON(&ctx, resp, err)
	}
}

func _SpaceService_RemovePageLink_HTTP_Handler(srv SpaceServiceHTTPServer) func(g *gin.Context) {
	return func(g *gin.Context) {
		req := &RemovePageLinkReq{}
		var err error
		ctx := api.NewContext(g)
		err = parseReq(&ctx, req)
		err = checkValidate(err)
		if err != nil {
			setRetJSON(&ctx, nil, err)
			return
		}
		resp, err := srv.RemovePageLink(&ctx, req)
		setRetJSON(&ctx, resp, err)
	}
}

func _SpaceService_CreateNewPageLink_HTTP_Handler(srv SpaceServiceHTTPServer) func(g *gin.Context) {
	return func(g *gin.Context) {
		req := &CreateNewPageLinkReq{}
		var err error
		ctx := api.NewContext(g)
		err = parseReq(&ctx, req)
		err = checkValidate(err)
		if err != nil {
			setRetJSON(&ctx, nil, err)
			return
		}
		resp, err := srv.CreateNewPageLink(&ctx, req)
		setRetJSON(&ctx, resp, err)
	}
}