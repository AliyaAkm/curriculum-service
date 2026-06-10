package teacherstats

import (
	"time"

	"github.com/google/uuid"
)

type Filter struct {
	TeacherID  uuid.UUID
	IsAdmin    bool
	CourseID   *uuid.UUID
	PeriodDays int
}

type Statistics struct {
	Summary     Summary           `json:"summary"`
	Practice    PracticeStats     `json:"practice"`
	Quiz        QuizStats         `json:"quiz"`
	Activity    []ActivityDay     `json:"activity"`
	Funnel      []FunnelStep      `json:"funnel"`
	QuizHeatmap []QuizHeatmapItem `json:"quiz_heatmap"`
	Courses     []CourseStats     `json:"courses"`
	NewStudents []NewStudent      `json:"new_students"`
	ReviewQueue []ReviewItem      `json:"review_queue"`
}

type Summary struct {
	Courses            int        `json:"courses"`
	Students           int        `json:"students"`
	NewStudents        int        `json:"new_students"`
	ActiveStudents     int        `json:"active_students"`
	CompletedStudents  int        `json:"completed_students"`
	AvgProgressPercent int        `json:"avg_progress_percent"`
	TotalXPAwarded     int        `json:"total_xp_awarded"`
	LastActivityAt     *time.Time `json:"last_activity_at,omitempty"`
}

type PracticeStats struct {
	AutoRuns                 int `json:"auto_runs"`
	AutoSubmissions          int `json:"auto_submissions"`
	AutoSuccessfulSubmits    int `json:"auto_successful_submits"`
	ManualSubmissions        int `json:"manual_submissions"`
	ManualApproved           int `json:"manual_approved"`
	ManualPendingReview      int `json:"manual_pending_review"`
	ManualChangesRequested   int `json:"manual_changes_requested"`
	CompletedPractices       int `json:"completed_practices"`
	PassRatePercent          int `json:"pass_rate_percent"`
	AvgAttemptsPerCompletion int `json:"avg_attempts_per_completion"`
}

type QuizStats struct {
	AttemptedQuizzes int `json:"attempted_quizzes"`
	PassedQuizzes    int `json:"passed_quizzes"`
	AccuracyPercent  int `json:"accuracy_percent"`
}

type ActivityDay struct {
	Date                string `json:"date"`
	ActiveStudents      int    `json:"active_students"`
	PracticeSubmissions int    `json:"practice_submissions"`
	PracticeApproved    int    `json:"practice_approved"`
	QuizAttempts        int    `json:"quiz_attempts"`
	XPAwarded           int    `json:"xp_awarded"`
}

type FunnelStep struct {
	Key               string `json:"key"`
	Label             string `json:"label"`
	Value             int    `json:"value"`
	ConversionPercent int    `json:"conversion_percent"`
	DropOffPercent    int    `json:"drop_off_percent"`
}

type QuizHeatmapItem struct {
	QuizID          uuid.UUID `json:"quiz_id"`
	CourseID        uuid.UUID `json:"course_id"`
	CourseTitle     string    `json:"course_title"`
	LessonID        uuid.UUID `json:"lesson_id"`
	LessonTitle     string    `json:"lesson_title"`
	Position        int       `json:"position"`
	Question        string    `json:"question"`
	Attempts        int       `json:"attempts"`
	CorrectAttempts int       `json:"correct_attempts"`
	AccuracyPercent int       `json:"accuracy_percent"`
}

type CourseStats struct {
	CourseID              uuid.UUID  `json:"course_id"`
	Title                 string     `json:"title"`
	Students              int        `json:"students"`
	ActiveStudents        int        `json:"active_students"`
	CompletedStudents     int        `json:"completed_students"`
	TotalLessons          int        `json:"total_lessons"`
	AvgProgressPercent    int        `json:"avg_progress_percent"`
	PracticePendingReview int        `json:"practice_pending_review"`
	PracticeCompleted     int        `json:"practice_completed"`
	QuizAccuracyPercent   int        `json:"quiz_accuracy_percent"`
	XPAwarded             int        `json:"xp_awarded"`
	LastActivityAt        *time.Time `json:"last_activity_at,omitempty"`
}

type NewStudent struct {
	StudentID    uuid.UUID `json:"student_id"`
	StudentEmail string    `json:"student_email"`
	StudentName  string    `json:"student_name"`
	PhotoURL     string    `json:"photo_url"`
	CourseID     uuid.UUID `json:"course_id"`
	CourseTitle  string    `json:"course_title"`
	SubscribedAt time.Time `json:"subscribed_at"`
}

type ReviewItem struct {
	SubmissionID  uuid.UUID `json:"submission_id"`
	CourseID      uuid.UUID `json:"course_id"`
	CourseTitle   string    `json:"course_title"`
	PracticeID    uuid.UUID `json:"practice_id"`
	PracticeTitle string    `json:"practice_title"`
	StudentID     uuid.UUID `json:"student_id"`
	StudentEmail  string    `json:"student_email"`
	Status        string    `json:"status"`
	AttemptNumber int       `json:"attempt_number"`
	CreatedAt     time.Time `json:"created_at"`
}
