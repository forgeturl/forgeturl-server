package api

import (
	"net/http/httptest"
	"testing"
	"time"

	"forgeturl-server/api/common"
	"forgeturl-server/api/space"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	gapi "github.com/sunmi-OS/gocore/v2/api"
)

func guestPageContext() *gapi.Context {
	g, _ := gin.CreateTestContext(httptest.NewRecorder())
	g.Request = httptest.NewRequest("POST", "/space/getPage", nil)
	c := gapi.NewContext(g)
	return &c
}

func TestGuestCanReadEditableShareWithoutOwningIt(t *testing.T) {
	mock := sharedPageDB(t)
	mock.ExpectQuery("SELECT .* FROM `page` ").WithArgs("Eabc", 1).WillReturnRows(sqlmock.NewRows(
		[]string{"uid", "pid", "edit_pid", "title", "content", "created_at", "updated_at"},
	).AddRow(7, "Oabc", "Eabc", "Shared", `[{"title":"Folder","links":[{"title":"Link","url":"https://example.com"}]}]`, time.Unix(100, 0), time.Unix(100, 0)))
	resp, err := (spaceServiceImpl{}).GetPage(guestPageContext(), &space.GetPageReq{PageId: "Eabc"})
	if err != nil {
		t.Fatal(err)
	}
	page := resp.Page
	if page.PageId != "Eabc" || page.Title != "Shared" || len(page.Collections) != 1 || page.Collections[0].Links[0].Title != "Link" {
		t.Fatalf("guest did not receive shared content: %+v", page)
	}
	if page.IsSelf || page.OwnerPageId != "" || page.PageConf.CanDelete || !page.PageConf.CanEdit {
		t.Fatalf("link capability or ownership changed: %+v", page)
	}
}

func TestGuestCannotReadOwnedOrAdminPages(t *testing.T) {
	for _, id := range []string{"Oabc", "Aabc"} {
		_, err := (spaceServiceImpl{}).GetPage(guestPageContext(), &space.GetPageReq{PageId: id})
		if !common.IsErrNeedLogin(err) {
			t.Fatalf("%s: expected login, got %v", id, err)
		}
	}
}

func TestGuestCannotReadRevokedEditableShare(t *testing.T) {
	mock := sharedPageDB(t)
	mock.ExpectQuery("SELECT .* FROM `page` ").WithArgs("Eabc", 1).WillReturnRows(sqlmock.NewRows([]string{"pid"}))
	_, err := (spaceServiceImpl{}).GetPage(guestPageContext(), &space.GetPageReq{PageId: "Eabc"})
	if !common.IsErrNotFound(err) {
		t.Fatalf("revoked share must remain inaccessible: %v", err)
	}
}

func TestGuestStillNeedsLoginForEveryPageMutation(t *testing.T) {
	s := spaceServiceImpl{}
	checks := map[string]func() error{
		"update": func() error {
			_, err := s.UpdatePage(guestPageContext(), &space.UpdatePageReq{PageId: "Eabc", Mask: 7})
			return err
		},
		"transfer": func() error {
			_, err := s.TransferCollection(guestPageContext(), &space.TransferCollectionReq{SourcePageId: "Eabc"})
			return err
		},
		"delete": func() error {
			_, err := s.DeletePage(guestPageContext(), &space.DeletePageReq{PageId: "Eabc"})
			return err
		},
		"create link": func() error {
			_, err := s.AddPageLink(guestPageContext(), &space.AddPageLinkReq{PageId: "Eabc", PageType: "readonly"})
			return err
		},
		"revoke link": func() error {
			_, err := s.RemovePageLink(guestPageContext(), &space.RemovePageLinkReq{PageId: "Eabc", PageType: "edit"})
			return err
		},
	}
	for name, check := range checks {
		t.Run(name, func(t *testing.T) {
			if err := check(); !common.IsErrNeedLogin(err) {
				t.Fatalf("expected login: %v", err)
			}
		})
	}
}
