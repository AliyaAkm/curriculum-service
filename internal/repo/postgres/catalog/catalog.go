package postgres

import (
	"context"
	"curriculum-service/internal/domain/category"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CatalogRepo struct {
	pool *pgxpool.Pool
}

func NewCatalogRepo(pool *pgxpool.Pool) *CatalogRepo {
	return &CatalogRepo{pool: pool}
}

func (r *CatalogRepo) SearchCourses(ctx context.Context, filter category.CourseSearchFilter) ([]category.CourseCard, int, error) {
	builder := newCourseQueryBuilder(filter)
	countWhereSQL := builder.whereSQL()
	if builder.exactQuery != "" && builder.prefixQuery != "" {
		countWhereSQL += fmt.Sprintf("\nAND %s::text IS NOT NULL\nAND %s::text IS NOT NULL", builder.exactQuery, builder.prefixQuery)
	}

	countSQL := `
SELECT COUNT(*)
FROM courses c
JOIN course_localizations cl ON cl.course_id = c.id AND cl.locale = $1
WHERE ` + countWhereSQL

	var total int
	if err := r.pool.QueryRow(ctx, countSQL, builder.args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limitPH := builder.add(filter.PageSize)
	offsetPH := builder.add((filter.Page - 1) * filter.PageSize)

	searchSQL := fmt.Sprintf(`
SELECT
	c.id,
	c.slug,
	c.status,
	c.level,
	c.duration_category,
	c.expected_hours,
	c.rating,
	c.rating_count,
	c.students_count,
	c.lessons_count,
	c.has_certificate,
	c.cover_image_url,
	c.author_name,
	cl.title,
	cl.subtitle,
	cl.short_description,
	COALESCE(array_remove(array_agg(DISTINCT t.slug), NULL), '{}') AS topic_slugs,
	COALESCE(array_remove(array_agg(DISTINCT COALESCE(tl.name, t.slug)), NULL), '{}') AS topic_names,
	COALESCE(array_remove(array_agg(DISTINCT COALESCE(tagl.name, tag.slug)), NULL), '{}') AS tags,
	c.published_at
FROM courses c
JOIN course_localizations cl ON cl.course_id = c.id AND cl.locale = $1
LEFT JOIN course_topics ct ON ct.course_id = c.id
LEFT JOIN topics t ON t.id = ct.topic_id
LEFT JOIN topic_localizations tl ON tl.topic_id = t.id AND tl.locale = $1
LEFT JOIN course_tags ctag ON ctag.course_id = c.id
LEFT JOIN tags tag ON tag.id = ctag.tag_id
LEFT JOIN tag_localizations tagl ON tagl.tag_id = tag.id AND tagl.locale = $1
WHERE %s
GROUP BY
	c.id,
	c.slug,
	c.status,
	c.level,
	c.duration_category,
	c.expected_hours,
	c.rating,
	c.rating_count,
	c.students_count,
	c.lessons_count,
	c.has_certificate,
	c.cover_image_url,
	c.author_name,
	cl.title,
	cl.subtitle,
	cl.short_description,
	c.published_at
ORDER BY %s
LIMIT %s OFFSET %s
`, builder.whereSQL(), builder.orderBySQL(), limitPH, offsetPH)

	rows, err := r.pool.Query(ctx, searchSQL, builder.args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]category.CourseCard, 0, filter.PageSize)
	for rows.Next() {
		var (
			item             category.CourseCard
			status           string
			level            string
			durationCategory string
			publishedAt      time.Time
			topicSlugs       []string
			topicNames       []string
			tags             []string
		)

		if err := rows.Scan(
			&item.ID,
			&item.Slug,
			&status,
			&level,
			&durationCategory,
			&item.ExpectedHours,
			&item.Rating,
			&item.RatingCount,
			&item.StudentsCount,
			&item.LessonsCount,
			&item.HasCertificate,
			&item.CoverImageURL,
			&item.AuthorName,
			&item.Title,
			&item.Subtitle,
			&item.ShortDescription,
			&topicSlugs,
			&topicNames,
			&tags,
			&publishedAt,
		); err != nil {
			return nil, 0, err
		}

		item.Status = category.CourseStatus(status)
		item.Level = category.CourseLevel(level)
		item.DurationCategory = category.DurationCategory(durationCategory)
		item.TopicSlugs = topicSlugs
		item.TopicNames = topicNames
		item.Tags = tags
		item.PublishedAt = publishedAt

		items = append(items, item)
	}

	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}

	return items, total, nil
}

