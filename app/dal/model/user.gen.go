// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameUser = "user"

// User mapped from table <user>
type User struct {
	ID            int64          `gorm:"column:id;type:int(11);primaryKey" json:"id"`
	DisplayName   string         `gorm:"column:display_name;type:varchar(128);not null" json:"display_name"`
	Username      string         `gorm:"column:username;type:varchar(64);not null" json:"username"`
	Email         string         `gorm:"column:email;type:varchar(100);not null" json:"email"`
	Avatar        string         `gorm:"column:avatar;type:varchar(1024);not null" json:"avatar"`
	Status        int64          `gorm:"column:status;type:int(11);not null" json:"status"`
	LastLoginDate time.Time      `gorm:"column:last_login_date;type:timestamp" json:"last_login_date"`
	PageIds       string         `gorm:"column:page_ids;type:varchar(2048);not null" json:"page_ids"`
	Provider      string         `gorm:"column:provider;type:varchar(32);not null" json:"provider"`
	ExternalID    string         `gorm:"column:external_id;type:varchar(128);not null" json:"external_id"`
	IPInfo        string         `gorm:"column:ip_info;type:varchar(255);not null" json:"ip_info"`
	IsAdmin       int32          `gorm:"column:is_admin;type:tinyint(1);not null" json:"is_admin"`
	SuspendedAt   time.Time      `gorm:"column:suspended_at;type:timestamp" json:"suspended_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp" json:"deleted_at"`
	CreatedAt     *time.Time     `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
