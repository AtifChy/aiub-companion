package notice

import (
	"fmt"
	"time"
)

type NotificationPayload struct {
	ID    string
	Title string
	Body  string
}

func BuildNotificationPayload(notices []Notice) NotificationPayload {
	count := len(notices)
	if count == 0 {
		return NotificationPayload{}
	}

	if count == 1 {
		n := notices[0]
		return NotificationPayload{
			ID:    n.ID,
			Title: n.Title,
			Body:  n.Summary,
		}
	}

	return NotificationPayload{
		ID:    fmt.Sprintf("sync-%d", time.Now().UnixMilli()),
		Title: fmt.Sprintf("%d new notices available", count),
		Body:  "Click to view the latest notices",
	}
}
