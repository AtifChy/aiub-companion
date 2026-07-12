-- Triggers for offered_courses
CREATE TRIGGER IF NOT EXISTS offered_courses_ai AFTER INSERT ON offered_courses BEGIN
  INSERT INTO offered_courses_fts (rowid, course_title, section, faculty, class_type, day, start_time, end_time, room, department)
  VALUES (new.rowid, new.course_title, new.section, new.faculty, new.class_type, new.day, new.start_time, new.end_time, new.room, new.department);
END;

CREATE TRIGGER IF NOT EXISTS offered_courses_ad AFTER DELETE ON offered_courses BEGIN
  INSERT INTO offered_courses_fts (offered_courses_fts, rowid, course_title, section, faculty, class_type, day, start_time, end_time, room, department)
  VALUES ('delete', old.rowid, old.course_title, old.section, old.faculty, old.class_type, old.day, old.start_time, old.end_time, old.room, old.department);
END;

CREATE TRIGGER IF NOT EXISTS offered_courses_au AFTER UPDATE ON offered_courses BEGIN
  INSERT INTO offered_courses_fts (offered_courses_fts, rowid, course_title, section, faculty, class_type, day, start_time, end_time, room, department)
  VALUES ('delete', old.rowid, old.course_title, old.section, old.faculty, old.class_type, old.day, old.start_time, old.end_time, old.room, old.department);
  INSERT INTO offered_courses_fts (rowid, course_title, section, faculty, class_type, day, start_time, end_time, room, department)
  VALUES (new.rowid, new.course_title, new.section, new.faculty, new.class_type, new.day, new.start_time, new.end_time, new.room, new.department);
END;

-- Rebuild FTS index from offered_courses table on every startup
INSERT INTO offered_courses_fts (offered_courses_fts) VALUES ('rebuild');
