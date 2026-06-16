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
					"notice_fetch_count": 10,
					"on_startup": true
				},
				"launch": {
					"start_minimized": true,
					"close_to_tray": false
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
				"theme": "invalid-theme-value"
			}`,
			wantErr: true,
		},
		{
			name: "Invalid log level",
			json: `{
				"log_level": "TRACE"
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
