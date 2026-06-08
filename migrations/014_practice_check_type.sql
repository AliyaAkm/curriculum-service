ALTER TABLE practice_tasks
ADD COLUMN IF NOT EXISTS check_type text NOT NULL DEFAULT 'auto';

UPDATE practice_tasks
SET check_type = 'auto'
WHERE check_type IS NULL OR trim(check_type) = '';

ALTER TABLE practice_tasks
DROP CONSTRAINT IF EXISTS practice_tasks_check_type_check;

ALTER TABLE practice_tasks
ADD CONSTRAINT practice_tasks_check_type_check
CHECK (check_type IN ('auto', 'manual'));

CREATE INDEX IF NOT EXISTS idx_practice_tasks_check_type
ON practice_tasks(check_type);
