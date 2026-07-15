package meta

type Service struct{}

func NewService() *Service {
	return &Service{}
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
