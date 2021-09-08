package storage

import (
	"errors"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
)

// Storager - interface for storager
type Storager interface {
	GetStats(ID uint) (*types.Stats, error)
	SaveStats(ID uint, stats *types.Stats) error
	GetAllStats() ([]map[uint]types.Stats, error)
	ReadFromFile() (*map[uint]types.Stats, error)
	WriteToFile(*map[uint]types.Stats) error
	Close() error
}

// Store - struct, where types.Storage saved
type Store struct {
	Storage types.Storage
	//File    *os.File
	//Scanner *bufio.Scanner
}

// NewStorage - create new storage
func NewStorage() *Store {
	return &Store{Storage: types.Storage{}}
}
//func NewStorage() (*Store, error) {
//	file, err := os.OpenFile("_test.json", os.O_RDONLY|os.O_CREATE, 0777)
//	if err != nil {
//		return nil, err
//	}
//
//	newStorage := Store{
//		File: file,
//		Scanner: bufio.NewScanner(file),
//	}
//
//	data, err := newStorage.ReadFromFile()
//	defer func(newStorage *Store) {
//		err := newStorage.Close()
//		if err != nil {
//			return
//		}
//	}(&newStorage)
//	if err != nil {
//		return nil, err
//	}
//
//	return &Store{Storage: *data, File: file, Scanner: bufio.NewScanner(file)}, nil
//}

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

// Close - func to close File
//func (s *Store) Close() error {
//	return s.File.Close()
//}

// ReadFromFile - read data from file with Scanner
//func (s *Store) ReadFromFile() (*types.Storage, error) {
//	if !s.Scanner.Scan() {
//		return nil, s.Scanner.Err()
//	}
//	data := s.Scanner.Bytes()
//	storage := &types.Storage{}
//	err := json.Unmarshal(data, storage)
//	if err != nil {
//		return nil, err
//	}
//	return storage, nil
//}
