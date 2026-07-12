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
