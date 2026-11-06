package api

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"forgeturl-server/api/common"
	"forgeturl-server/api/space"
	"forgeturl-server/dal"
	"forgeturl-server/dal/query"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func sharedPageDB(t *testing.T) sqlmock.Sqlmock {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{DisableAutomaticPing: true, SkipDefaultTransaction: true})
	if err != nil {
		t.Fatal(err)
	}
	original := dal.Q
	dal.Q = query.Use(gdb)
	t.Cleanup(func() {
		dal.Q = original
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
		_ = db.Close()
	})
	return mock
}

func lockSpace(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .* FROM `user` .*FOR UPDATE").WithArgs(int64(42), 1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))
}

func sharedRow(owner int64) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"uid", "pid", "edit_pid", "readonly_pid", "title", "created_at", "updated_at"}).AddRow(owner, "Oabc", "Eabc", "Rabc", "Original", time.Unix(100, 0), time.Unix(100, 0))
}

func TestSharedPageRejectsUnauthenticatedAndMalformedRequests(t *testing.T) {
	ctx := context.Background()
	if _, err := addSharedPage(ctx, 0, "Eabc"); !common.IsErrNeedLogin(err) {
		t.Fatalf("expected login: %v", err)
	}
	if _, err := removeSharedPage(ctx, 0, "Eabc"); !common.IsErrNeedLogin(err) {
		t.Fatalf("expected login: %v", err)
	}
	for _, id := range []string{"", "Oabc", "T123", "X123", "E", "E../bad"} {
		if _, err := addSharedPage(ctx, 42, id); !common.IsErrBadRequest(err) {
			t.Errorf("add %q: %v", id, err)
		}
		if _, err := removeSharedPage(ctx, 42, id); !common.IsErrBadRequest(err) {
			t.Errorf("remove %q: %v", id, err)
		}
	}
}

func TestSharedPageJoinRetainsTokenAndOriginalPermissions(t *testing.T) {
	for _, id := range []string{"Eabc", "Rabc"} {
		t.Run(id, func(t *testing.T) {
			mock := sharedPageDB(t)
			lockSpace(mock)
			mock.ExpectQuery("SELECT .* FROM `page` .*FOR UPDATE").WithArgs(id, 1).WillReturnRows(sharedRow(7))
			mock.ExpectQuery("SELECT .* FROM `user_page` .*ORDER BY").WithArgs(int64(42)).WillReturnRows(sqlmock.NewRows([]string{"pid", "sort"}).AddRow("Oold", 0))
			mock.ExpectExec("INSERT INTO `user_page` .*ON DUPLICATE KEY UPDATE").WithArgs(int64(42), id, int64(1), nil, sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(10, 1))
			mock.ExpectCommit()
			resp, err := addSharedPage(context.Background(), 42, id)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(resp.PageIds, []string{"Oold", id}) || resp.AlreadyAdded || resp.PageBrief.PageId != id || resp.PageBrief.IsSelf {
				t.Fatalf("wrong association: %+v", resp)
			}
			if resp.PageBrief.PageConf.CanEdit != (id == "Eabc") || resp.PageBrief.PageConf.CanDelete {
				t.Fatalf("wrong permissions: %+v", resp.PageBrief.PageConf)
			}
		})
	}
}

func TestSharedPageOwnPageReturnsOwnerEntryWithoutAdding(t *testing.T) {
	mock := sharedPageDB(t)
	lockSpace(mock)
	mock.ExpectQuery("SELECT .* FROM `page` .*FOR UPDATE").WithArgs("Rabc", 1).WillReturnRows(sharedRow(42))
	mock.ExpectQuery("SELECT .* FROM `user_page` .*ORDER BY").WithArgs(int64(42)).WillReturnRows(sqlmock.NewRows([]string{"pid"}).AddRow("Oabc"))
	mock.ExpectCommit()
	resp, err := addSharedPage(context.Background(), 42, "Rabc")
	if err != nil || !resp.AlreadyAdded || resp.PageBrief.PageId != "Oabc" || len(resp.PageIds) != 1 {
		t.Fatalf("got %+v %v", resp, err)
	}
}

