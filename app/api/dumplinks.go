package api

import (
	"encoding/json"

	"forgeturl-server/api/common"
	"forgeturl-server/api/dumplinks"
	"forgeturl-server/dal"
	"forgeturl-server/dal/model"

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

func (d dumplinksSerivceImpl) ImportBookmarks(ctx *api.Context, req *dumplinks.ImportBookmarksReq) (*dumplinks.ImportBookmarksResp, error) {
	// Get user ID from context
	userID, exists := ctx.Get("uid")
	if !exists {
		return nil, common.ErrForbidden()
	}
	uid, ok := userID.(int64)
	if !ok {
		return nil, common.ErrForbidden()
	}

	// Convert bookmarks to JSON content
	bookmarkContent := &BookmarkContent{
		Folders: req.Folders,
	}
	content, err := json.Marshal(bookmarkContent)
	if err != nil {
		return nil, common.ErrBadRequest()
	}

	// Create or update page record
	var page *model.Page

	// Check if user already has a bookmark page
	// 如果存在需要删除先
	_, err = dal.Page.GetSelfPage(ctx.Request.Context(), uid)
	if common.IsErrNotFound(err) {
		// Create new page for bookmarks with proper page IDs
		pid := genOwnerPageId()
		dal.UniquePid.Create(ctx, pid)
		page = &model.Page{
			UID:     uid,
			Title:   "Chrome Bookmarks",
			Content: string(content),
			Pid:     pid,
			//ReadonlyPid: genReadOnlyPageId(),
			//EditPid:     genEditPageId(),
			//AdminPid:    genOwnerPageId(), // Reusing owner page ID for admin
		}
		err = dal.Page.Create(ctx.Request.Context(), page)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, common.ErrBadRequest("Your self page already exist, please remove them first")
	}

	return &dumplinks.ImportBookmarksResp{}, nil
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
