package api

import (
	"forgeturl-server/api/dumplinks"

	"github.com/sunmi-OS/gocore/v2/api"
)

type dumplinksSerivceImpl struct {
}

func NewDumplinksService() dumplinks.DumplinksServiceHTTPServer {
	return &dumplinksSerivceImpl{}
}

func (d dumplinksSerivceImpl) ImportBookmarks(context *api.Context, req *dumplinks.ImportBookmarksReq) (*dumplinks.ImportBookmarksResp, error) {

	//TODO implement me
	panic("implement me")
}

func (d dumplinksSerivceImpl) ExportBookmarks(context *api.Context, req *dumplinks.ExportBookmarksReq) (*dumplinks.ExportBookmarksResp, error) {
	//TODO implement me
	panic("implement me")
}
