package achievement

import (
	"time"

	"github.com/google/uuid"
)

type LocalizedText struct {
	EN string
	RU string
	KK string
}

type Achievement struct {
	ID          uuid.UUID
	Code        string
	Title       LocalizedText
	Description LocalizedText
	IconKey     string
	MetricKey   string
	Goal        int
	Progress    int
	Unlocked    bool
	UnlockedAt  *time.Time
	SortOrder   int
}
