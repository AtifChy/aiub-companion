// Package meta defines the application metadata and version information.
package meta

const (
	ID          = "io.github.atifchy.aiub-companion"
	Name        = "aiub-companion"
	DisplayName = "AIUB Companion"
	Description = "Desktop companion app for AIUB notices and tools"
)

var Version string

func SetVersion(v string) {
	Version = v
}
