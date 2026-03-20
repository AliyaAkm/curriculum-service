CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    level TEXT NOT NULL CHECK (level IN ('beginner', 'intermediate', 'advanced')),
    duration_category TEXT NOT NULL CHECK (duration_category IN ('quick', 'focused', 'deep')),
    expected_hours INTEGER NOT NULL CHECK (expected_hours > 0),
    rating DOUBLE PRECISION NOT NULL DEFAULT 0 CHECK (rating >= 0 AND rating <= 5),
    rating_count INTEGER NOT NULL DEFAULT 0 CHECK (rating_count >= 0),
    students_count INTEGER NOT NULL DEFAULT 0 CHECK (students_count >= 0),
    lessons_count INTEGER NOT NULL DEFAULT 0 CHECK (lessons_count >= 0),
    has_certificate BOOLEAN NOT NULL DEFAULT FALSE,
    cover_image_url TEXT NOT NULL DEFAULT '',
    author_name TEXT NOT NULL DEFAULT '',
    published_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT published_courses_have_date CHECK (status <> 'published' OR published_at IS NOT NULL)
);

CREATE TABLE IF NOT EXISTS course_localizations (
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    locale TEXT NOT NULL CHECK (locale IN ('ru', 'en', 'kz')),
    title TEXT NOT NULL,
    subtitle TEXT NOT NULL DEFAULT '',
    short_description TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    syllabus JSONB NOT NULL DEFAULT '[]'::jsonb,
    PRIMARY KEY (course_id, locale)
);

CREATE TABLE IF NOT EXISTS topics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS topic_localizations (
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    locale TEXT NOT NULL CHECK (locale IN ('ru', 'en', 'kz')),
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (topic_id, locale)
);

CREATE TABLE IF NOT EXISTS course_topics (
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE RESTRICT,
    PRIMARY KEY (course_id, topic_id)
);

CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tag_localizations (
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    locale TEXT NOT NULL CHECK (locale IN ('ru', 'en', 'kz')),
    name TEXT NOT NULL,
    PRIMARY KEY (tag_id, locale)
);

CREATE TABLE IF NOT EXISTS course_tags (
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE RESTRICT,
    PRIMARY KEY (course_id, tag_id)
);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_courses_set_updated_at ON courses;
CREATE TRIGGER trg_courses_set_updated_at
BEFORE UPDATE ON courses
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE INDEX IF NOT EXISTS idx_courses_status ON courses(status);
CREATE INDEX IF NOT EXISTS idx_courses_level ON courses(level);
CREATE INDEX IF NOT EXISTS idx_courses_duration_category ON courses(duration_category);
CREATE INDEX IF NOT EXISTS idx_courses_has_certificate ON courses(has_certificate);
CREATE INDEX IF NOT EXISTS idx_courses_published_at ON courses(published_at DESC);
CREATE INDEX IF NOT EXISTS idx_course_topics_topic_id ON course_topics(topic_id);
CREATE INDEX IF NOT EXISTS idx_course_tags_tag_id ON course_tags(tag_id);

CREATE INDEX IF NOT EXISTS idx_course_localizations_locale ON course_localizations(locale);
CREATE INDEX IF NOT EXISTS idx_topic_localizations_locale ON topic_localizations(locale);
CREATE INDEX IF NOT EXISTS idx_tag_localizations_locale ON tag_localizations(locale);

CREATE INDEX IF NOT EXISTS idx_course_localizations_search_trgm
    ON course_localizations
    USING GIN ((lower(coalesce(title, '') || ' ' || coalesce(short_description, '') || ' ' || coalesce(description, ''))) gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_topic_localizations_name_trgm
    ON topic_localizations
    USING GIN ((lower(name)) gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_tag_localizations_name_trgm
    ON tag_localizations
    USING GIN ((lower(name)) gin_trgm_ops);
