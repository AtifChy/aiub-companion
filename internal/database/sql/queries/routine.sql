-- name: ClearOfferedCourses :exec
DELETE FROM offered_courses;

-- name: InsertOfferedCourse :exec
INSERT INTO
  offered_courses (class_id, course_code, course_title, section, faculty, class_type, day, start_time, end_time, room, department)
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT DO NOTHING;

-- name: ListOfferedCourses :many
SELECT
  *
FROM
  offered_courses;

-- name: AddToUserRoutine :exec
INSERT INTO
  user_routine (class_id)
VALUES
  (?)
ON CONFLICT DO NOTHING;

-- name: RemoveFromUserRoutine :exec
DELETE FROM
  user_routine
WHERE
  class_id = ?;

-- name: ListUserRoutine :many
SELECT
  o.*
FROM
  user_routine u
  JOIN offered_courses o ON u.class_id = o.class_id
ORDER BY
  o.day ASC, o.start_time ASC;
