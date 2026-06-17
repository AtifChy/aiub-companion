// Package notice provides notice domain for the application.
package notice

// Notice represents a notice with its details and metadata.
type Notice struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	FullTitle   string       `json:"fullTitle"`
	Summary     string       `json:"summary"`
	Content     string       `json:"content"`
	PostedDate  string       `json:"postedDate"`
	Category    string       `json:"category"`
	Attachments []Attachment `json:"attachments"`
	IsCached    bool         `json:"isCached"`
	IsUrgent    bool         `json:"isUrgent"`
	IsPinned    bool         `json:"isPinned"`
	IsRead      bool         `json:"isRead"`
}

// NoticeDetails contains the only the details of a notice.
type NoticeDetails struct {
	FullTitle   string       `json:"title"`
	Content     string       `json:"content"`
	Attachments []Attachment `json:"attachments"`
}

// Attachment represents an attachment associated with a notice.
type Attachment struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	URL       string `json:"url"`
	LocalPath string `json:"localPath"`
}

// Filter represents query criteria for retrieving notices.
type Filter struct {
	Urgent   *bool
	Pinned   *bool
	Unread   *bool
	Category string
	Search   string
	Limit    int
	Offset   int
}
