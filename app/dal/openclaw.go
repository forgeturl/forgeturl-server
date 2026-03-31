package dal

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"forgeturl-server/api/common"
	"forgeturl-server/dal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	openClawAPIKeyTable = "openclaw_api_key"
	tmpBookmarkTable    = "tmp_bookmark"

	openClawAPIKeyEnabled = 1
	tmpBookmarkEnabled    = 1
)

type openClawImpl struct{}

var OpenClaw = &openClawImpl{}

type OpenClawAPIKey struct {
	ID        int64      `gorm:"column:id"`
	UID       int64      `gorm:"column:uid"`
	APIKey    string     `gorm:"column:api_key"`
	Status    int32      `gorm:"column:status"`
	ExpiresAt *time.Time `gorm:"column:expires_at"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

func (*OpenClawAPIKey) TableName() string {
	return openClawAPIKeyTable
}

type TmpBookmark struct {
	ID            int64      `gorm:"column:id"`
	UID           int64      `gorm:"column:uid"`
	Title         string     `gorm:"column:title"`
	TitleHash     string     `gorm:"column:title_hash"`
	URL           string     `gorm:"column:url"`
	NormalizedURL string     `gorm:"column:normalized_url"`
	URLHash       string     `gorm:"column:url_hash"`
	Status        int32      `gorm:"column:status"`
	CreatedAt     *time.Time `gorm:"column:created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
}

func (*TmpBookmark) TableName() string {
	return tmpBookmarkTable
}

func (*openClawImpl) Transaction(ctx context.Context, fc func(tx *gorm.DB) error) error {
	return db.WithContext(ctx).Transaction(fc)
}

func (*openClawImpl) GetAPIKeyByUID(ctx context.Context, uid int64) (*OpenClawAPIKey, error) {
	var info OpenClawAPIKey
	err := db.WithContext(ctx).Table(openClawAPIKeyTable).
		Where("uid = ? AND status = ?", uid, openClawAPIKeyEnabled).
		Order("id DESC").
		First(&info).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, transGormErr(err)
	}
	return &info, nil
}

func (*openClawImpl) RegenerateAPIKey(ctx context.Context, uid int64) (*OpenClawAPIKey, error) {
	newKey := genOpenClawAPIKey()
	var oldKeys []OpenClawAPIKey
	info := &OpenClawAPIKey{
		UID:    uid,
		APIKey: newKey,
		Status: openClawAPIKeyEnabled,
	}

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err0 := tx.Table(openClawAPIKeyTable).Where("uid = ?", uid).Find(&oldKeys).Error; err0 != nil {
			return err0
		}
		if err0 := tx.Table(openClawAPIKeyTable).Where("uid = ?", uid).Delete(&OpenClawAPIKey{}).Error; err0 != nil {
			return err0
		}
		if err0 := tx.Table(openClawAPIKeyTable).Create(info).Error; err0 != nil {
			return err0
		}
		return nil
	})
	if err != nil {
		return nil, transGormErr(err)
	}

	for _, k := range oldKeys {
		if k.APIKey != "" {
			_ = C.DelOpenClawAPIKey(ctx, k.APIKey)
		}
	}
	_ = C.SetOpenClawAPIKey(ctx, newKey, uid)
	return info, nil
}

func (*openClawImpl) GetUIDByAPIKey(ctx context.Context, apiKey string) int64 {
	if apiKey == "" {
		return 0
	}
	if uid := C.GetOpenClawAPIKey(ctx, apiKey); uid > 0 {
		return uid
	}

	var info OpenClawAPIKey
	err := db.WithContext(ctx).Table(openClawAPIKeyTable).
		Where("api_key = ? AND status = ?", apiKey, openClawAPIKeyEnabled).
		First(&info).Error
	if err != nil {
		return 0
	}
	if info.ExpiresAt != nil && info.ExpiresAt.Before(time.Now()) {
		return 0
	}
	_ = C.SetOpenClawAPIKey(ctx, apiKey, info.UID)
	return info.UID
}

func (*openClawImpl) UpsertTmpBookmark(ctx context.Context, uid int64, title, rawURL string, tx ...*gorm.DB) (*TmpBookmark, error) {
	title = strings.TrimSpace(title)
	rawURL = strings.TrimSpace(rawURL)
	normURL, err := NormalizeBookmarkURL(rawURL)
	if err != nil {
		return nil, err
	}

	titleHash := hashText(strings.ToLower(title))
	urlHash := hashText(normURL)
	dbRef := db.WithContext(ctx)
	if len(tx) > 0 && tx[0] != nil {
		dbRef = tx[0].WithContext(ctx)
	}

	var existed TmpBookmark
	err = dbRef.Table(tmpBookmarkTable).
		Where("uid = ? AND title_hash = ? AND url_hash = ? AND status = ?", uid, titleHash, urlHash, tmpBookmarkEnabled).
		First(&existed).Error
	if err == nil {
		return &existed, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, transGormErr(err)
	}

	info := &TmpBookmark{
		UID:           uid,
		Title:         title,
		TitleHash:     titleHash,
		URL:           rawURL,
		NormalizedURL: normURL,
		URLHash:       urlHash,
		Status:        tmpBookmarkEnabled,
	}
	err = dbRef.Table(tmpBookmarkTable).Create(info).Error
	if err == nil {
		return info, nil
	}
	if !errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, transGormErr(err)
	}

	err = dbRef.Table(tmpBookmarkTable).
		Where("uid = ? AND title_hash = ? AND url_hash = ? AND status = ?", uid, titleHash, urlHash, tmpBookmarkEnabled).
		First(&existed).Error
	if err != nil {
		return nil, transGormErr(err)
	}
	return &existed, nil
}

