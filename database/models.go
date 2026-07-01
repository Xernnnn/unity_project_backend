package database

import (
	"time"
)

// User model
type User struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement;not null;<-:create"`
	Name      string    `gorm:"column:name;not null"`
	Email     string    `gorm:"column:email;uniqueIndex;not null"`
	Password  string    `gorm:"column:password;not null"`
	Role      string    `gorm:"column:role;type:varchar(50);not null;default:'customer'"` // customer, admin
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

type LearningNote struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement;not null;<-:create"`
	UserID    int       `gorm:"column:user_id;not null;index"`
	Title     string    `gorm:"column:title;not null"`
	Content   string    `gorm:"column:content;type:text;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

type Todo struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement"`
	Title       string    `gorm:"column:title;not null"`
	Description string    `gorm:"column:description"`
	IsDone      bool      `gorm:"column:is_done;default:false"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
}
