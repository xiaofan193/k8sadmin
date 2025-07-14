package model

import (
	"time"
)

type User struct {
	ID        uint64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	Name      string     `gorm:"column:name;type:varchar(100);not null" json:"name"`
	CreatedAt *time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;not null" json:"createdAt"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;not null" json:"updatedAt"`
}

// TableName table name
func (m *User) TableName() string {
	return "user"
}

// UserColumnNames Whitelist for custom query fields to prevent sql injection attacks
var UserColumnNames = map[string]bool{
	"id":         true,
	"name":       true,
	"created_at": true,
	"updated_at": true,
}
