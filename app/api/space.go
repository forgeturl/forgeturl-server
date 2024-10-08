package api

import (
	"forgeturl-server/api/space"
	"github.com/sunmi-OS/gocore/v2/api"
)

type spaceServiceImpl struct {
}

func NewSpaceService() space.SpaceServiceHTTPServer {
	return &spaceServiceImpl{}
}

func (s spaceServiceImpl) GetMySpace(context *api.Context, req *space.GetMySpaceReq) (*space.GetMySpaceResp, error) {
	ctx := context.Request.Context()
	_ = ctx
	//TODO implement me
	panic("implement me")
}

func (s spaceServiceImpl) ChangeSpacePageSequence(context *api.Context, req *space.ChangeSpacePageSequenceReq) (*space.ChangeSpacePageSequenceResp, error) {
	ctx := context.Request.Context()
	_ = ctx
	//TODO implement me
	panic("implement me")
}

func (s spaceServiceImpl) CreateTmpPage(context *api.Context, req *space.CreateTmpPageReq) (*space.CreateTmpPageResp, error) {
	ctx := context.Request.Context()
	_ = ctx
	//TODO implement me
	panic("implement me")
}

func (s spaceServiceImpl) GetPages(context *api.Context, req *space.GetPagesReq) (*space.GetPagesResp, error) {
	ctx := context.Request.Context()
	_ = ctx
	//TODO implement me
	panic("implement me")
}

func (s spaceServiceImpl) GetPage(context *api.Context, req *space.GetPageReq) (*space.GetPageResp, error) {
	ctx := context.Request.Context()
	_ = ctx
	//TODO implement me
	panic("implement me")
}

func (s spaceServiceImpl) UpdatePage(context *api.Context, req *space.UpdatePageReq) (*space.UpdatePageResp, error) {
	ctx := context.Request.Context()
	_ = ctx
	//TODO implement me
	panic("implement me")
}

func (s spaceServiceImpl) DeletePage(context *api.Context, req *space.DeletePageReq) (*space.DeletePageResp, error) {
	ctx := context.Request.Context()
	_ = ctx
	//TODO implement me
	panic("implement me")
}

func (s spaceServiceImpl) RemovePageLink(context *api.Context, req *space.RemovePageLinkReq) (*space.RemovePageLinkResp, error) {
	ctx := context.Request.Context()
	_ = ctx
	//TODO implement me
	panic("implement me")
}

func (s spaceServiceImpl) CreateNewPageLink(context *api.Context, req *space.CreateNewPageLinkReq) (*space.CreateNewPageLinkResp, error) {
	ctx := context.Request.Context()
	_ = ctx
	//TODO implement me
	panic("implement me")
}