func (r *CatalogRepo) GetFilterOptions(ctx context.Context, locale category.Locale) (category.FilterOptions, error) {
	topics, err := r.listTopicOptions(ctx, locale)
	if err != nil {
		return category.FilterOptions{}, err
	}

	levels, err := r.listLevelOptions(ctx)
	if err != nil {
		return category.FilterOptions{}, err
	}

	durations, err := r.listDurationOptions(ctx)
	if err != nil {
		return category.FilterOptions{}, err
	}

	certificateAvailable, err := r.listCertificateAvailability(ctx)
	if err != nil {
		return category.FilterOptions{}, err
	}

	return category.FilterOptions{
		Topics:               topics,
		Levels:               levels,
		Durations:            durations,
		CertificateAvailable: certificateAvailable,
	}, nil
}

func (r *CatalogRepo) listTopicOptions(ctx context.Context, locale category.Locale) ([]category.TopicFilterOption, error) {
	const query = `
SELECT
	t.slug,
	COALESCE(tl.name, t.slug) AS name,
	COUNT(DISTINCT c.id) AS courses_count
FROM topics t
LEFT JOIN topic_localizations tl ON tl.topic_id = t.id AND tl.locale = $1
LEFT JOIN course_topics ct ON ct.topic_id = t.id
LEFT JOIN courses c ON c.id = ct.course_id AND c.status = 'published'
GROUP BY t.slug, COALESCE(tl.name, t.slug)
HAVING COUNT(DISTINCT c.id) > 0
ORDER BY LOWER(COALESCE(tl.name, t.slug))
`

	rows, err := r.pool.Query(ctx, query, string(locale))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]category.TopicFilterOption, 0)
	for rows.Next() {
		var item category.TopicFilterOption
		if err := rows.Scan(&item.Slug, &item.Name, &item.CoursesCount); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return items, nil
}

func (r *CatalogRepo) listLevelOptions(ctx context.Context) ([]category.FilterValueOption, error) {
	const query = `
SELECT level, COUNT(*)
FROM courses
WHERE status = 'published'
GROUP BY level
`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]category.FilterValueOption, 0)
	for rows.Next() {
		var item category.FilterValueOption
		if err := rows.Scan(&item.Value, &item.CoursesCount); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return items, nil
}

func (r *CatalogRepo) listDurationOptions(ctx context.Context) ([]category.FilterValueOption, error) {
	const query = `
SELECT duration_category, COUNT(*)
FROM courses
WHERE status = 'published'
GROUP BY duration_category
`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]category.FilterValueOption, 0)
	for rows.Next() {
		var item category.FilterValueOption
		if err := rows.Scan(&item.Value, &item.CoursesCount); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return items, nil
}

func (r *CatalogRepo) listCertificateAvailability(ctx context.Context) (bool, error) {
	const query = `
SELECT EXISTS(
	SELECT 1
	FROM courses
	WHERE status = 'published'
	  AND has_certificate = TRUE
)
`

	var available bool
	if err := r.pool.QueryRow(ctx, query).Scan(&available); err != nil {
		return false, err
	}

	return available, nil
}

type courseQueryBuilder struct {
	args        []any
	conditions  []string
	exactQuery  string
	prefixQuery string
	likeQuery   string
}

