package storage

import (
	"errors"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
)

// ErrCantGetStats custom error if there isn't ID in Store
var ErrCantGetStats = errors.New("can't get stats by ID")

// Store struct, where types.Stats saved
type Store struct {
	Storage types.Storage
}

// NewStorage create new Store
func NewStorage() *Store {
	return &Store{Storage: types.Storage{}}
}

// GetStatsByID get types.Stats by ID
func (s *Store) GetStatsByID(ID uint) (*types.Stats, error) {
	statsByID, ok := s.Storage[ID]
	if !ok {
		return nil, ErrCantGetStats
	}
	return &statsByID, nil
}

// SaveStats save types.Stats in Storage by ID
func (s *Store) SaveStats(ID uint, stats types.Stats) error {
	s.Storage[ID] = stats
	return nil
}

// GetAllStats return all types.Stats of Storage
func (s *Store) GetAllStats() (*Store, error) {
	return s, nil
}
