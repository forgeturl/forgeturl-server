package dal

import (
	"context"
	"reflect"
	"testing"

	"forgeturl-server/dal/model"
	"forgeturl-server/dal/query"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSaveUserPageIdsRebuildsExistingAssociations(t *testing.T) {
	testDB, err := gorm.Open(sqlite.Open("file:user_page_test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, err := testDB.DB()
	if err != nil {
		t.Fatalf("get sqlite connection: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	if err = testDB.AutoMigrate(&model.UserPage{}); err != nil {
		t.Fatalf("migrate user_page: %v", err)
	}

	originalQ := Q
	Q = query.Use(testDB)
	t.Cleanup(func() { Q = originalQ })

	const uid int64 = 42
	if err = testDB.Create(&model.UserPage{UID: uid, Pid: "O_existing", Sort: 0}).Error; err != nil {
		t.Fatalf("create existing association: %v", err)
	}

	want := []string{"O_new", "O_existing"}
	if err = UserPage.SaveUserPageIds(context.Background(), uid, want); err != nil {
		t.Fatalf("save multiple page ids: %v", err)
	}

	got, err := UserPage.GetUserPageIds(context.Background(), uid)
	if err != nil {
		t.Fatalf("get page ids: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("page ids mismatch: got %v, want %v", got, want)
	}

	var total int64
	if err = testDB.Unscoped().Model(&model.UserPage{}).Where("uid = ?", uid).Count(&total).Error; err != nil {
		t.Fatalf("count associations: %v", err)
	}
	if total != int64(len(want)) {
		t.Fatalf("unexpected soft-deleted rows: got %d rows, want %d", total, len(want))
	}
}
