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
}

// NewStorage create new Store
func NewStorage(cfg *types.ServerConfig) (*Store, error) {
	statsStorage := &types.Storage{}
	statsStorage.GaugeStorage = map[string]types.Stats{}
	statsStorage.CounterStorage = map[string]int64{}
	_, err := os.Stat(cfg.FileStoragePath)
	if !os.IsNotExist(err) {
		err := ReadFromFile(statsStorage, cfg.FileStoragePath)
		if err != nil {
			return nil, err
		}
	}

	return &Store{
		Storage: *statsStorage,
	}, nil
}

// GetGaugeStatsByID ...
func (s *Store) GetGaugeStatsByID(ID string) (*types.Stats, error) {
	statsByID, ok := s.Storage.GaugeStorage[ID]
	if !ok {
		return nil, types.ErrCantGetStats
	}
	return &statsByID, nil
}

// GetCounterStatsByID ...
func (s *Store) GetCounterStatsByID(ID string) (int64, error) {
	value, ok := s.Storage.CounterStorage[ID]
	if !ok {
		return 0, types.ErrCantGetStats
	}
	return value, nil
}

// UpdateGaugeStats ...
func (s *Store) UpdateGaugeStats(ID string, stats types.Stats) error {
	s.Storage.GaugeStorage[ID] = stats
	return nil
}

// UpdateCounterStats ...
func (s *Store) UpdateCounterStats(ID string, stats types.Stats) error {
	if _, found := s.Storage.CounterStorage[ID]; !found {
		s.Storage.CounterStorage[ID] = int64(stats.Value)
	} else {
		s.Storage.CounterStorage[ID] += int64(stats.Value)
	}
	return nil
}

// GetAllStats return all types.Stats of Storage
func (s *Store) GetAllStats() (*types.Storage, error) {
	return &s.Storage, nil
}

// SaveToFile func to save types.Storage to file before shutdown server
func SaveToFile(data []byte, fileName string) error {
	writeFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer func(writeFile *os.File) {
		err := writeFile.Close()
		if err != nil {
			return
		}
	}(writeFile)

	writer := bufio.NewWriter(writeFile)

	_, err = writer.Write(data)
	if err != nil {
		return err
	}

	err = writer.WriteByte('\n')
	if err != nil {
		return err
	}
	return writer.Flush()
}

// ReadFromFile func to read stats if exist types.ServerConfig to types.Storage
func ReadFromFile(storage *types.Storage, fileName string) error {
	readFile, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	_, err = readFile.Stat()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(readFile)
	if !scanner.Scan() {
		return scanner.Err()
	}

	data := scanner.Bytes()
	err = json.Unmarshal(data, storage)
	if err != nil {
		return err
	}

	err = readFile.Close()
	if err != nil {
		return err
	}
	err = os.Remove(fileName)
	if err != nil {
		return err
	}
	return nil
}
