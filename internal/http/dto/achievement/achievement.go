package achievement

import "time"

type LocalizedText struct {
	EN string `json:"en"`
	RU string `json:"ru"`
	KK string `json:"kk"`
}

type Achievement struct {
	ID          string        `json:"id"`
	Title       LocalizedText `json:"title"`
	Description LocalizedText `json:"description"`
	IconKey     string        `json:"icon_key"`
	Goal        int           `json:"goal"`
	Progress    int           `json:"progress"`
	Unlocked    bool          `json:"unlocked"`
	UnlockedAt  *time.Time    `json:"unlocked_at,omitempty"`
}
