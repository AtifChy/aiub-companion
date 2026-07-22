package config

import "testing"

func TestValidateSchema(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name: "Valid config",
			json: `{
				"appearance": { "theme": "dark" },
				"updates": { "interval": "weekly" },
				"logging": { "level": "DEBUG" },
				"sync": {
					"interval": 15,
					"fetch_count": 10,
					"on_startup": true
				},
				"launch": {
					"auto_start": false,
					"start_minimized": true,
					"close_to_tray": false,
					"keep_alive": false,
					"restore_window": false,
					"sidebar_open": true
				},
				"notifications": {
					"enabled": true
				}
			}`,
			wantErr: false,
		},
		{
			name: "Invalid theme value",
			json: `{
				"appearance": { "theme": "invalid-theme-value" },
				"updates": { "interval": "weekly" },
				"logging": { "level": "DEBUG" },
				"sync": { "interval": 15, "fetch_count": 10, "on_startup": true },
				"launch": { "start_minimized": true, "close_to_tray": false, "keep_alive": false, "restore_window": false, "sidebar_open": true },
				"notifications": { "enabled": true }
			}`,
			wantErr: true,
		},
		{
			name: "Invalid log level",
			json: `{
				"appearance": { "theme": "dark" },
				"updates": { "interval": "weekly" },
				"logging": { "level": "TRACE" },
				"sync": { "interval": 15, "fetch_count": 10, "on_startup": true },
				"launch": { "start_minimized": true, "close_to_tray": false, "keep_alive": false, "restore_window": false, "sidebar_open": true },
				"notifications": { "enabled": true }
			}`,
			wantErr: true,
		},
		{
			name: "Malformed JSON syntax",
			json: `{
				"appearance": { "theme": "dark" },
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate([]byte(tt.json))
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
