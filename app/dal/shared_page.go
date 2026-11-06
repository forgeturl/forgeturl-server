package dal

import (
	"context"
	"slices"

	"forgeturl-server/api/common"
	"forgeturl-server/dal/model"
	"forgeturl-server/dal/query"
	"gorm.io/gorm/clause"
)

// Lock a stable user row, including when the space has no associations. Every
// operation which rewrites membership/order must acquire this lock first.
func WithUserPageLock(ctx context.Context, uid int64, action func(*query.Query) error) error {
	return Q.Transaction(func(tx *query.Query) error {
		u := tx.User
		if _, err := u.WithContext(ctx).Select(u.ID).Where(u.ID.Eq(uid)).Clauses(clause.Locking{Strength: "UPDATE"}).First(); err != nil {
			return transGormErr(err)
		}
		return action(tx)
	})
}

// AppendSharedPage must run inside WithUserPageLock. Upsert revives legacy
// soft-deleted rows without colliding with UNIQUE(uid, pid).
func (*userPageImpl) AppendSharedPage(ctx context.Context, uid int64, pid string, tx *query.Query) ([]string, bool, error) {
	u := tx.UserPage
	rows, err := u.WithContext(ctx).Select(u.Pid, u.Sort).Where(u.UID.Eq(uid)).Order(u.Sort.Asc(), u.ID.Asc()).Find()
	if err != nil {
		return nil, false, transGormErr(err)
	}
	ids := make([]string, 0, len(rows)+1)
	var next int64
	already := false
	for _, row := range rows {
		ids = append(ids, row.Pid)
		if row.Sort >= next {
			next = row.Sort + 1
		}
		if row.Pid == pid {
			already = true
		}
	}
	if already {
		return ids, true, nil
	}
	err = u.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uid"}, {Name: "pid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"deleted_at": nil, "sort": next}),
	}).Create(&model.UserPage{UID: uid, Pid: pid, Sort: next})
	if err != nil {
		return nil, false, transGormErr(err)
	}
	return append(ids, pid), false, nil
}

// New clients compare their original list atomically. Legacy clients can still
// add/reorder, but omission no longer removes unseen concurrent additions.
func ReconcilePageOrder(current, requested, base []string, checkBase bool) ([]string, error) {
	if checkBase && !slices.Equal(current, base) {
		return nil, common.ErrConflict("space changed; refresh before sorting")
	}
	seen := make(map[string]bool, len(requested))
	result := make([]string, 0, len(current)+len(requested))
	for _, id := range requested {
		if seen[id] {
			return nil, common.ErrBadRequest("duplicate page id")
		}
		seen[id] = true
		result = append(result, id)
	}
	for _, id := range current {
		if !seen[id] {
			result = append(result, id)
		}
	}
	return result, nil
}
