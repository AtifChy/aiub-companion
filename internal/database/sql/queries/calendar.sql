-- name: GetCalendarCache :one
SELECT data, semester, scraped_at
FROM calendar_cache
WHERE calendar_type = ?
LIMIT 1;

-- name: UpsertCalendarCache :exec
INSERT INTO calendar_cache (calendar_type, semester, data)
VALUES (?, ?, ?)
ON CONFLICT(calendar_type) DO UPDATE
SET
  semester = EXCLUDED.semester,
  data = EXCLUDED.data,
  scraped_at = CURRENT_TIMESTAMP;

