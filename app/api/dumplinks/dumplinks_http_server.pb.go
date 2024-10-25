// Code generated by protoc-gen-go-gin. DO NOT EDIT.
// versions:
// - protoc-gen-go-gin v1.0.3
// - protoc            v4.24.2
// source: api/proto/dumplinks.proto

package dumplinks

import (
	gin "github.com/gin-gonic/gin"
	api "github.com/sunmi-OS/gocore/v2/api"
)

// DumplinksServiceHTTPServer is the server API for DumplinksService service.
type DumplinksServiceHTTPServer interface {
	ImportBookmarks(*api.Context, *ImportBookmarksReq) (*ImportBookmarksResp, error)
	ExportBookmarks(*api.Context, *ExportBookmarksReq) (*ExportBookmarksResp, error)
}

func RegisterDumplinksServiceHTTPServer(s *gin.Engine, srv DumplinksServiceHTTPServer) {
	r := s.Group("/")
	r.POST("/dumplinks/importBookmarks", _DumplinksService_ImportBookmarks_HTTP_Handler(srv))
	r.POST("/dumplinks/exportBookmarks", _DumplinksService_ExportBookmarks_HTTP_Handler(srv))
}

func _DumplinksService_ImportBookmarks_HTTP_Handler(srv DumplinksServiceHTTPServer) func(g *gin.Context) {
	return func(g *gin.Context) {
		req := &ImportBookmarksReq{}
		var err error
		ctx := api.NewContext(g)
		err = parseReq(&ctx, req)
		err = checkValidate(err)
		if err != nil {
			setRetJSON(&ctx, nil, err)
			return
		}
		resp, err := srv.ImportBookmarks(&ctx, req)
		setRetJSON(&ctx, resp, err)
	}
}

func _DumplinksService_ExportBookmarks_HTTP_Handler(srv DumplinksServiceHTTPServer) func(g *gin.Context) {
	return func(g *gin.Context) {
		req := &ExportBookmarksReq{}
		ctx := api.NewContext(g)
		resp, err := srv.ExportBookmarks(&ctx, req)
		setRetJSON(&ctx, resp, err)
	}
}
