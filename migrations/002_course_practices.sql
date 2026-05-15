CREATE TABLE IF NOT EXISTS course_practices (
    id uuid PRIMARY KEY,
    lesson_id uuid NOT NULL UNIQUE REFERENCES course_lessons(id) ON DELETE CASCADE,
    position integer NOT NULL DEFAULT 1,
    title_en text NOT NULL DEFAULT '',
    title_ru text NOT NULL DEFAULT '',
    title_kk text NOT NULL DEFAULT '',
    summary_en text NOT NULL DEFAULT '',
    summary_ru text NOT NULL DEFAULT '',
    summary_kk text NOT NULL DEFAULT '',
    brief_en text NOT NULL DEFAULT '',
    brief_ru text NOT NULL DEFAULT '',
    brief_kk text NOT NULL DEFAULT '',
    starter_code text NOT NULL DEFAULT '',
    success_criteria jsonb NOT NULL DEFAULT '[]'::jsonb,
    knowledge_checks jsonb NOT NULL DEFAULT '[]'::jsonb,
    prompt_suggestion_en text NOT NULL DEFAULT '',
    prompt_suggestion_ru text NOT NULL DEFAULT '',
    prompt_suggestion_kk text NOT NULL DEFAULT '',
    xp_reward integer NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT course_practices_position_positive CHECK (position > 0),
    CONSTRAINT course_practices_xp_reward_non_negative CHECK (xp_reward >= 0)
);

CREATE INDEX IF NOT EXISTS idx_course_practices_lesson_id
    ON course_practices(lesson_id);