func TestRevokedShareCannotJoinButCanBeRemovedFromEmptySpace(t *testing.T) {
	mock := sharedPageDB(t)
	lockSpace(mock)
	mock.ExpectQuery("SELECT .* FROM `page` .*FOR UPDATE").WithArgs("Eabc", 1).WillReturnRows(sqlmock.NewRows([]string{"pid"}))
	mock.ExpectRollback()
	if _, err := addSharedPage(context.Background(), 42, "Eabc"); !common.IsErrNotFound(err) {
		t.Fatalf("got %v", err)
	}
	for range 2 {
		lockSpace(mock)
		mock.ExpectQuery("SELECT .* FROM `page` .*FOR UPDATE").WithArgs("Eabc", 1).WillReturnRows(sqlmock.NewRows([]string{"pid"}))
		mock.ExpectExec("DELETE FROM `user_page` WHERE .*uid.*pid").WithArgs(int64(42), "Eabc").WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectQuery("SELECT .* FROM `user_page` .*ORDER BY").WithArgs(int64(42)).WillReturnRows(sqlmock.NewRows([]string{"pid"}))
		mock.ExpectCommit()
		resp, err := removeSharedPage(context.Background(), 42, "Eabc")
		if err != nil || len(resp.PageIds) != 0 {
			t.Fatalf("got %+v %v", resp, err)
		}
	}
}

func TestMySpaceMissingShareIsPlaceholderButDatabaseFailurePropagates(t *testing.T) {
	for _, brokenDB := range []bool{false, true} {
		t.Run(map[bool]string{false: "missing", true: "database failure"}[brokenDB], func(t *testing.T) {
			mock := sharedPageDB(t)
			mock.ExpectQuery("SELECT .* FROM `user` ").WithArgs(int64(42), 1).WillReturnRows(sqlmock.NewRows([]string{"id", "display_name"}).AddRow(42, "B"))
			mock.ExpectQuery("SELECT .* FROM `user_page` .*ORDER BY").WithArgs(int64(42)).WillReturnRows(sqlmock.NewRows([]string{"pid"}).AddRow("Eabc").AddRow("Rabc"))
			first := mock.ExpectQuery("SELECT .* FROM `page` ").WithArgs("Eabc", 1)
			if brokenDB {
				first.WillReturnError(errors.New("database offline"))
			} else {
				first.WillReturnRows(sqlmock.NewRows([]string{"pid"}))
				mock.ExpectQuery("SELECT .* FROM `page` ").WithArgs("Rabc", 1).WillReturnRows(sharedRow(7))
			}
			resp, err := getMySpace(context.Background(), 42)
			if brokenDB {
				if err == nil {
					t.Fatal("database errors must not be swallowed")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(resp.PageBriefs) != 2 || !resp.PageBriefs[0].IsNotFound || resp.PageBriefs[0].PageConf.CanEdit || resp.PageBriefs[1].Title != "Original" {
				t.Fatalf("got %+v", resp)
			}
		})
	}
}

func TestStaleSortingCannotOverwriteJoinedPage(t *testing.T) {
	mock := sharedPageDB(t)
	lockSpace(mock)
	mock.ExpectQuery("SELECT .* FROM `user_page` .*ORDER BY").WithArgs(int64(42)).WillReturnRows(sqlmock.NewRows([]string{"pid"}).AddRow("Oabc").AddRow("Eabc"))
	mock.ExpectRollback()
	_, err := savePageIds(context.Background(), 42, &space.SavePageIdsReq{PageIds: []string{"Oabc"}, BasePageIds: []string{"Oabc"}, CheckBase: true})
	if !common.IsErrConflict(err) {
		t.Fatalf("must reject stale sorting: %v", err)
	}
}

func TestJoinExistingTokenDoesNotWriteOrChangeSort(t *testing.T) {
	mock := sharedPageDB(t)
	lockSpace(mock)
	mock.ExpectQuery("SELECT .* FROM `page` .*FOR UPDATE").WithArgs("Eabc", 1).WillReturnRows(sharedRow(7))
	mock.ExpectQuery("SELECT .* FROM `user_page` .*ORDER BY").WithArgs(int64(42)).WillReturnRows(sqlmock.NewRows([]string{"pid", "sort"}).AddRow("Eabc", 3).AddRow("Oold", 9))
	mock.ExpectCommit()
	resp, err := addSharedPage(context.Background(), 42, "Eabc")
	if err != nil || !resp.AlreadyAdded || !reflect.DeepEqual(resp.PageIds, []string{"Eabc", "Oold"}) {
		t.Fatalf("got %+v %v", resp, err)
	}
}
