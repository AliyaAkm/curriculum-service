ALTER TABLE course_lessons
ADD COLUMN IF NOT EXISTS position integer NOT NULL DEFAULT 0;

UPDATE course_lessons AS lesson
SET position = numbered.position
FROM (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY module_id
            ORDER BY created_at ASC, id ASC
        )::integer AS position
    FROM course_lessons
) AS numbered
WHERE lesson.id = numbered.id
  AND lesson.position = 0;
