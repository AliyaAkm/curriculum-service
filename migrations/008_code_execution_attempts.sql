CREATE TABLE IF NOT EXISTS code_execution_attempts (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    course_id uuid REFERENCES courses(id) ON DELETE SET NULL,
    lesson_id uuid REFERENCES course_lessons(id) ON DELETE SET NULL,
    practice_id text NOT NULL,
    run_type text NOT NULL DEFAULT 'run',
    language text NOT NULL,
    passed boolean NOT NULL DEFAULT false,
    error_type text NOT NULL DEFAULT '',
    error_message text NOT NULL DEFAULT '',
    output text NOT NULL DEFAULT '',
    duration_ms integer NOT NULL DEFAULT 0,
    code_hash text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT code_execution_attempts_run_type_check CHECK (run_type IN ('run', 'submit')),
    CONSTRAINT code_execution_attempts_duration_non_negative CHECK (duration_ms >= 0)
);

CREATE INDEX IF NOT EXISTS idx_code_execution_attempts_user_created
ON code_execution_attempts(user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_code_execution_attempts_practice_id
ON code_execution_attempts(practice_id);

CREATE INDEX IF NOT EXISTS idx_code_execution_attempts_lesson_id
ON code_execution_attempts(lesson_id);

CREATE INDEX IF NOT EXISTS idx_code_execution_attempts_course_id
ON code_execution_attempts(course_id);
