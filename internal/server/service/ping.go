package service

import "context"

// GophKeeperPingStorager is the interface that wraps the basic methods for working with ping.
type GophKeeperPingStorager interface {
	Ping(ctx context.Context) error
}

// GophKeeperPingService is the service that wraps the basic methods for working with ping.
type GophKeeperPingService struct {
	storage GophKeeperPingStorager
}

// NewGophKeeperService creates a new GophKeeperPingService.
func NewGophKeeperService(storage GophKeeperPingStorager) *GophKeeperPingService {
	return &GophKeeperPingService{
		storage: storage,
	}
}

// Ping checks if the service is available.
func (s *GophKeeperPingService) Ping(ctx context.Context) error {
	return s.storage.Ping(ctx)
}
