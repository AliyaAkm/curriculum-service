package codeattempt

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"curriculum-service/internal/domain"
	codeattemptdomain "curriculum-service/internal/domain/codeattempt"
	practicedomain "curriculum-service/internal/domain/practice"

	"github.com/google/uuid"
)

func (u *UseCase) Run(ctx context.Context, req codeattemptdomain.RunRequest) (*codeattemptdomain.RunResult, error) {
	if strings.TrimSpace(req.PracticeID) == "" || strings.TrimSpace(req.Code) == "" {
		return nil, domain.ErrValidation
	}

	practiceID, err := uuid.Parse(req.PracticeID)
	if err != nil {
		return nil, domain.ErrValidation
	}
	if u.practices == nil {
		return nil, domain.ErrInternal
	}
	practice, err := u.practices.GetByID(ctx, practiceID)
	if err != nil {
		return nil, err
	}

	req.CourseID = &practice.CourseID
	req.LessonID = &practice.LessonID
	if strings.TrimSpace(req.Language) == "" {
		req.Language = practice.Language
	}

	runType := req.RunType
	if runType == "" {
		runType = codeattemptdomain.RunTypeRun
	}
	if runType != codeattemptdomain.RunTypeRun && runType != codeattemptdomain.RunTypeSubmit {
		return nil, domain.ErrValidation
	}
	if practice.CheckType == practicedomain.CheckTypeManual && runType == codeattemptdomain.RunTypeSubmit {
		return nil, domain.ErrPracticeAutoSubmitNotAllowed
	}
	if !strings.EqualFold(strings.TrimSpace(req.Language), strings.TrimSpace(practice.Language)) {
		return nil, domain.ErrValidation
	}
	canStart, err := u.repo.CanStartPractice(ctx, req.UserID, practice.ID)
	if err != nil {
		return nil, err
	}
	if !canStart {
		return nil, domain.ErrPracticePrerequisitesNotMet
	}

	started := time.Now()
	runnerResult, err := u.runner.Run(ctx, practice.Language, req.Code)
	durationMS := int(time.Since(started).Milliseconds())

	errorType := ""
	errorMessage := runnerResult.Error
	passed := err == nil && runnerResult.Passed
	if err != nil {
		errorType = "runner_error"
		errorMessage = err.Error()
	} else if runnerResult.Error != "" {
		errorType = "execution_error"
	} else if runType == codeattemptdomain.RunTypeSubmit && !sameOutput(runnerResult.Output, practice.ExpectedOutput) {
		passed = false
		errorType = "wrong_answer"
		errorMessage = "output does not match expected output"
	}

	hash := sha256.Sum256([]byte(req.Code))
	attempt, saveErr := u.repo.CreateAttempt(ctx, codeattemptdomain.Attempt{
		UserID:       req.UserID,
		CourseID:     req.CourseID,
		LessonID:     req.LessonID,
		PracticeID:   practice.ID.String(),
		RunType:      runType,
		Language:     practice.Language,
		Passed:       passed,
		ErrorType:    errorType,
		ErrorMessage: errorMessage,
		Output:       runnerResult.Output,
		DurationMS:   durationMS,
		CodeHash:     hex.EncodeToString(hash[:]),
		XPReward:     practice.XPReward,
	})
	if saveErr != nil {
		return nil, saveErr
	}
	if err != nil {
		return nil, err
	}

	return &codeattemptdomain.RunResult{
		AttemptID:  attempt.ID,
		Output:     runnerResult.Output,
		Error:      errorMessage,
		Passed:     passed,
		ErrorType:  errorType,
		DurationMS: durationMS,
		XPAwarded:  attempt.XPAwarded,
	}, nil
}

func sameOutput(actual, expected string) bool {
	normalize := func(value string) string {
		value = strings.ReplaceAll(value, "\r\n", "\n")
		return strings.TrimSpace(value)
	}
	return normalize(actual) == normalize(expected)
}
