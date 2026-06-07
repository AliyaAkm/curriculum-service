CREATE TABLE IF NOT EXISTS practice_tasks (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    lesson_id uuid NOT NULL REFERENCES course_lessons(id) ON DELETE CASCADE,
    position integer NOT NULL DEFAULT 1,
    title text NOT NULL DEFAULT '',
    description text NOT NULL DEFAULT '',
    language text NOT NULL DEFAULT 'java',
    starter_code text NOT NULL DEFAULT '',
    expected_output text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT practice_tasks_position_positive CHECK (position > 0),
    CONSTRAINT practice_tasks_language_not_blank CHECK (length(trim(language)) > 0),
    UNIQUE (lesson_id, position)
);

CREATE INDEX IF NOT EXISTS idx_practice_tasks_lesson_id
ON practice_tasks(lesson_id);
