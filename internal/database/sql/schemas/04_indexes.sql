CREATE INDEX IF NOT EXISTS idx_notices_category ON notices (category);

CREATE INDEX IF NOT EXISTS idx_notices_urgency ON notices (is_urgent);

CREATE INDEX IF NOT EXISTS idx_notices_date_source ON notices (posted_date DESC, source_order DESC);

CREATE INDEX IF NOT EXISTS idx_attachments_notice_id ON notice_attachments (notice_id);

CREATE INDEX IF NOT EXISTS idx_user_states_pinned ON user_states (is_pinned);

CREATE INDEX IF NOT EXISTS idx_user_states_read ON user_states (is_read);

CREATE INDEX IF NOT EXISTS idx_offered_search ON offered_courses (course_title, faculty);
