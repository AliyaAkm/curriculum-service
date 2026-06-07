CREATE OR REPLACE VIEW user_xp_events AS
SELECT
    ucp.id AS event_id,
    ucp.user_id,
    cm.course_id,
    ucp.lesson_id,
    'lesson'::text AS source_type,
    ucp.id::text AS source_id,
    ucp.xp::bigint AS xp,
    ucp.created_at AS activity_at
FROM user_course_points ucp
INNER JOIN course_lessons cl ON cl.id = ucp.lesson_id
INNER JOIN course_modules cm ON cm.id = cl.module_id

UNION ALL

SELECT
    pxa.id AS event_id,
    pxa.user_id,
    pxa.course_id,
    pxa.lesson_id,
    'practice'::text AS source_type,
    pxa.practice_id::text AS source_id,
    pxa.xp::bigint AS xp,
    pxa.created_at AS activity_at
FROM practice_xp_awards pxa;

UPDATE users
SET level = levels.computed_level
FROM (
    SELECT
        u.id,
        (1 + (COALESCE(SUM(uxe.xp), 0)::bigint / 180))::integer AS computed_level
    FROM users u
    LEFT JOIN user_xp_events uxe ON uxe.user_id = u.id
    GROUP BY u.id
) AS levels
WHERE users.id = levels.id
  AND users.level IS DISTINCT FROM levels.computed_level;
