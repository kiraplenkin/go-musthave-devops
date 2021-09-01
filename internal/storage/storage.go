package storage

import (
	"errors"
)

// Stats ...
type Stats struct {
	StatsType  string
	StatsValue string
}

// Storager ...
type Storager interface {
	GetStatsByID(uint) Stats
	SaveStats(uint, Stats) error
	ExistId(string) bool
}

// Store ...
type Store struct {
	Storage map[uint]Stats
}

// GetStatsByID ...
func (s *Store) GetStatsByID(ID uint) (Stats, error) {
	statsById, ok := s.Storage[ID]
	if !ok {
		return Stats{}, errors.New("can't get stats by Id")
	}
	return statsById, nil
}

// SaveStats ...
func (s *Store) SaveStats(ID uint, stats Stats) error {
	s.Storage[ID] = stats
	return nil
}

// ExistId ...
func (s *Store) ExistId(ID uint) bool {
	_, ok := s.Storage[ID]
	return ok
}

// New ...
func New() Store {
	s := make(map[uint]Stats)
	var store Store
	store.Storage = s
	return store
}
