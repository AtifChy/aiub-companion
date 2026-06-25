package settings

import "testing"

func TestValidateSchema(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name: "Valid settings",
			json: `{
				"theme": "dark",
				"log_level": "DEBUG",
				"sync": {
					"interval_minutes": 15,
					"fetch_count": 10,
					"on_startup": true
				},
				"launch": {
					"start_minimized": true,
					"close_to_tray": false,
					"restore_window": false,
					"sidebar_open": true
				},
				"notifications": {
					"enabled": true
				},
				"window": {
					"width": 1024,
					"height": 768,
					"x": -1,
					"y": -1,
					"maximized": false
				}
			}`,
			wantErr: false,
		},
		{
			name: "Invalid theme value",
			json: `{
				"theme": "invalid-theme-value",
				"log_level": "DEBUG",
				"sync": { "interval_minutes": 15, "fetch_count": 10, "on_startup": true },
				"launch": { "start_minimized": true, "close_to_tray": false, "restore_window": false, "sidebar_open": true },
				"notifications": { "enabled": true },
				"window": { "width": 1024, "height": 768, "x": -1, "y": -1, "maximized": false }
			}`,
			wantErr: true,
		},
		{
			name: "Invalid log level",
			json: `{
				"theme": "dark",
				"log_level": "TRACE",
				"sync": { "interval_minutes": 15, "fetch_count": 10, "on_startup": true },
				"launch": { "start_minimized": true, "close_to_tray": false, "restore_window": false, "sidebar_open": true },
				"notifications": { "enabled": true },
				"window": { "width": 1024, "height": 768, "x": -1, "y": -1, "maximized": false }
			}`,
			wantErr: true,
		},
		{
			name: "Malformed JSON syntax",
			json: `{
				"theme": "dark",
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
