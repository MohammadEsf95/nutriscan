package users

import "time"

type User struct {
	ID        int64 `gorm:"primaryKey"`
	ChatID    int64
	CreatedAt time.Time
	Username  string
	FirstName string
	LastName  string
}
