ALTER TABLE course_subscription
ADD COLUMN IF NOT EXISTS started_at timestamptz,
ADD COLUMN IF NOT EXISTS last_activity_at timestamptz,
ADD COLUMN IF NOT EXISTS completed_at timestamptz,
ADD COLUMN IF NOT EXISTS current_lesson_id uuid REFERENCES course_lessons(id);

ALTER TABLE lesson_quiz_attempts
ADD COLUMN IF NOT EXISTS updated_at timestamptz NOT NULL DEFAULT now();

WITH duplicate_groups AS (
    SELECT
        user_id,
        course_id,
        MIN(id::text)::uuid AS keep_id,
        MIN(started_at) AS started_at,
        MAX(last_activity_at) AS last_activity_at,
        MIN(completed_at) AS completed_at
    FROM course_subscription
    GROUP BY user_id, course_id
    HAVING COUNT(*) > 1
),
latest_current_lesson AS (
    SELECT DISTINCT ON (cs.user_id, cs.course_id)
        cs.user_id,
        cs.course_id,
        cs.current_lesson_id
    FROM course_subscription cs
    INNER JOIN duplicate_groups dg
        ON dg.user_id = cs.user_id
       AND dg.course_id = cs.course_id
    WHERE cs.current_lesson_id IS NOT NULL
    ORDER BY
        cs.user_id,
        cs.course_id,
        cs.last_activity_at DESC NULLS LAST,
        cs.started_at DESC NULLS LAST,
        cs.id
)
UPDATE course_subscription keep
SET started_at = COALESCE(keep.started_at, dg.started_at),
    last_activity_at = CASE
        WHEN keep.last_activity_at IS NULL THEN dg.last_activity_at
        WHEN dg.last_activity_at IS NULL THEN keep.last_activity_at
        WHEN keep.last_activity_at >= dg.last_activity_at THEN keep.last_activity_at
        ELSE dg.last_activity_at
    END,
    completed_at = COALESCE(keep.completed_at, dg.completed_at),
    current_lesson_id = COALESCE(keep.current_lesson_id, latest_current_lesson.current_lesson_id)
FROM duplicate_groups dg
LEFT JOIN latest_current_lesson
    ON latest_current_lesson.user_id = dg.user_id
   AND latest_current_lesson.course_id = dg.course_id
WHERE keep.id = dg.keep_id;

DELETE FROM course_subscription cs
USING (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY user_id, course_id
            ORDER BY id
        ) AS rn
    FROM course_subscription
) ranked
WHERE cs.id = ranked.id
  AND ranked.rn > 1;

CREATE UNIQUE INDEX IF NOT EXISTS uq_course_subscription_user_course
ON course_subscription(user_id, course_id);

DELETE FROM user_course_points ucp
USING (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY user_id, lesson_id
            ORDER BY updated_at DESC NULLS LAST, created_at DESC NULLS LAST, id ASC
        ) AS rn
    FROM user_course_points
) ranked
WHERE ucp.id = ranked.id
  AND ranked.rn > 1;

CREATE UNIQUE INDEX IF NOT EXISTS uq_user_course_points_user_lesson
ON user_course_points(user_id, lesson_id);

CREATE INDEX IF NOT EXISTS idx_course_subscription_user_id
ON course_subscription(user_id);

CREATE INDEX IF NOT EXISTS idx_user_course_points_user_id
ON user_course_points(user_id);

CREATE INDEX IF NOT EXISTS idx_user_course_points_lesson_id
ON user_course_points(lesson_id);
