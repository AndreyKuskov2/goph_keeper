package service

import "context"

type GophKeeperPingStorager interface {
	Ping(ctx context.Context) error
}

type GophKeeperPingService struct {
	storage GophKeeperPingStorager
}

func NewGophKeeperService(storage GophKeeperPingStorager) *GophKeeperPingService {
	return &GophKeeperPingService{
		storage: storage,
	}
}

func (s *GophKeeperPingService) Ping(ctx context.Context) error {
	return s.storage.Ping(ctx)
}