func (*openClawImpl) ListTmpBookmarks(ctx context.Context, uid int64, limit int) ([]*TmpBookmark, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	var list []*TmpBookmark
	err := db.WithContext(ctx).Table(tmpBookmarkTable).
		Where("uid = ? AND status = ?", uid, tmpBookmarkEnabled).
		Order("id DESC").
		Limit(limit).
		Find(&list).Error
	if err != nil {
		return nil, transGormErr(err)
	}
	return list, nil
}

func (*openClawImpl) ExistsTmpBookmark(ctx context.Context, uid int64, title, rawURL string) (*TmpBookmark, bool, error) {
	title = strings.TrimSpace(title)
	rawURL = strings.TrimSpace(rawURL)
	normURL, err := NormalizeBookmarkURL(rawURL)
	if err != nil {
		return nil, false, err
	}

	titleHash := hashText(strings.ToLower(title))
	urlHash := hashText(normURL)

	var info TmpBookmark
	err = db.WithContext(ctx).Table(tmpBookmarkTable).
		Where("uid = ? AND title_hash = ? AND url_hash = ? AND status = ?", uid, titleHash, urlHash, tmpBookmarkEnabled).
		First(&info).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, transGormErr(err)
	}
	return &info, true, nil
}

func (*openClawImpl) GetTmpBookmarkByID(ctx context.Context, uid, id int64, tx ...*gorm.DB) (*TmpBookmark, error) {
	dbRef := db.WithContext(ctx)
	if len(tx) > 0 && tx[0] != nil {
		dbRef = tx[0].WithContext(ctx)
	}
	var info TmpBookmark
	err := dbRef.Table(tmpBookmarkTable).
		Where("id = ? AND uid = ? AND status = ?", id, uid, tmpBookmarkEnabled).
		First(&info).Error
	if err != nil {
		return nil, transGormErr(err)
	}
	return &info, nil
}

func (*openClawImpl) DeleteTmpBookmarkByID(ctx context.Context, uid, id int64, tx ...*gorm.DB) error {
	dbRef := db.WithContext(ctx)
	if len(tx) > 0 && tx[0] != nil {
		dbRef = tx[0].WithContext(ctx)
	}
	err := dbRef.Table(tmpBookmarkTable).
		Where("id = ? AND uid = ? AND status = ?", id, uid, tmpBookmarkEnabled).
		Delete(&TmpBookmark{}).Error
	if err != nil {
		return transGormErr(err)
	}
	return nil
}

func (*openClawImpl) LockOwnerPage(ctx context.Context, ownerPid string, tx *gorm.DB) (*model.Page, error) {
	var page model.Page
	err := tx.WithContext(ctx).Table(model.TableNamePage).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("pid = ?", ownerPid).
		First(&page).Error
	if err != nil {
		return nil, transGormErr(err)
	}
	return &page, nil
}

func (*openClawImpl) UpdateOwnerPageContentByVersion(ctx context.Context, ownerPid string, version int64, newContent string, tx *gorm.DB) error {
	result := tx.WithContext(ctx).Table(model.TableNamePage).
		Where("pid = ? AND version = ?", ownerPid, version).
		Updates(map[string]any{
			"content": newContent,
			"version": version + 1,
		})
	if result.Error != nil {
		return transGormErr(result.Error)
	}
	if result.RowsAffected == 0 {
		return common.ErrUpdateMissNeedRefreshPage()
	}
	return nil
}

func NormalizeBookmarkURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("url is empty")
	}
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("invalid url")
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("invalid url")
	}

	scheme := strings.ToLower(parsed.Scheme)
	host := strings.ToLower(parsed.Hostname())
	port := parsed.Port()
	if (scheme == "http" && port == "80") || (scheme == "https" && port == "443") {
		port = ""
	}
	if port != "" {
		host = host + ":" + port
	}

	path := parsed.EscapedPath()
	if path == "" {
		path = "/"
	}
	if path != "/" {
		path = strings.TrimRight(path, "/")
		if path == "" {
			path = "/"
		}
	}

	queryValues := parsed.Query()
	rawQuery := queryValues.Encode()

	normalized := &url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     path,
		RawQuery: rawQuery,
	}
	return normalized.String(), nil
}

func hashText(text string) string {
	sum := sha256.Sum256([]byte(text))
	return hex.EncodeToString(sum[:])
}

func genOpenClawAPIKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 48)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("crypto/rand failed: %v", err))
	}
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return "sk-" + string(b)
}
