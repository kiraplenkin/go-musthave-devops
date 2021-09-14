package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"os"
)

// ErrCantGetStats custom error if there isn't ID in Store
var ErrCantGetStats = errors.New("can't get stats by ID")

// Store struct, where types.Stats saved
type Store struct {
	Storage types.Storage
	File    *os.File
	writer  *bufio.Writer
}

// NewStorage create new Store
// if cfg.FileStoragePath exists - read stats from it
// if not - create new cfg.FileStoragePath
func NewStorage(cfg *types.ServerConfig) (*Store, error) {
	statsStorage := &types.Storage{}
	_, err := os.Stat(cfg.FileStoragePath)
	if !os.IsNotExist(err) {
		rFile, err := os.OpenFile(cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0777)
		if err != nil {
			return nil, err
		}
		_, err = rFile.Stat()
		if err != nil {
			fmt.Println("Not stats")
		}
		scanner := bufio.NewScanner(rFile)
		if !scanner.Scan() {
			return nil, scanner.Err()
		}

		data := scanner.Bytes()
		err = json.Unmarshal(data, statsStorage)
		if err != nil {
			return nil, err
		}

		func(rFile *os.File) {
			err := rFile.Close()
			if err != nil {
				return
			}
		}(rFile)
		err = os.Remove(cfg.FileStoragePath)
		if err != nil {
			return nil, err
		}
	}

	wFile, err := os.OpenFile(cfg.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &Store{
		Storage: *statsStorage,
		File:    wFile,
		writer:  bufio.NewWriter(wFile),
	}, nil
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
func (s *Store) GetAllStats() (*types.Storage, error) {
	return &s.Storage, nil
}

// SaveToFile save stats to file before shutdown server
func (s *Store) SaveToFile() error {
	data, err := json.Marshal(&s.Storage)
	if err != nil {
		return err
	}
	if _, err = s.writer.Write(data); err != nil {
		return err
	}
	if err = s.writer.WriteByte('\n'); err != nil {
		return err
	}

	return s.writer.Flush()
}
