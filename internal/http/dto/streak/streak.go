package streak

/*type DailyStreakRequest struct {
	UserID uuid.UUID `json:"user_id"`
}*/

type DailyStreakResponse struct {
	Streak int64 `json:"streak"`
}