func newCourseQueryBuilder(filter category.CourseSearchFilter) *courseQueryBuilder {
	builder := &courseQueryBuilder{
		args:       []any{string(filter.Locale)},
		conditions: []string{"c.status = 'published'"},
	}

	if filter.Query != "" {
		builder.exactQuery = builder.add(filter.Query)
		builder.prefixQuery = builder.add(filter.Query + "%")
		builder.likeQuery = builder.add("%" + filter.Query + "%")

		builder.conditions = append(builder.conditions, fmt.Sprintf(`(
	cl.title ILIKE %[1]s
	OR cl.short_description ILIKE %[1]s
	OR EXISTS (
		SELECT 1
		FROM course_localizations search_cl
		WHERE search_cl.course_id = c.id
		  AND search_cl.locale = $1
		  AND (
			search_cl.title ILIKE %[1]s
			OR search_cl.short_description ILIKE %[1]s
			OR search_cl.description ILIKE %[1]s
		  )
	)
	OR EXISTS (
		SELECT 1
		FROM course_topics search_ct
		JOIN topics search_t ON search_t.id = search_ct.topic_id
		LEFT JOIN topic_localizations search_tl ON search_tl.topic_id = search_t.id AND search_tl.locale = $1
		WHERE search_ct.course_id = c.id
		  AND (
			search_t.slug ILIKE %[1]s
			OR COALESCE(search_tl.name, search_t.slug) ILIKE %[1]s
		  )
	)
	OR EXISTS (
		SELECT 1
		FROM course_tags search_ctag
		JOIN tags search_tag ON search_tag.id = search_ctag.tag_id
		LEFT JOIN tag_localizations search_tagl ON search_tagl.tag_id = search_tag.id AND search_tagl.locale = $1
		WHERE search_ctag.course_id = c.id
		  AND (
			search_tag.slug ILIKE %[1]s
			OR COALESCE(search_tagl.name, search_tag.slug) ILIKE %[1]s
		  )
	)
)`, builder.likeQuery))
	}

	if len(filter.TopicSlugs) > 0 {
		topicsPH := builder.add(filter.TopicSlugs)
		builder.conditions = append(builder.conditions, fmt.Sprintf(`EXISTS (
	SELECT 1
	FROM course_topics filter_ct
	JOIN topics filter_t ON filter_t.id = filter_ct.topic_id
	WHERE filter_ct.course_id = c.id
	  AND filter_t.slug = ANY(%s)
)`, topicsPH))
	}

	if len(filter.Levels) > 0 {
		levelsPH := builder.add(levelsToStrings(filter.Levels))
		builder.conditions = append(builder.conditions, fmt.Sprintf("c.level = ANY(%s)", levelsPH))
	}

	if filter.MinRating > 0 {
		minRatingPH := builder.add(filter.MinRating)
		builder.conditions = append(builder.conditions, fmt.Sprintf("c.rating >= %s", minRatingPH))
	}

	if len(filter.Durations) > 0 {
		durationsPH := builder.add(durationsToStrings(filter.Durations))
		builder.conditions = append(builder.conditions, fmt.Sprintf("c.duration_category = ANY(%s)", durationsPH))
	}

	if filter.WithCertificate != nil && *filter.WithCertificate {
		builder.conditions = append(builder.conditions, "c.has_certificate = TRUE")
	}

	return builder
}

func (b *courseQueryBuilder) add(value any) string {
	b.args = append(b.args, value)
	return fmt.Sprintf("$%d", len(b.args))
}

func (b *courseQueryBuilder) whereSQL() string {
	return strings.Join(b.conditions, "\nAND ")
}

func (b *courseQueryBuilder) orderBySQL() string {
	if b.likeQuery == "" {
		return "c.rating DESC, c.published_at DESC NULLS LAST, cl.title ASC"
	}

	return fmt.Sprintf(`
CASE
	WHEN LOWER(cl.title) = LOWER(%s) THEN 0
	WHEN cl.title ILIKE %s THEN 1
	WHEN cl.title ILIKE %s THEN 2
	ELSE 3
END,
c.rating DESC,
c.published_at DESC NULLS LAST,
cl.title ASC
`, b.exactQuery, b.prefixQuery, b.likeQuery)
}

func levelsToStrings(levels []category.CourseLevel) []string {
	items := make([]string, 0, len(levels))
	for _, level := range levels {
		items = append(items, string(level))
	}
	return items
}

func durationsToStrings(durations []category.DurationCategory) []string {
	items := make([]string, 0, len(durations))
	for _, duration := range durations {
		items = append(items, string(duration))
	}
	return items
}
