-- name: ClearOfferedCourses :exec
DELETE FROM offered_courses;

-- name: InsertOfferedCourse :exec
INSERT INTO
  offered_courses (class_id, course_code, course_title, section, faculty, department)
VALUES
  (?, ?, ?, ?, ?, ?)
ON CONFLICT DO NOTHING;

-- name: InsertClassSchedule :exec
INSERT INTO
  class_schedule (class_id, class_type, day, start_time, end_time, room)
VALUES
  (?, ?, ?, ?, ?, ?);

-- name: ListOfferedCourses :many
SELECT
  o.class_id,
  o.course_code,
  o.course_title,
  o.section,
  o.faculty,
  o.department,
  s.class_type,
  s.day,
  s.start_time,
  s.end_time,
  s.room
FROM
  offered_courses o
  JOIN class_schedule s ON o.class_id = s.class_id;

-- name: AddToUserRoutine :exec
INSERT INTO
  user_routine (class_id)
VALUES
  (?)
ON CONFLICT(class_id) DO NOTHING;

-- name: RemoveFromUserRoutine :exec
DELETE FROM
  user_routine
WHERE
  class_id = ?;

-- name: ListUserRoutine :many
SELECT
  o.class_id,
  o.course_code,
  o.course_title,
  o.section,
  o.faculty,
  o.department,
  s.class_type,
  s.day,
  s.start_time,
  s.end_time,
  s.room
FROM
  user_routine u
  JOIN offered_courses o ON u.class_id = o.class_id
  JOIN class_schedule s ON o.class_id = s.class_id
ORDER BY
  CASE s.day
    WHEN 'Sunday' THEN 1
    WHEN 'Monday' THEN 2
    WHEN 'Tuesday' THEN 3
    WHEN 'Wednesday' THEN 4
    WHEN 'Thursday' THEN 5
    WHEN 'Friday' THEN 6
    WHEN 'Saturday' THEN 7
    ELSE 8
  END,
  s.start_time ASC;
