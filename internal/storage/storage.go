package storage

import (
	"errors"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
)

// Storager - interface for storager
type Storager interface {
	GetStats(ID uint) (*types.Stats, error)
	SaveStats(ID uint, stats types.Stats) error
	GetAllStats() ([]map[uint]types.Stats, error)
}

// Store - struct, where types.Stats saved
type Store struct {
	Storage map[uint]types.Stats
}

// NewStorage - create new storage
func NewStorage() *Store {
	var storage = make(map[uint]types.Stats)
	return &Store{Storage: storage}
}

// GetStatsByID - method of Store struct to get Stats by ID
func (s *Store) GetStatsByID(ID uint) (*types.Stats, error) {
	statsByID, ok := s.Storage[ID]
	if !ok {
		return &types.Stats{}, errors.New("can't get stats by ID")
	}
	return &statsByID, nil
}

// SaveStats - method of Store struct to save stats in storage by ID
func (s *Store) SaveStats(ID uint, stats types.Stats) error {
	s.Storage[ID] = stats
	return nil
}

// GetAllStats - method of Store struct what return all values of Storage
func (s *Store) GetAllStats() (*Store, error) {
	return s, nil
}
