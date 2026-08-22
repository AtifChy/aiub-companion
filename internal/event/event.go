// Package event provides a set of constants representing different events that can occur within the application.
package event

const (
	EventConfigChanged = "config:changed"

	EventNoticeSyncing = "notice:syncing"
	EventNoticeSynced  = "notice:synced"
	EventNoticeOpen    = "notice:open"

	EventMainWindowShow    = "main-window:show"
	EventMainWindowClosing = "main-window:closing"

	EventWindowMaximized = "window:maximized"
)
