-- FTS5 content table (references notices, no duplicate storage)
CREATE VIRTUAL TABLE IF NOT EXISTS notices_fts USING fts5 (
  title,
  full_title,
  summary,
  content,
  content='notices',
  content_rowid='rowid',
  tokenize='trigram'
);

CREATE VIRTUAL TABLE IF NOT EXISTS offered_courses_fts USING fts5 (
  course_title,
  section,
  faculty,
  class_type,
  day,
  start_time,
  end_time,
  room,
  department,
  content='offered_courses',
  content_rowid='rowid',
  tokenize='trigram'
);
