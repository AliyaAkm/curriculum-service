CREATE TABLE IF NOT EXISTS lesson_quiz_attempts (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    quiz_id uuid NOT NULL REFERENCES lesson_quizzes(id) ON DELETE CASCADE,
    user_id uuid NOT NULL,
    selected_answer_index integer NOT NULL,
    is_correct boolean NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_lesson_quiz_attempts_quiz_id
    ON lesson_quiz_attempts(quiz_id);

CREATE INDEX IF NOT EXISTS idx_lesson_quiz_attempts_user_id
    ON lesson_quiz_attempts(user_id);

CREATE UNIQUE INDEX IF NOT EXISTS uq_lesson_quiz_attempts_user_quiz
    ON lesson_quiz_attempts(user_id, quiz_id);
