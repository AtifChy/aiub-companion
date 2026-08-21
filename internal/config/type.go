package config

import (
	"encoding/json/v2"
	"fmt"
	"log/slog"
)

type Color string

const (
	ColorLight  Color = "light"
	ColorDark   Color = "dark"
	ColorSystem Color = "system"
)

func (c *Color) UnmarshalText(text []byte) error {
	color := Color(text)
	switch color {
	case ColorLight, ColorDark, ColorSystem:
		*c = color
		return nil
	default:
		return fmt.Errorf("invalid color value: %s", text)
	}
}

type Theme string

type appearance struct {
	Color Color `json:"color" jsonschema:"enum=light,enum=dark,enum=system"`
	Theme Theme `json:"theme"`
}

type notification struct {
	Enabled bool `json:"enabled"`
}

type launch struct {
	AutoStart      bool `json:"auto_start"`
	StartMinimized bool `json:"start_minimized"`
	CloseToTray    bool `json:"close_to_tray"`
	KeepAlive      bool `json:"keep_alive"`
	RestoreWindow  bool `json:"restore_window"`
	SidebarOpen    bool `json:"sidebar_open"`
}

type positiveInt int

func (p *positiveInt) UnmarshalJSON(data []byte) error {
	var value int
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	if value <= 0 {
		return fmt.Errorf("value must be a positive integer, got %d", value)
	}
	*p = positiveInt(value)
	return nil
}

type syn_ struct {
	Interval   positiveInt `json:"interval"`
	FetchCount positiveInt `json:"fetch_count"`
	OnStartup  bool        `json:"on_startup"`
}

type UpdateInterval string

const (
	UpdateIntervalNever   UpdateInterval = "never"
	UpdateIntervalDaily   UpdateInterval = "daily"
	UpdateIntervalWeekly  UpdateInterval = "weekly"
	UpdateIntervalMonthly UpdateInterval = "monthly"
)

func (u *UpdateInterval) UnmarshalText(text []byte) error {
	interval := UpdateInterval(text)
	switch interval {
	case UpdateIntervalNever, UpdateIntervalDaily, UpdateIntervalWeekly, UpdateIntervalMonthly:
		*u = interval
		return nil
	default:
		return fmt.Errorf("invalid update interval value: %s", text)
	}
}

type updates struct {
	Interval UpdateInterval `json:"interval" jsonschema:"enum=never,enum=daily,enum=weekly,enum=monthly"`
}

type logging struct {
	Level slog.Level `json:"level" jsonschema:"enum=DEBUG,enum=INFO,enum=WARN,enum=ERROR"`
}
