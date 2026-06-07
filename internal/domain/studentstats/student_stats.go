package studentstats

import (
	"time"

	"github.com/google/uuid"
)

type Statistics struct {
	Summary         Summary          `json:"summary"`
	Quiz            QuizStats        `json:"quiz"`
	Practice        PracticeStats    `json:"practice"`
	AI              AIStats          `json:"ai"`
	Activity        []ActivityDay    `json:"activity"`
	Topics          []TopicProgress  `json:"topics"`
	CourseProgress  []CourseDetail   `json:"courseProgress"`
	Recommendations []Recommendation `json:"recommendations"`
}

type Summary struct {
	UserID            uuid.UUID  `json:"user_id"`
	TotalXP           int        `json:"total_xp"`
	Level             int        `json:"level"`
	CurrentStreak     int        `json:"current_streak"`
	MaxStreak         int        `json:"max_streak"`
	StartedCourses    int        `json:"started_courses"`
	ActiveCourses     int        `json:"active_courses"`
	CompletedCourses  int        `json:"completed_courses"`
	TotalLessons      int        `json:"total_lessons"`
	CompletedLessons  int        `json:"completed_lessons"`
	ProgressPercent   int        `json:"progress_percent"`
	Certificates      int        `json:"certificates"`
	Achievements      int        `json:"achievements"`
	TotalAchievements int        `json:"total_achievements"`
	LastActivityAt    *time.Time `json:"last_activity_at,omitempty"`
}

type QuizStats struct {
	AttemptedQuizzes int `json:"attempted_quizzes"`
	PassedQuizzes    int `json:"passed_quizzes"`
	AccuracyPercent  int `json:"accuracy_percent"`
}

type PracticeStats struct {
	Runs               int        `json:"runs"`
	Submissions        int        `json:"submissions"`
	SuccessfulSubmits  int        `json:"successful_submits"`
	AttemptedPractices int        `json:"attempted_practices"`
	CompletedPractices int        `json:"completed_practices"`
	PassRatePercent    int        `json:"pass_rate_percent"`
	XPEarned           int        `json:"xp_earned"`
	LastAttemptAt      *time.Time `json:"last_attempt_at,omitempty"`
}

type AIStats struct {
	Available         bool       `json:"available"`
	TotalRequests     int        `json:"total_requests"`
	CompletedRequests int        `json:"completed_requests"`
	FailedRequests    int        `json:"failed_requests"`
	Chats             int        `json:"chats"`
	UserMessages      int        `json:"user_messages"`
	AssistantMessages int        `json:"assistant_messages"`
	InputTokens       int        `json:"input_tokens"`
	OutputTokens      int        `json:"output_tokens"`
	AvgLatencyMS      int        `json:"avg_latency_ms"`
	LastActivityAt    *time.Time `json:"last_activity_at,omitempty"`
}

type ActivityDay struct {
	Date              string `json:"date"`
	LessonsCompleted  int    `json:"lessons_completed"`
	QuizzesPassed     int    `json:"quizzes_passed"`
	PracticeAttempts  int    `json:"practice_attempts"`
	PracticeCompleted int    `json:"practice_completed"`
	AIRequests        int    `json:"ai_requests"`
	XP                int    `json:"xp"`
}

type TopicProgress struct {
	TopicID          uuid.UUID `json:"topic_id"`
	Code             string    `json:"code"`
	Name             string    `json:"name"`
	TotalLessons     int       `json:"total_lessons"`
	CompletedLessons int       `json:"completed_lessons"`
	ProgressPercent  int       `json:"progress_percent"`
	XP               int       `json:"xp"`
}

type CourseDetail struct {
	CourseID          uuid.UUID      `json:"course_id"`
	Title             string         `json:"title"`
	StartedAt         *time.Time     `json:"started_at,omitempty"`
	LastActivityAt    *time.Time     `json:"last_activity_at,omitempty"`
	CompletedAt       *time.Time     `json:"completed_at,omitempty"`
	CurrentLessonID   *uuid.UUID     `json:"current_lesson_id,omitempty"`
	TotalLessons      int            `json:"total_lessons"`
	CompletedLessons  int            `json:"completed_lessons"`
	ProgressPercent   int            `json:"progress_percent"`
	PracticeRuns      int            `json:"practice_runs"`
	PracticeCompleted int            `json:"practice_completed"`
	Modules           []ModuleDetail `json:"modules" gorm:"-"`
}

type ModuleDetail struct {
	ModuleID         uuid.UUID `json:"module_id"`
	Title            string    `json:"title"`
	Position         int       `json:"position"`
	TotalLessons     int       `json:"total_lessons"`
	CompletedLessons int       `json:"completed_lessons"`
	ProgressPercent  int       `json:"progress_percent"`
	IsOpen           bool      `json:"is_open"`
}

type Recommendation struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Message string `json:"message"`
}
