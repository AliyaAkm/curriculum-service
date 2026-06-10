CREATE TABLE IF NOT EXISTS student_lesson_progress (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    lesson_id uuid NOT NULL REFERENCES course_lessons(id) ON DELETE CASCADE,
    course_id uuid NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    theory_completed_at timestamptz,
    completed_at timestamptz,
    last_activity_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (user_id, lesson_id)
);

CREATE INDEX IF NOT EXISTS idx_student_lesson_progress_user_course
ON student_lesson_progress(user_id, course_id, completed_at);

CREATE INDEX IF NOT EXISTS idx_student_lesson_progress_lesson
ON student_lesson_progress(lesson_id, completed_at);

INSERT INTO student_lesson_progress (
    id,
    user_id,
    lesson_id,
    course_id,
    theory_completed_at,
    completed_at,
    last_activity_at,
    created_at,
    updated_at
)
SELECT
    uuid_generate_v4(),
    ucp.user_id,
    ucp.lesson_id,
    cm.course_id,
    ucp.created_at,
    ucp.created_at,
    ucp.created_at,
    ucp.created_at,
    COALESCE(ucp.updated_at, ucp.created_at)
FROM user_course_points ucp
INNER JOIN course_lessons cl ON cl.id = ucp.lesson_id
INNER JOIN course_modules cm ON cm.id = cl.module_id
ON CONFLICT (user_id, lesson_id) DO NOTHING;
