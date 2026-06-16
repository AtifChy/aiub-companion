-- name: GetNoticeWithState :one
SELECT
  n.*,
  COALESCE(s.is_pinned, 0) AS is_pinned,
  COALESCE(s.is_read, 0) AS is_read
FROM
  notices n
  LEFT JOIN user_states s ON n.id = s.notice_id
WHERE
  n.id = ?
LIMIT
  1;

-- name: SearchNoticesWithState :many
SELECT
  n.id, n.title, n.summary, n.full_title, n.content, n.posted_date, n.category, n.is_cached, n.is_urgent, n.source_order, n.created_at, n.updated_at,
  COALESCE(s.is_pinned, 0) AS is_pinned,
  COALESCE(s.is_read, 0) AS is_read
FROM
  notices n
  LEFT JOIN user_states s ON n.id = s.notice_id
WHERE
  (
    CAST(sqlc.narg(search) AS TEXT) IS NULL
    OR n.rowid IN (
      SELECT rowid FROM notices_fts WHERE notices_fts = sqlc.narg(search)
    )
  )
  AND (
    CAST(sqlc.narg(category) AS TEXT) IS NULL
    OR n.category = sqlc.narg(category)
  )
  AND (
    CAST(sqlc.narg(urgent) AS BOOL) IS NULL
    OR n.is_urgent = sqlc.narg(urgent)
  )
  AND (
    CAST(sqlc.narg(pinned) AS BOOL) IS NULL
    OR COALESCE(s.is_pinned, 0) = sqlc.narg(pinned)
  )
  AND (
    CAST(sqlc.narg(unread) AS BOOL) IS NULL
    OR COALESCE(s.is_read, 0) != sqlc.narg(unread)
  )
ORDER BY
  n.posted_date DESC,
  n.source_order DESC
LIMIT
  @limit
OFFSET
  @offset;

-- name: GetLatestNoticeInfo :one
SELECT
  CAST(COALESCE(MAX(posted_date), '') AS TEXT) AS posted_date,
  CAST(COALESCE(MAX(source_order), 0) AS INTEGER) AS source_order
FROM
  notices;

-- name: CountNotices :one
SELECT
  COUNT(1)
FROM
  notices;

-- name: UpsertNotice :exec
INSERT INTO
  notices (id, title, summary, posted_date, category, is_urgent, source_order)
VALUES
  (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE
  SET
    title = EXCLUDED.title,
    summary = EXCLUDED.summary,
    posted_date = EXCLUDED.posted_date,
    category = EXCLUDED.category,
    is_urgent = EXCLUDED.is_urgent,
    source_order = EXCLUDED.source_order
  WHERE
    COALESCE(notices.title, '') != COALESCE(EXCLUDED.title, '');

-- name: InsertNotice :execrows
INSERT INTO
  notices (id, title, summary, posted_date, category, is_urgent, source_order)
VALUES
  (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO NOTHING;

-- name: UpdateNotice :exec
UPDATE notices
SET
  title = @title,
  summary = @summary,
  posted_date = @posted_date,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = @id
  AND (
    title != @title
    OR summary != @summary
    OR posted_date != @posted_date
  );

-- name: UpdateNoticeDetails :exec
UPDATE notices
SET
  full_title = ?,
  content = ?,
  is_cached = 1,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ?;

-- name: UpsertNoticeAttachment :exec
INSERT INTO
  notice_attachments (id, notice_id, label, url)
VALUES
  (?, ?, ?, ?)
ON CONFLICT(id) DO NOTHING;

-- name: GetNoticeAttachments :many
SELECT
  *
FROM
  notice_attachments
WHERE
  notice_id = ?;

-- name: SetPinState :exec
INSERT INTO
  user_states (notice_id, is_pinned, pinned_at)
VALUES
  (
    ?,
    ?,
    CASE
      WHEN @is_pinned = 1 THEN CURRENT_TIMESTAMP
      ELSE NULL
    END
  )
ON CONFLICT(notice_id) DO UPDATE
SET
  is_pinned = EXCLUDED.is_pinned,
  pinned_at = EXCLUDED.pinned_at;

-- name: SetReadState :exec
INSERT INTO
  user_states (notice_id, is_read, last_read_at)
VALUES
  (
    ?,
    ?,
    CASE
      WHEN @is_read = 1 THEN CURRENT_TIMESTAMP
      ELSE NULL
    END)
ON CONFLICT(notice_id) DO UPDATE
SET
  is_read = EXCLUDED.is_read,
  last_read_at = EXCLUDED.last_read_at;

-- name: ClearNotices :exec
DELETE FROM notices;
