package lesson

import (
	"curriculum-service/internal/domain"
	"curriculum-service/internal/domain/keypoint"
	"curriculum-service/internal/domain/lesson"
	"curriculum-service/internal/domain/locale"
	"curriculum-service/internal/domain/outcome"
	"curriculum-service/internal/domain/summary"
	"curriculum-service/internal/domain/theorycontent"
	"curriculum-service/internal/domain/title"
	lesson2 "curriculum-service/internal/http/dto/lesson"
	"curriculum-service/internal/http/respond"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func (h *Handler) GetAllLessons(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		respond.JSON(c, http.StatusBadRequest, "empty module id")
		return
	}
	uuidID, err := uuid.Parse(id)
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid module id")
		return
	}
	resp, err := h.client.GetAllLessons(c.Request.Context(), uuidID)
	if err != nil {
		writeCatalogError(c, err)
		return
	}

	c.JSON(http.StatusOK, convertLessons(resp))

}

func (h *Handler) CreateLesson(c *gin.Context) {

	request := lesson2.LessonRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid body")
		return
	}
	locale, err := h.localClient.GetAllLocales(c)
	if err != nil {
		respond.JSON(c, http.StatusInternalServerError, "get locale")
		return
	}

	result, err := h.client.CreateLesson(c.Request.Context(), convertLessonRequest(request, uuid.New(), locale))
	if err != nil {
		writeCatalogError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertLesson(result))
}

func (h *Handler) GetLessonByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid lesson id")
		return
	}

	result, err := h.client.GetLessonByID(c.Request.Context(), id)
	if err != nil {
		writeCatalogError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertLesson(result))
}

func (h *Handler) UpdateLesson(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid lesson id")
		return
	}

	request := lesson2.LessonRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid body")
		return
	}

	locale, err := h.localClient.GetAllLocales(c)
	if err != nil {
		respond.JSON(c, http.StatusInternalServerError, "get locale")
		return
	}

	result, err := h.client.UpdateLesson(c.Request.Context(), id, convertLessonRequest(request, id, locale))
	if err != nil {
		writeCatalogError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertLesson(result))
}

func writeCatalogError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		respond.Error(c, http.StatusBadRequest, "validation", "invalid request query")
	default:
		_ = c.Error(err)
		respond.Error(c, http.StatusInternalServerError, "internal", domain.ErrInternal.Error())
	}
}

func convertLessons(resp []lesson.LessonModel) []lesson2.Lesson {
	lessons := make([]lesson2.Lesson, len(resp))

	for i := range resp {
		lessons[i] = convertLesson(&resp[i])
	}

	return lessons
}

func convertLesson(resp *lesson.LessonModel) lesson2.Lesson {
	if resp == nil {
		return lesson2.Lesson{}
	}

	item := lesson2.Lesson{
		ID:              resp.ID,
		ModuleID:        resp.ModuleID,
		DurationMinutes: resp.DurationMinutes,
		XPReward:        resp.XPReward,
		CodeSnippet:     resp.CodeSnippet,
		ExampleOutput:   resp.ExampleOutput,
		CreatedAt:       resp.CreatedAt,
		UpdatedAt:       resp.UpdatedAt,
	}

	for _, v := range resp.Titles {
		setLocaleValue(&item.Title, v.Locale.Code, v.Name)
	}

	for _, v := range resp.Summaries {
		setLocaleValue(&item.Summary, v.Locale.Code, v.Name)
	}

	for _, v := range resp.Outcomes {
		setLocaleValue(&item.OutCome, v.Locale.Code, v.Name)
	}

	for _, v := range resp.TheoryContents {
		setLocaleValue(&item.TheoryContent, v.Locale.Code, v.Name)
	}

	item.KeyPoints = buildKeyPoints(resp.KeyPoints)

	return item
}

func buildKeyPoints(rows []keypoint.LessonKeyPointModel) []lesson2.Locale {
	enList := make([]string, 0)
	ruList := make([]string, 0)
	kkList := make([]string, 0)

	for _, v := range rows {
		switch strings.ToLower(v.Locale.Code) {
		case "en":
			enList = append(enList, v.Name)
		case "ru":
			ruList = append(ruList, v.Name)
		case "kk":
			kkList = append(kkList, v.Name)
		}
	}

	maxLen := len(enList)
	if len(ruList) > maxLen {
		maxLen = len(ruList)
	}
	if len(kkList) > maxLen {
		maxLen = len(kkList)
	}

	result := make([]lesson2.Locale, maxLen)

	for i := 0; i < maxLen; i++ {
		if i < len(enList) {
			result[i].EN = enList[i]
		}
		if i < len(ruList) {
			result[i].RU = ruList[i]
		}
		if i < len(kkList) {
			result[i].KK = kkList[i]
		}
	}

	return result
}

func setLocaleValue(dst *lesson2.Locale, localeCode, value string) {
	switch strings.ToLower(localeCode) {
	case "en":
		dst.EN = value
	case "ru":
		dst.RU = value
	case "kk":
		dst.KK = value
	}
}
func stringPtr(v string) *string {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return &v
}

