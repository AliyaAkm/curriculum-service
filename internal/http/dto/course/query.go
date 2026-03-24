package course

type GetCoursesQuery struct {
	Search           string   `form:"search" json:"search"`
	Topic            string   `form:"topic" json:"topic"`
	Level            string   `form:"level" json:"level"`
	MinRating        *float64 `form:"min_rating" json:"min_rating"`
	DurationCategory string   `form:"duration_category" json:"duration_category"`
	HasCertificate   *bool    `form:"has_certificate" json:"has_certificate"`
	Page             int      `form:"page" json:"page"`
	Limit            int      `form:"limit" json:"limit"`
}
