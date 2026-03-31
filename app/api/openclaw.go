package api

import (
	"forgeturl-server/api/common"
	"forgeturl-server/api/openclaw"
	"forgeturl-server/api/space"
	"forgeturl-server/dal"
	"forgeturl-server/pkg/middleware"
	"strings"

	"github.com/bytedance/sonic"
	"gorm.io/gorm"

	gapi "github.com/sunmi-OS/gocore/v2/api"
)

type openClawServiceImpl struct{}

func NewOpenClawService() openclaw.OpenClawServiceHTTPServer {
	return &openClawServiceImpl{}
}

func (s openClawServiceImpl) GetApiKey(context *gapi.Context, req *openclaw.GetApiKeyReq) (*openclaw.GetApiKeyResp, error) {
	ctx := context.Request.Context()
	uid := middleware.GetLoginUid(context)
	if uid == 0 {
		return nil, common.ErrNeedLogin("")
	}
	info, err := dal.OpenClaw.GetAPIKeyByUID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return &openclaw.GetApiKeyResp{
			HasKey: false,
			ApiKey: "",
		}, nil
	}
	return &openclaw.GetApiKeyResp{
		HasKey: true,
		ApiKey: info.APIKey,
	}, nil
}

func (s openClawServiceImpl) RegenerateApiKey(context *gapi.Context, req *openclaw.RegenerateApiKeyReq) (*openclaw.RegenerateApiKeyResp, error) {
	ctx := context.Request.Context()
	uid := middleware.GetLoginUid(context)
	if uid == 0 {
		return nil, common.ErrNeedLogin("")
	}
	info, err := dal.OpenClaw.RegenerateAPIKey(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &openclaw.RegenerateApiKeyResp{ApiKey: info.APIKey}, nil
}

func (s openClawServiceImpl) AddTmpBookmark(context *gapi.Context, req *openclaw.AddTmpBookmarkReq) (*openclaw.AddTmpBookmarkResp, error) {
	ctx := context.Request.Context()
	uid := middleware.GetLoginUid(context)
	if uid == 0 {
		return nil, common.ErrNeedLogin("")
	}
	info, err := dal.OpenClaw.UpsertTmpBookmark(ctx, uid, req.Title, req.Url)
	if err != nil {
		return nil, common.ErrBadRequest(err.Error())
	}
	return &openclaw.AddTmpBookmarkResp{Bookmark: toTmpBookmarkDTO(info)}, nil
}

func (s openClawServiceImpl) ListTmpBookmarks(context *gapi.Context, req *openclaw.ListTmpBookmarksReq) (*openclaw.ListTmpBookmarksResp, error) {
	ctx := context.Request.Context()
	uid := middleware.GetLoginUid(context)
	if uid == 0 {
		return nil, common.ErrNeedLogin("")
	}
	list, err := dal.OpenClaw.ListTmpBookmarks(ctx, uid, int(req.Limit))
	if err != nil {
		return nil, err
	}
	bookmarks := make([]*openclaw.TmpBookmark, 0, len(list))
	for _, item := range list {
		bookmarks = append(bookmarks, toTmpBookmarkDTO(item))
	}
	return &openclaw.ListTmpBookmarksResp{Bookmarks: bookmarks}, nil
}

func (s openClawServiceImpl) ExistsTmpBookmark(context *gapi.Context, req *openclaw.ExistsTmpBookmarkReq) (*openclaw.ExistsTmpBookmarkResp, error) {
	ctx := context.Request.Context()
	uid := middleware.GetLoginUid(context)
	if uid == 0 {
		return nil, common.ErrNeedLogin("")
	}
	info, exists, err := dal.OpenClaw.ExistsTmpBookmark(ctx, uid, req.Title, req.Url)
	if err != nil {
		return nil, common.ErrBadRequest(err.Error())
	}
	resp := &openclaw.ExistsTmpBookmarkResp{Exists: exists}
	if exists {
		resp.Bookmark = toTmpBookmarkDTO(info)
	}
	return resp, nil
}

func (s openClawServiceImpl) MoveTmpBookmarkToPage(context *gapi.Context, req *openclaw.MoveTmpBookmarkToPageReq) (*openclaw.EmptyResp, error) {
	ctx := context.Request.Context()
	uid := middleware.GetLoginUid(context)
	if uid == 0 {
		return nil, common.ErrNeedLogin("")
	}

	userInfo, err := dal.User.Get(ctx, uid)
	if err != nil {
		return nil, err
	}
	ownerPid, err := canEditPage(ctx, userInfo, req.PageId)
	if err != nil {
		return nil, err
	}

	err = dal.OpenClaw.Transaction(ctx, func(tx *gorm.DB) error {
		bookmark, err0 := dal.OpenClaw.GetTmpBookmarkByID(ctx, uid, req.BookmarkId, tx)
		if err0 != nil {
			return err0
		}
		page, err0 := dal.OpenClaw.LockOwnerPage(ctx, ownerPid, tx)
		if err0 != nil {
			return err0
		}

		collections, err0 := parseCollections(page.Content)
		if err0 != nil {
			return common.ErrInternalServerError(err0.Error())
		}

		appendBookmarkToCollections(&collections, bookmark.Title, bookmark.URL, req.CollectionTitle)
		newContent, err0 := sonic.MarshalString(collections)
		if err0 != nil {
			return common.ErrInternalServerError(err0.Error())
		}

		if err0 = dal.OpenClaw.UpdateOwnerPageContentByVersion(ctx, ownerPid, page.Version, newContent, tx); err0 != nil {
			return err0
		}
		if err0 = dal.OpenClaw.DeleteTmpBookmarkByID(ctx, uid, req.BookmarkId, tx); err0 != nil {
			return err0
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &openclaw.EmptyResp{}, nil
}

func (s openClawServiceImpl) MovePageBookmarkToTmp(context *gapi.Context, req *openclaw.MovePageBookmarkToTmpReq) (*openclaw.EmptyResp, error) {
	ctx := context.Request.Context()
	uid := middleware.GetLoginUid(context)
	if uid == 0 {
		return nil, common.ErrNeedLogin("")
	}

	userInfo, err := dal.User.Get(ctx, uid)
	if err != nil {
		return nil, err
	}
	ownerPid, err := canEditPage(ctx, userInfo, req.PageId)
	if err != nil {
		return nil, err
	}

	err = dal.OpenClaw.Transaction(ctx, func(tx *gorm.DB) error {
		page, err0 := dal.OpenClaw.LockOwnerPage(ctx, ownerPid, tx)
		if err0 != nil {
			return err0
		}

		collections, err0 := parseCollections(page.Content)
		if err0 != nil {
			return common.ErrInternalServerError(err0.Error())
		}

		removed := removeBookmarkFromCollections(&collections, req.CollectionTitle, req.Title, req.Url)
		if !removed {
			return common.ErrNotFound("bookmark not found in page collections")
		}

		newContent, err0 := sonic.MarshalString(collections)
		if err0 != nil {
			return common.ErrInternalServerError(err0.Error())
		}

		if err0 = dal.OpenClaw.UpdateOwnerPageContentByVersion(ctx, ownerPid, page.Version, newContent, tx); err0 != nil {
			return err0
		}
		if _, err0 = dal.OpenClaw.UpsertTmpBookmark(ctx, uid, req.Title, req.Url, tx); err0 != nil {
			return common.ErrBadRequest(err0.Error())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &openclaw.EmptyResp{}, nil
}

func toTmpBookmarkDTO(info *dal.TmpBookmark) *openclaw.TmpBookmark {
	if info == nil {
		return nil
	}
	resp := &openclaw.TmpBookmark{
		Id:    info.ID,
		Title: info.Title,
		Url:   info.URL,
	}
	if info.CreatedAt != nil {
		resp.CreateAt = info.CreatedAt.Unix()
	}
	if info.UpdatedAt != nil {
		resp.UpdateAt = info.UpdatedAt.Unix()
	}
	return resp
}

func parseCollections(content string) ([]*space.Collections, error) {
	collections := make([]*space.Collections, 0)
	if strings.TrimSpace(content) == "" {
		return collections, nil
	}
	if err := sonic.UnmarshalString(content, &collections); err != nil {
		return nil, err
	}
	return collections, nil
}

func appendBookmarkToCollections(collections *[]*space.Collections, title, rawURL, preferredCollectionTitle string) {
	co := *collections
	title = strings.TrimSpace(title)
	rawURL = strings.TrimSpace(rawURL)
	targetTitle := strings.TrimSpace(preferredCollectionTitle)
	if targetTitle == "" {
		targetTitle = "Temporary Inbox"
	}

	targetIdx := -1
	for i, c := range co {
		if strings.TrimSpace(c.Title) == targetTitle {
			targetIdx = i
			break
		}
	}

	if targetIdx == -1 {
		co = append(co, &space.Collections{
			Title: targetTitle,
			Links: make([]*space.Link, 0),
		})
		targetIdx = len(co) - 1
	}

	target := co[targetIdx]
	for _, l := range target.Links {
		if sameBookmark(l.Title, l.Url, title, rawURL) {
			*collections = co
			return
		}
	}
	target.Links = append(target.Links, &space.Link{
		Title: title,
		Url:   rawURL,
	})
	co[targetIdx] = target
	*collections = co
}

func removeBookmarkFromCollections(collections *[]*space.Collections, collectionTitle, title, rawURL string) bool {
	co := *collections
	title = strings.TrimSpace(title)
	rawURL = strings.TrimSpace(rawURL)
	collectionTitle = strings.TrimSpace(collectionTitle)

	for i, c := range co {
		if collectionTitle != "" && strings.TrimSpace(c.Title) != collectionTitle {
			continue
		}
		for linkIdx, l := range c.Links {
			if !sameBookmark(l.Title, l.Url, title, rawURL) {
				continue
			}
			c.Links = append(c.Links[:linkIdx], c.Links[linkIdx+1:]...)
			co[i] = c
			*collections = co
			return true
		}
	}
	return false
}

func sameBookmark(leftTitle, leftURL, rightTitle, rightURL string) bool {
	leftTitle = strings.TrimSpace(leftTitle)
	rightTitle = strings.TrimSpace(rightTitle)
	if leftTitle != rightTitle {
		return false
	}

	leftNorm, leftErr := dal.NormalizeBookmarkURL(leftURL)
	rightNorm, rightErr := dal.NormalizeBookmarkURL(rightURL)
	if leftErr == nil && rightErr == nil {
		return leftNorm == rightNorm
	}
	return strings.TrimSpace(leftURL) == strings.TrimSpace(rightURL)
}
