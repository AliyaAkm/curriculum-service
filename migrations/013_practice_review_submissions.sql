CREATE TABLE IF NOT EXISTS practice_review_submissions (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    practice_id uuid NOT NULL REFERENCES practice_tasks(id) ON DELETE CASCADE,
    student_id uuid NOT NULL,
    course_id uuid NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    lesson_id uuid NOT NULL REFERENCES course_lessons(id) ON DELETE CASCADE,
    status text NOT NULL DEFAULT 'submitted',
    code text NOT NULL,
    language text NOT NULL,
    output text NOT NULL DEFAULT '',
    error text NOT NULL DEFAULT '',
    error_type text NOT NULL DEFAULT '',
    exit_code integer,
    duration_ms integer,
    teacher_comment text NOT NULL DEFAULT '',
    reviewed_by uuid,
    reviewed_at timestamptz,
    attempt_number integer NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT practice_review_submissions_status_check
        CHECK (status IN ('submitted', 'in_review', 'changes_requested', 'approved')),
    CONSTRAINT practice_review_submissions_attempt_positive CHECK (attempt_number > 0),
    CONSTRAINT practice_review_submissions_code_not_blank CHECK (length(trim(code)) > 0),
    CONSTRAINT practice_review_submissions_language_not_blank CHECK (length(trim(language)) > 0)
);

CREATE INDEX IF NOT EXISTS idx_practice_review_submissions_student_created
ON practice_review_submissions(student_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_practice_review_submissions_teacher_queue
ON practice_review_submissions(course_id, status, created_at);

CREATE INDEX IF NOT EXISTS idx_practice_review_submissions_practice_student
ON practice_review_submissions(practice_id, student_id, attempt_number DESC);

CREATE UNIQUE INDEX IF NOT EXISTS uq_practice_review_submissions_active
ON practice_review_submissions(student_id, practice_id)
WHERE status IN ('submitted', 'in_review');

CREATE TABLE IF NOT EXISTS student_practice_progress (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id uuid NOT NULL,
    practice_id uuid NOT NULL REFERENCES practice_tasks(id) ON DELETE CASCADE,
    course_id uuid NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    lesson_id uuid NOT NULL REFERENCES course_lessons(id) ON DELETE CASCADE,
    status text NOT NULL DEFAULT 'in_progress',
    started_at timestamptz NOT NULL DEFAULT now(),
    completed_at timestamptz,
    last_attempt_at timestamptz,
    attempts_count integer NOT NULL DEFAULT 0,
    approved_submission_id uuid REFERENCES practice_review_submissions(id) ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT student_practice_progress_status_check
        CHECK (status IN ('in_progress', 'submitted', 'changes_requested', 'completed')),
    CONSTRAINT student_practice_progress_attempts_non_negative CHECK (attempts_count >= 0),
    UNIQUE (student_id, practice_id)
);

CREATE INDEX IF NOT EXISTS idx_student_practice_progress_student_course
ON student_practice_progress(student_id, course_id, status);

CREATE INDEX IF NOT EXISTS idx_student_practice_progress_course_status
ON student_practice_progress(course_id, status, last_attempt_at DESC);

WITH attempt_progress AS (
    SELECT
        cea.user_id AS student_id,
        pt.id AS practice_id,
        cm.course_id AS course_id,
        pt.lesson_id AS lesson_id,
        MIN(cea.created_at) AS started_at,
        MAX(cea.created_at) AS last_attempt_at,
        COUNT(*)::integer AS attempts_count,
        MIN(cea.created_at) FILTER (WHERE cea.run_type = 'submit' AND cea.passed = TRUE) AS first_success_at
    FROM code_execution_attempts cea
    INNER JOIN practice_tasks pt ON pt.id::text = cea.practice_id
    INNER JOIN course_lessons cl ON cl.id = pt.lesson_id
    INNER JOIN course_modules cm ON cm.id = cl.module_id
    GROUP BY cea.user_id, pt.id, cm.course_id, pt.lesson_id
),
awards AS (
    SELECT
        user_id AS student_id,
        practice_id,
        MIN(created_at) AS awarded_at
    FROM practice_xp_awards
    GROUP BY user_id, practice_id
)
INSERT INTO student_practice_progress (
    id,
    student_id,
    practice_id,
    course_id,
    lesson_id,
    status,
    started_at,
    completed_at,
    last_attempt_at,
    attempts_count
)
SELECT
    uuid_generate_v4(),
    ap.student_id,
    ap.practice_id,
    ap.course_id,
    ap.lesson_id,
    CASE WHEN COALESCE(ap.first_success_at, awards.awarded_at) IS NOT NULL THEN 'completed' ELSE 'in_progress' END,
    ap.started_at,
    COALESCE(ap.first_success_at, awards.awarded_at),
    ap.last_attempt_at,
    ap.attempts_count
FROM attempt_progress ap
LEFT JOIN awards ON awards.student_id = ap.student_id AND awards.practice_id = ap.practice_id
ON CONFLICT (student_id, practice_id) DO NOTHING;

ALTER TABLE practice_xp_awards
ADD COLUMN IF NOT EXISTS submission_id uuid REFERENCES practice_review_submissions(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_practice_xp_awards_submission_id
ON practice_xp_awards(submission_id);
