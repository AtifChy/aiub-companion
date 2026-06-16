package app

type Info struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetInfo() Info {
	return Info{
		Name:    DisplayName,
		Version: Version,
	}
}
