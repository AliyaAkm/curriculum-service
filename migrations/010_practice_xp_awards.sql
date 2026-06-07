ALTER TABLE practice_tasks
ADD COLUMN IF NOT EXISTS xp_reward integer NOT NULL DEFAULT 25;

CREATE TABLE IF NOT EXISTS practice_xp_awards (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    practice_id uuid NOT NULL REFERENCES practice_tasks(id) ON DELETE CASCADE,
    course_id uuid REFERENCES courses(id) ON DELETE SET NULL,
    lesson_id uuid REFERENCES course_lessons(id) ON DELETE SET NULL,
    xp integer NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT practice_xp_awards_xp_positive CHECK (xp > 0),
    UNIQUE (user_id, practice_id)
);

ALTER TABLE code_execution_attempts
ADD COLUMN IF NOT EXISTS xp_awarded integer NOT NULL DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_practice_xp_awards_user_created
ON practice_xp_awards(user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_practice_xp_awards_practice_id
ON practice_xp_awards(practice_id);

CREATE INDEX IF NOT EXISTS idx_practice_xp_awards_course_id
ON practice_xp_awards(course_id);