func convertLessonRequest(req lesson2.LessonRequest, lessonID uuid.UUID, locales []locale.Locale) *lesson.LessonModel {
	localesMap := make(map[string]uuid.UUID, 0)
	for i := range locales {
		localesMap[locales[i].Code] = locales[i].ID
	}
	return &lesson.LessonModel{
		ModuleID:        req.ModuleID,
		DurationMinutes: req.DurationMinutes,
		XPReward:        req.XPReward,
		CodeSnippet:     stringPtr(req.CodeSnippet),
		ExampleOutput:   stringPtr(req.ExampleOutput),

		Titles:         buildLessonTitles(req.Title, localesMap, lessonID),
		Summaries:      buildLessonSummaries(req.Summary, localesMap, lessonID),
		Outcomes:       buildLessonOutcomes(req.OutCome, localesMap, lessonID),
		TheoryContents: buildLessonTheoryContents(req.TheoryContent, localesMap, lessonID),
		KeyPoints:      buildLessonKeyPoints(req.KeyPoints, localesMap, lessonID),
	}
}

func buildLessonTitles(src lesson2.Locale, localesMap map[string]uuid.UUID, lessonID uuid.UUID) []title.LessonTitleModel {
	titles := make([]title.LessonTitleModel, 0)
	titles = append(titles, title.LessonTitleModel{
		ID:       uuid.New(),
		Name:     src.EN,
		LocaleID: localesMap["en"],
		LessonID: lessonID,
	})
	titles = append(titles, title.LessonTitleModel{
		ID:       uuid.New(),
		Name:     src.RU,
		LocaleID: localesMap["ru"],
		LessonID: lessonID,
	})
	titles = append(titles, title.LessonTitleModel{
		ID:       uuid.New(),
		Name:     src.KK,
		LocaleID: localesMap["kk"],
		LessonID: lessonID,
	})
	return titles
}

func buildLessonSummaries(src lesson2.Locale, localesMap map[string]uuid.UUID, lessonID uuid.UUID) []summary.LessonSummaryModel {
	summaries := make([]summary.LessonSummaryModel, 0)
	summaries = append(summaries, summary.LessonSummaryModel{
		ID:       uuid.New(),
		Name:     src.EN,
		LocaleID: localesMap["en"],
		LessonID: lessonID,
	})
	summaries = append(summaries, summary.LessonSummaryModel{
		ID:       uuid.New(),
		Name:     src.RU,
		LocaleID: localesMap["ru"],
		LessonID: lessonID,
	})
	summaries = append(summaries, summary.LessonSummaryModel{
		ID:       uuid.New(),
		Name:     src.KK,
		LocaleID: localesMap["kk"],
		LessonID: lessonID,
	})
	return summaries
}

func buildLessonOutcomes(src lesson2.Locale, localesMap map[string]uuid.UUID, lessonID uuid.UUID) []outcome.LessonOutcomeModel {
	outcomes := make([]outcome.LessonOutcomeModel, 0)
	outcomes = append(outcomes, outcome.LessonOutcomeModel{
		ID:       uuid.New(),
		Name:     src.EN,
		LocaleID: localesMap["en"],
		LessonID: lessonID,
	})

	outcomes = append(outcomes, outcome.LessonOutcomeModel{
		ID:       uuid.New(),
		Name:     src.RU,
		LocaleID: localesMap["ru"],
		LessonID: lessonID,
	})

	outcomes = append(outcomes, outcome.LessonOutcomeModel{
		ID:       uuid.New(),
		Name:     src.KK,
		LocaleID: localesMap["kk"],
		LessonID: lessonID,
	})
	return outcomes
}

func buildLessonTheoryContents(src lesson2.Locale, localesMap map[string]uuid.UUID, lessonID uuid.UUID) []theorycontent.LessonTheoryContentModel {
	theoryContents := make([]theorycontent.LessonTheoryContentModel, 0)
	theoryContents = append(theoryContents, theorycontent.LessonTheoryContentModel{
		ID:       uuid.New(),
		Name:     src.EN,
		LocaleID: localesMap["en"],
		LessonID: lessonID,
	})
	theoryContents = append(theoryContents, theorycontent.LessonTheoryContentModel{
		ID:       uuid.New(),
		Name:     src.RU,
		LocaleID: localesMap["ru"],
		LessonID: lessonID,
	})
	theoryContents = append(theoryContents, theorycontent.LessonTheoryContentModel{
		ID:       uuid.New(),
		Name:     src.KK,
		LocaleID: localesMap["kk"],
		LessonID: lessonID,
	})
	return theoryContents
}

func buildLessonKeyPoints(src []lesson2.Locale, localesMap map[string]uuid.UUID, lessonID uuid.UUID) []keypoint.LessonKeyPointModel {
	keypoints := make([]keypoint.LessonKeyPointModel, 0)
	for i := range src {
		keypoints = append(keypoints, keypoint.LessonKeyPointModel{
			ID:       uuid.New(),
			Name:     src[i].EN,
			LocaleID: localesMap["en"],
			LessonID: lessonID,
		})
		keypoints = append(keypoints, keypoint.LessonKeyPointModel{
			ID:       uuid.New(),
			Name:     src[i].RU,
			LocaleID: localesMap["ru"],
			LessonID: lessonID,
		})
		keypoints = append(keypoints, keypoint.LessonKeyPointModel{
			ID:       uuid.New(),
			Name:     src[i].KK,
			LocaleID: localesMap["kk"],
			LessonID: lessonID,
		})
	}
	return keypoints
}
