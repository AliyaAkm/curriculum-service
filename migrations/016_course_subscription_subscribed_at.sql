ALTER TABLE course_subscription
ADD COLUMN IF NOT EXISTS subscribed_at timestamptz;

UPDATE course_subscription
SET subscribed_at = COALESCE(started_at, last_activity_at, NOW())
WHERE subscribed_at IS NULL;

ALTER TABLE course_subscription
ALTER COLUMN subscribed_at SET DEFAULT NOW();

ALTER TABLE course_subscription
ALTER COLUMN subscribed_at SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_course_subscription_course_subscribed_at
ON course_subscription(course_id, subscribed_at DESC);
