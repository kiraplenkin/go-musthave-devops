package storage

import (
	"bufio"
	"database/sql"
	"encoding/json"
	_ "github.com/jackc/pgx/v4"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"github.com/pressly/goose/v3"
	"os"
)

// Store struct, where types.Storage saved
type Store struct {
	Storage types.Storage
	writer  *bufio.Writer
	db      *sql.DB
	cfg     types.Config
}

// NewStorage create new Store with types.Storage and writer
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

	db, err := sql.Open("postgres", cfg.Database)
	if err != nil {
		return nil, err
	}

	err = goose.Up(db, "migrations")
	if err != nil {
		panic(err)
	}
	return &Store{
		Storage: *statsStorage,
		writer:  bufio.NewWriter(file),
		db:      db,
		cfg:     *cfg,
	}, nil

}

// GetGaugeStatsByID return gauge metric from GaugeStorage by ID
func (s *Store) GetGaugeStatsByID(ID string) (*types.Stats, error) {
	statsByID, ok := s.Storage.GaugeStorage[ID]
	if !ok {
		return nil, types.ErrCantGetStats
	}
	return &statsByID, nil
}

// GetCounterStatsByID return counter metric from CounterStorage by ID
func (s *Store) GetCounterStatsByID(ID string) (int64, error) {
	value, ok := s.Storage.CounterStorage[ID]
	if !ok {
		return 0, types.ErrCantGetStats
	}
	return value, nil
}

// UpdateGaugeStats replace GaugeStorage by ID
func (s *Store) UpdateGaugeStats(ID string, stats types.Stats) error {
	s.Storage.GaugeStorage[ID] = stats
	return nil
}

// UpdateCounterStats increase CounterStorage by ID if exist or create new
func (s *Store) UpdateCounterStats(ID string, stats types.Stats) error {
	if _, found := s.Storage.CounterStorage[ID]; !found {
		s.Storage.CounterStorage[ID] = int64(stats.Value)
	} else {
		s.Storage.CounterStorage[ID] += int64(stats.Value)
	}
	return nil
}

// GetAllStats return all types.Stats from types.Storage
func (s *Store) GetAllStats() (*types.Storage, error) {
	return &s.Storage, nil
}

// Save types.Storage to file or db
func (s *Store) Save() error {
	if s.cfg.Database == "" {
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
	} else {
		for key, metric := range s.Storage.GaugeStorage {
			_, err := s.db.Exec("INSERT INTO metrics (id, mtype, value) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET mtype=$2, value=$3;", key, metric.Type, metric.Value)
			if err != nil {
				return err
			}
		}
		for key, delta := range s.Storage.CounterStorage {
			_, err := s.db.Exec("INSERT INTO metrics (id, mtype, delta) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET mtype=$2, delta=$3;", key, "counter", delta)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Store) Ping() error {
	err := s.db.Ping()
	if err != nil {
		return err
	}
	return nil
}
