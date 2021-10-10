package storage

import (
	"bufio"
	"encoding/json"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"os"
)

// Store struct, where types.Stats saved
type Store struct {
	Storage types.Storage
	writer  *bufio.Writer
}

// NewStorage create new Store
func NewStorage(cfg *types.Config) (*Store, error) {
	statsStorage := &types.Storage{
		GaugeStorage:   map[string]types.Stats{},
		CounterStorage: map[string]int64{},
	}
	if cfg.Restore {
		_, err := os.Stat(cfg.FileStoragePath)
		if !os.IsNotExist(err) {
			readFile, err := os.OpenFile(cfg.FileStoragePath, os.O_RDONLY, 0644)
			if err != nil {
				return nil, err
			}
			_, err = readFile.Stat()
			if err != nil {
				return nil, err
			}
			scanner := bufio.NewScanner(readFile)
			if !scanner.Scan() {
				return nil, scanner.Err()
			}

			data := scanner.Bytes()
			err = json.Unmarshal(data, &statsStorage)
			if err != nil {
				return nil, err
			}

			err = readFile.Close()
			if err != nil {
				return nil, err
			}
		}
	}
	file, err := os.OpenFile(cfg.FileStoragePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Store{
		Storage: *statsStorage,
		writer:  bufio.NewWriter(file),
	}, nil
}

// GetGaugeStatsByID ...
func (s *Store) GetGaugeStatsByID(ID string) (*types.Stats, error) {
	//s.Mu.Lock()
	//defer s.Mu.Unlock()

	statsByID, ok := s.Storage.GaugeStorage[ID]
	if !ok {
		return nil, types.ErrCantGetStats
	}
	return &statsByID, nil
}

// GetCounterStatsByID ...
func (s *Store) GetCounterStatsByID(ID string) (int64, error) {
	//s.Mu.Lock()
	//defer s.Mu.Unlock()

	value, ok := s.Storage.CounterStorage[ID]
	if !ok {
		return 0, types.ErrCantGetStats
	}
	return value, nil
}

// UpdateGaugeStats ...
func (s *Store) UpdateGaugeStats(ID string, stats types.Stats) error {
	//s.Mu.Lock()
	//defer s.Mu.Unlock()

	s.Storage.GaugeStorage[ID] = stats
	return nil
}

// UpdateCounterStats ...
func (s *Store) UpdateCounterStats(ID string, stats types.Stats) error {
	//s.Mu.Lock()
	//s.Mu.Unlock()

	if _, found := s.Storage.CounterStorage[ID]; !found {
		s.Storage.CounterStorage[ID] = int64(stats.Value)
	} else {
		s.Storage.CounterStorage[ID] += int64(stats.Value)
	}
	return nil
}

// GetAllStats return all types.Stats of Storage
func (s *Store) GetAllStats() (*types.Storage, error) {
	//s.Mu.Lock()
	//defer s.Mu.Unlock()

	return &s.Storage, nil
}

// WriteToFile ...
func (s *Store) WriteToFile() error {
	//s.Mu.Lock()
	//defer s.Mu.Unlock()

	data, err := json.Marshal(&s.Storage)
	if err != nil {
		return err
	}
	_, err = s.writer.Write(data)
	if err != nil {
		return err
	}
	err = s.writer.WriteByte('\n')
	if err != nil {
		return err
	}
	return s.writer.Flush()
}
