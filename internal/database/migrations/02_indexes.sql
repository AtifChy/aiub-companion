-- =============================================================================
-- notices
-- =============================================================================

-- ListNoticesWithState: optional filter  WHERE n.category = ?
CREATE INDEX IF NOT EXISTS idx_notices_category
  ON notices (category);

-- ListNoticesWithState: optional filter  WHERE n.is_urgent = ?
CREATE INDEX IF NOT EXISTS idx_notices_urgency
  ON notices (is_urgent);

-- ListNoticesWithState ORDER BY + GetLatestNoticeInfo MAX() aggregates.
-- Composite so SQLite can satisfy the sort without a separate sort step.
CREATE INDEX IF NOT EXISTS idx_notices_date_source
  ON notices (posted_date DESC, source_order DESC);

-- =============================================================================
-- notice_attachments
-- =============================================================================

-- GetNoticeAttachments: WHERE notice_id = ?
-- (notice_id is not the PK here, so a dedicated index is required)
CREATE INDEX IF NOT EXISTS idx_attachments_notice_id
  ON notice_attachments (notice_id);

-- =============================================================================
-- class_schedule
-- =============================================================================

-- ListOfferedCourses: JOIN class_schedule s ON o.class_id = s.class_id
-- ListUserRoutine:    JOIN class_schedule s ON o.class_id = s.class_id
-- class_id is a non-PK FK column; without this index both queries do a full
-- table scan of class_schedule for every course in offered_courses.
CREATE INDEX IF NOT EXISTS idx_class_schedule_class_id
  ON class_schedule (class_id);
