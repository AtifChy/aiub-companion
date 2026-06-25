package config

const (
	ID          = "io.github.atifchy.aiub-companion"
	AppName     = "aiub-companion"
	DisplayName = "AIUB Companion"
	Description = "Desktop companion app for AIUB notices and tools"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

type BuildInfo struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
}

func (s *Service) GetBuildInfo() BuildInfo {
	return BuildInfo{
		Name:      DisplayName,
		Version:   Version,
		BuildTime: BuildTime,
	}
}
