CREATE TABLE IF NOT EXISTS notices (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  summary TEXT,
  full_title TEXT,
  content TEXT,
  posted_date TEXT NOT NULL,
  category TEXT NOT NULL DEFAULT 'general',
  is_cached INTEGER NOT NULL DEFAULT 0,
  is_urgent INTEGER NOT NULL DEFAULT 0,
  source_order INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS notice_attachments (
  id TEXT PRIMARY KEY,
  notice_id TEXT NOT NULL,
  label TEXT NOT NULL,
  url TEXT NOT NULL,
  local_path TEXT,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (notice_id) REFERENCES notices (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_states (
  notice_id TEXT PRIMARY KEY,
  is_pinned INTEGER NOT NULL DEFAULT 0,
  is_read INTEGER NOT NULL DEFAULT 0,
  pinned_at TEXT,
  last_read_at TEXT,
  FOREIGN KEY (notice_id) REFERENCES notices (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS calendar_cache (
  calendar_type TEXT PRIMARY KEY,
  data TEXT NOT NULL,
  semester TEXT NOT NULL,
  scraped_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS offered_courses (
  class_id TEXT NOT NULL,
  course_code TEXT,
  course_title TEXT NOT NULL,
  section TEXT NOT NULL,
  faculty TEXT,
  class_type TEXT,
  day TEXT,
  start_time TEXT,
  end_time TEXT,
  room TEXT,
  department TEXT,
  PRIMARY KEY (class_id, day, start_time)
);

CREATE TABLE IF NOT EXISTS user_routine (
  class_id TEXT NOT NULL,
  added_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
