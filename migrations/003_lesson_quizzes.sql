CREATE TABLE IF NOT EXISTS lesson_quizzes (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    lesson_id uuid NOT NULL REFERENCES course_lessons(id) ON DELETE CASCADE,
    position integer NOT NULL DEFAULT 1,
    correct_answer_index integer NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT lesson_quizzes_position_positive CHECK (position > 0),
    CONSTRAINT lesson_quizzes_correct_answer_index_non_negative CHECK (correct_answer_index >= 0)
);

CREATE TABLE IF NOT EXISTS lesson_quiz_texts (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    quiz_id uuid NOT NULL REFERENCES lesson_quizzes(id) ON DELETE CASCADE,
    locale_id uuid NOT NULL REFERENCES course_locales(id),
    question text NOT NULL DEFAULT '',
    explanation text NOT NULL DEFAULT '',
    UNIQUE (quiz_id, locale_id)
);

CREATE TABLE IF NOT EXISTS lesson_quiz_options (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    quiz_id uuid NOT NULL REFERENCES lesson_quizzes(id) ON DELETE CASCADE,
    position integer NOT NULL,
    UNIQUE (quiz_id, position),
    CONSTRAINT lesson_quiz_options_position_positive CHECK (position > 0)
);

CREATE TABLE IF NOT EXISTS lesson_quiz_option_texts (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    option_id uuid NOT NULL REFERENCES lesson_quiz_options(id) ON DELETE CASCADE,
    locale_id uuid NOT NULL REFERENCES course_locales(id),
    text text NOT NULL DEFAULT '',
    UNIQUE (option_id, locale_id)
);

CREATE INDEX IF NOT EXISTS idx_lesson_quizzes_lesson_id
    ON lesson_quizzes(lesson_id);

CREATE INDEX IF NOT EXISTS idx_lesson_quiz_texts_quiz_id
    ON lesson_quiz_texts(quiz_id);

CREATE INDEX IF NOT EXISTS idx_lesson_quiz_options_quiz_id
    ON lesson_quiz_options(quiz_id);

CREATE INDEX IF NOT EXISTS idx_lesson_quiz_option_texts_option_id
    ON lesson_quiz_option_texts(option_id);
