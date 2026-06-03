CREATE TABLE IF NOT EXISTS course_certificates (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL,
    course_id uuid NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    certificate_number text NOT NULL UNIQUE,
    pdf_object_key text NOT NULL,
    issued_at timestamptz NOT NULL,
    completed_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_course_certificates_user_course
ON course_certificates(user_id, course_id);

CREATE INDEX IF NOT EXISTS idx_course_certificates_user_id
ON course_certificates(user_id);
