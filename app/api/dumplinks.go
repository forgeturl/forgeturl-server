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

type BookmarkContent struct {
	Folders []*dumplinks.Folder `json:"folders"`
}

func (d dumplinksSerivceImpl) ExportBookmarks(ctx *api.Context, req *dumplinks.ExportBookmarksReq) (*dumplinks.ExportBookmarksResp, error) {
	// Get user ID from context
	//userID, exists := ctx.Get("uid")
	//if !exists {
	//	return nil, common.ErrBadRequest()
	//}
	//uid, ok := userID.(int64)
	//if !ok {
	//	return nil, common.ErrBadRequest()
	//}
	//
	//// Query the database for the user's bookmarks
	//q := query.Use(query.DB)
	//page, err := q.Page.WithContext(ctx.Request.Context()).
	//	Where(q.Page.UID.Eq(uid), q.Page.Title.Eq("Chrome Bookmarks")).
	//	First()
	//
	//if err != nil {
	//	if err == gorm.ErrRecordNotFound {
	//		// No bookmarks found, return empty response
	//		return &dumplinks.ExportBookmarksResp{
	//			Folders: []*dumplinks.Folder{},
	//		}, nil
	//	}
	//	return nil, ecode.ErrSystem
	//}

	// Parse the content JSON
	var bookmarkContent BookmarkContent
	// if err := json.Unmarshal([]byte(page.Content), &bookmarkContent); err != nil {
	//	return nil, ecode.ErrSystem
	//}

	return &dumplinks.ExportBookmarksResp{
		Folders: bookmarkContent.Folders,
	}, nil
}
