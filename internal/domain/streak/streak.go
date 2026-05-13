package streak

import (
	"github.com/google/uuid"
	"time"
)

type DailyStreak struct {
	ID        uuid.UUID `gorm:"column:id"`
	UserID    uuid.UUID `gorm:"column:user_id"`
	Streak    int64     `gorm:"column:streak"`
	LastLogin time.Time `gorm:"column:last_login"`
}

func (DailyStreak) TableName() string {
	return "daily_streak"
}
