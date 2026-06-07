WITH eligible AS (
    SELECT DISTINCT ON (cea.user_id, pt.id)
        cea.id AS attempt_id,
        cea.user_id,
        pt.id AS practice_id,
        cm.course_id,
        pt.lesson_id,
        pt.xp_reward
    FROM code_execution_attempts cea
    INNER JOIN practice_tasks pt ON pt.id::text = cea.practice_id
    INNER JOIN course_lessons cl ON cl.id = pt.lesson_id
    INNER JOIN course_modules cm ON cm.id = cl.module_id
    WHERE cea.run_type = 'submit'
      AND cea.passed = true
      AND pt.xp_reward > 0
    ORDER BY cea.user_id, pt.id, cea.created_at ASC
),
inserted AS (
    INSERT INTO practice_xp_awards (
        id,
        user_id,
        practice_id,
        course_id,
        lesson_id,
        xp
    )
    SELECT
        uuid_generate_v4(),
        eligible.user_id,
        eligible.practice_id,
        eligible.course_id,
        eligible.lesson_id,
        eligible.xp_reward
    FROM eligible
    ON CONFLICT (user_id, practice_id) DO NOTHING
    RETURNING user_id, practice_id, xp
)
UPDATE code_execution_attempts cea
SET xp_awarded = inserted.xp
FROM eligible
INNER JOIN inserted
    ON inserted.user_id = eligible.user_id
   AND inserted.practice_id = eligible.practice_id
WHERE cea.id = eligible.attempt_id
  AND cea.xp_awarded = 0;

UPDATE users
SET level = GREATEST(
    COALESCE(level, 1),
    1 + ((
        COALESCE((SELECT SUM(xp) FROM user_course_points WHERE user_id = users.id), 0) +
        COALESCE((SELECT SUM(xp) FROM practice_xp_awards WHERE user_id = users.id), 0)
    )::bigint / 180)::int
)
WHERE id IN (
    SELECT DISTINCT user_id
    FROM practice_xp_awards
);
