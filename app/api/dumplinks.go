package api

import (
	"encoding/json"
	"forgeturl-server/api/dumplinks"
	"forgeturl-server/dal/model"
	"forgeturl-server/dal/query"

	"github.com/sunmi-OS/gocore/v2/api"
	"github.com/sunmi-OS/gocore/v2/api/ecode"
	"gorm.io/gorm"
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
		return nil, ecode.ErrUnauthorized()
	}
	uid, ok := userID.(int64)
	if !ok {
		return nil, ecode.ErrUnauthorized
	}

	// Convert bookmarks to JSON content
	bookmarkContent := &BookmarkContent{
		Folders: req.Folders,
	}
	content, err := json.Marshal(bookmarkContent)
	if err != nil {
		return nil, ecode.ErrInvalidParam
	}

	// Create or update page record
	q := query.Use(query.DB)
	var page *model.Page

	// Check if user already has a bookmark page
	existingPage, err := q.Page.WithContext(ctx.Request.Context()).
		Where(q.Page.UID.Eq(uid), q.Page.Title.Eq("Chrome Bookmarks")).
		First()

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new page for bookmarks with proper page IDs
			page = &model.Page{
				UID:         uid,
				Title:       "Chrome Bookmarks",
				Content:     string(content),
				Pid:         genOwnerPageId(),
				ReadonlyPid: genReadOnlyPageId(),
				EditPid:     genEditPageId(),
				AdminPid:    genOwnerPageId(), // Reusing owner page ID for admin
			}
			err = q.Page.WithContext(ctx.Request.Context()).Create(page)
			if err != nil {
				return nil, ecode.ErrSystem
			}
		} else {
			return nil, ecode.ErrSystem
		}
	} else {
		// Update existing page
		_, err = q.Page.WithContext(ctx.Request.Context()).
			Where(q.Page.ID.Eq(existingPage.ID)).
			Update(q.Page.Content, string(content))
		if err != nil {
			return nil, ecode.ErrSystem
		}
	}

	return &dumplinks.ImportBookmarksResp{}, nil
}

func (d dumplinksSerivceImpl) ExportBookmarks(ctx *api.Context, req *dumplinks.ExportBookmarksReq) (*dumplinks.ExportBookmarksResp, error) {
	// Get user ID from context
	userID, exists := ctx.Get("uid")
	if !exists {
		return nil, ecode.ErrUnauthorized
	}
	uid, ok := userID.(int64)
	if !ok {
		return nil, ecode.ErrUnauthorized
	}

	// Query the database for the user's bookmarks
	q := query.Use(query.DB)
	page, err := q.Page.WithContext(ctx.Request.Context()).
		Where(q.Page.UID.Eq(uid), q.Page.Title.Eq("Chrome Bookmarks")).
		First()

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// No bookmarks found, return empty response
			return &dumplinks.ExportBookmarksResp{
				Folders: []*dumplinks.Folder{},
			}, nil
		}
		return nil, ecode.ErrSystem
	}

	// Parse the content JSON
	var bookmarkContent BookmarkContent
	if err := json.Unmarshal([]byte(page.Content), &bookmarkContent); err != nil {
		return nil, ecode.ErrSystem
	}

	return &dumplinks.ExportBookmarksResp{
		Folders: bookmarkContent.Folders,
	}, nil
}
