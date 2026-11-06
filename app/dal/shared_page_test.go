package dal

import (
	"context"
	"reflect"
	"testing"

	"forgeturl-server/dal/query"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func mockUserPages(t *testing.T) sqlmock.Sqlmock {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{DisableAutomaticPing: true, SkipDefaultTransaction: true})
	if err != nil {
		t.Fatal(err)
	}
	original := Q
	Q = query.Use(gdb)
	t.Cleanup(func() {
		Q = original
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
		_ = db.Close()
	})
	return mock
}

func TestAppendSharedPagePreservesOrderAndRestoresSoftDeletedID(t *testing.T) {
	mock := mockUserPages(t)
	mock.ExpectQuery("SELECT .* FROM `user_page` WHERE .* ORDER BY").WithArgs(int64(42)).WillReturnRows(sqlmock.NewRows([]string{"pid", "sort"}).AddRow("Oabc", 3).AddRow("Rabc", 8))
	// Upsert must clear deleted_at, otherwise a previously removed unique key
	// would prevent rejoining. The sort follows the highest existing position.
	mock.ExpectExec("INSERT INTO `user_page` .*ON DUPLICATE KEY UPDATE `deleted_at`=.*,`sort`=.*").WithArgs(int64(42), "Eabc", int64(9), nil, sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(10, 1))
	ids, exists, err := UserPage.AppendSharedPage(context.Background(), 42, "Eabc", Q)
	if err != nil || exists || !reflect.DeepEqual(ids, []string{"Oabc", "Rabc", "Eabc"}) {
		t.Fatalf("got %v %v %v", ids, exists, err)
	}
}

func TestAppendSharedPageAlreadyAddedKeepsSort(t *testing.T) {
	mock := mockUserPages(t)
	mock.ExpectQuery("SELECT .* FROM `user_page` WHERE .* ORDER BY").WithArgs(int64(42)).WillReturnRows(sqlmock.NewRows([]string{"pid", "sort"}).AddRow("Eabc", 3).AddRow("Oabc", 8))
	ids, exists, err := UserPage.AppendSharedPage(context.Background(), 42, "Eabc", Q)
	if err != nil || !exists || !reflect.DeepEqual(ids, []string{"Eabc", "Oabc"}) {
		t.Fatalf("got %v %v %v", ids, exists, err)
	}
}

func TestRemoveSharedPageIsPhysicalAndScopedToUser(t *testing.T) {
	mock := mockUserPages(t)
	for range 2 {
		mock.ExpectExec("DELETE FROM `user_page` WHERE .*uid.*pid").WithArgs(int64(42), "Eabc").WillReturnResult(sqlmock.NewResult(0, 0))
		if _, err := UserPage.DeleteUserPageId(context.Background(), 42, "Eabc", Q); err != nil {
			t.Fatal(err)
		}
	}
}

func TestPageOrderRejectsStaleBaseAndPreservesUnseenMemberships(t *testing.T) {
	current := []string{"Oabc", "Rabc", "Enew"}
	if _, err := ReconcilePageOrder(current, []string{"Rabc", "Oabc"}, []string{"Oabc", "Rabc"}, true); err == nil {
		t.Fatal("stale sort must conflict")
	}
	ids, err := ReconcilePageOrder(current, []string{"Rabc", "Oabc"}, nil, false)
	if err != nil || !reflect.DeepEqual(ids, []string{"Rabc", "Oabc", "Enew"}) {
		t.Fatalf("lost concurrent join: %v %v", ids, err)
	}
	if _, err = ReconcilePageOrder(current, []string{"Rabc", "Rabc"}, nil, false); err == nil {
		t.Fatal("duplicate IDs must fail normally")
	}
}

func TestEmptySpaceSaveDoesNotInsertEmptyBatch(t *testing.T) {
	mock := mockUserPages(t)
	mock.ExpectExec("DELETE FROM `user_page` WHERE").WithArgs(int64(42)).WillReturnResult(sqlmock.NewResult(0, 1))
	if err := UserPage.SaveUserPageIds(context.Background(), 42, nil); err != nil {
		t.Fatal(err)
	}
}
