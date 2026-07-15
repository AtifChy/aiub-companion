// Package autostart provides functionality to manage application autostart behavior across different operating systems.
package autostart

type Service struct{}

func NewService() *Service {
	return &Service{}
}
