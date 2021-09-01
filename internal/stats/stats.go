package stats

import (
	"errors"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
)

// Service - the struct for stats service
type Service struct {
	Storage storage.Store
}

// StatsService - the interface for stats service
type StatsService interface {
	GetStats(ID uint) (storage.Stats, error)
	PostStats(ID uint, stats storage.Stats) error
	GetAllStats() ([]storage.Stats, error)
}

// NewService - returns a new stats service
func NewService(store storage.Store) *Service {
	return &Service{
		Storage: store,
	}
}

// GetStats - get stats by ID from storage
func (s *Service) GetStats(ID uint) (storage.Stats, error) {
	if s.Storage.ExistId(ID) != true {
		return storage.Stats{}, errors.New("can't get stats by this ID")
	}
	stats, err := s.Storage.GetStatsByID(ID)
	if err != nil {
		return storage.Stats{}, err
	}
	return stats, nil
}

// GetAllStats - get all stats from storage
func (s *Service) GetAllStats() (storage.Store, error) {
	return s.Storage, nil
}

// PostStats - save stats to storage by ID
func (s *Service) PostStats(ID uint, stats storage.Stats) error {
	err := s.Storage.SaveStats(ID, stats)
	if err != nil {
		return err
	}
	return nil
}
