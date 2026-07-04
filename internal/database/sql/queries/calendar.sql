-- name: GetCalendarCache :one
SELECT data, scraped_at
FROM calendar_cache
WHERE calendar_type = ?
LIMIT 1;

-- name: UpsertCalendarCache :exec
INSERT INTO calendar_cache (calendar_type, data)
VALUES (?, ?)
ON CONFLICT(calendar_type) DO UPDATE
SET
  data = EXCLUDED.data,
  scraped_at = CURRENT_TIMESTAMP;

