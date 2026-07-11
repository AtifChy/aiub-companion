package config

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

type BuildInfo struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
}

func (s *Service) GetBuildInfo() BuildInfo {
	return BuildInfo{
		Name:      DisplayName,
		Version:   version,
		BuildTime: buildTime,
	}
}
