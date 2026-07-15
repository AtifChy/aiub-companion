// Package meta contains metadata about the application, such as its name, version, and description.
package meta

const (
	ID          = "io.github.atifchy.aiub-companion"
	AppName     = "aiub-companion"
	DisplayName = "AIUB Companion"
	Description = "Desktop companion app for AIUB notices and tools"
	Repo        = "atifchy/aiub-companion"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func Version() string {
	return version
}
