package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

// Handler stores pointers to service
type Handler struct {
	Router  *mux.Router
	Storage *storage.Store
	Cfg     *types.Config
	Mu      *sync.Mutex
}

// NewHandler returns a pointer to Handler
func NewHandler(s *storage.Store, cfg types.Config) *Handler {
	return &Handler{
		Storage: s,
		Cfg:     &cfg,
		Mu:      &sync.Mutex{},
	}
}

// SetupRouters sets up all the routes for server
func (h *Handler) SetupRouters() {
	h.Router = mux.NewRouter()
	h.Router.HandleFunc("/", h.GetAllStats).Methods(http.MethodGet)
	h.Router.HandleFunc("/update/", h.PostJSONStat).Methods(http.MethodPost)
	h.Router.HandleFunc("/updates/", h.PostJSONStats).Methods(http.MethodPost)
	h.Router.HandleFunc("/value/", h.GetStatsByTypeJSON).Methods(http.MethodPost)
	h.Router.HandleFunc("/update/{type}/{id}/{value}", h.PostURLStat).Methods(http.MethodPost)
	h.Router.HandleFunc("/value/{type}/{id}", h.GetStatsByType).Methods(http.MethodGet)
	h.Router.HandleFunc("/ping", h.Ping).Methods(http.MethodGet)
}

//GetAllStats handler that return all values from storage.Store
func (h Handler) GetAllStats(w http.ResponseWriter, _ *http.Request) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	allStats, err := h.Storage.GetAllStats()
	if err != nil {
		http.Error(w, "Can't get all stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/html")
	_, err = fmt.Fprintf(w, "%v", *allStats)
	if err != nil {
		return
	}
}

// GetStatsByTypeJSON ...
func (h Handler) GetStatsByTypeJSON(w http.ResponseWriter, r *http.Request) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	var requestStats types.Metrics
	err = json.Unmarshal(body, &requestStats)
	if err != nil {
		http.Error(w, "can't decode input json", http.StatusBadRequest)
		return
	}

	id := requestStats.ID
	statsType := requestStats.MType

	var responseStats types.Metrics
	responseStats.ID = id
	responseStats.MType = statsType

	switch statsType {
	case "gauge":
		stat, err := h.Storage.GetGaugeStatsByID(id)
		if err != nil {
			http.Error(w, "can't get gauge stat by this ID", http.StatusNotFound)
			return
		}
		responseStats.Value = &stat.Value
		if h.Cfg.Key != "" {
			hash := hmac.New(sha256.New, []byte(h.Cfg.Key))
			hash.Write([]byte(fmt.Sprintf("%s:gauge:%f", id, stat.Value)))
			dst := hash.Sum(nil)
			responseStats.Hash = fmt.Sprintf("%x", dst)
		}
	case "counter":
		value, err := h.Storage.GetCounterStatsByID(id)
		if err != nil {
			http.Error(w, "can't get counter value by this ID", http.StatusNotFound)
			return
		}
		responseStats.Delta = &value
		if h.Cfg.Key != "" {
			hash := hmac.New(sha256.New, []byte(h.Cfg.Key))
			hash.Write([]byte(fmt.Sprintf("%s:counter:%d", id, value)))
			dst := hash.Sum(nil)
			responseStats.Hash = fmt.Sprintf("%x", dst)
		}
	default:
		http.Error(w, "unknown type", http.StatusNotImplemented)
		return
	}

	resp, err := json.Marshal(responseStats)
	if err != nil {
		http.Error(w, "can't create JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		return
	}
}

// GetStatsByType ...
func (h Handler) GetStatsByType(w http.ResponseWriter, r *http.Request) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	statsType := mux.Vars(r)["type"]
	id := mux.Vars(r)["id"]
	switch statsType {
	case "gauge":
		stat, err := h.Storage.GetGaugeStatsByID(id)
		if err != nil {
			http.Error(w, "can't get gauge stat by this ID", http.StatusNotFound)
			return
		}
		w.Header().Set("content-type", "text/html")
		_, err = fmt.Fprintf(w, "%+v", stat.Value)
		if err != nil {
			return
		}
	case "counter":
		value, err := h.Storage.GetCounterStatsByID(id)
		if err != nil {
			http.Error(w, "can't get counter value by this ID", http.StatusNotFound)
		}
		_, err = fmt.Fprintf(w, "%+v", value)
		if err != nil {
			return
		}
	default:
		http.Error(w, "unknown type", http.StatusNotImplemented)
		return
	}
}

// PostJSONStat ...
func (h Handler) PostJSONStat(w http.ResponseWriter, r *http.Request) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	var requestStats types.Metrics
	err = json.Unmarshal(body, &requestStats)
	if err != nil {
		http.Error(w, "can't decode input json", http.StatusBadRequest)
		return
	}

	id := requestStats.ID
	statsType := requestStats.MType

	switch statsType {
	case "gauge":
		statsValue := requestStats.Value
		// todo create func
		if h.Cfg.Key != "" {
			hash := hmac.New(sha256.New, []byte(h.Cfg.Key))
			hash.Write([]byte(fmt.Sprintf("%s:gauge:%f", id, *statsValue)))
			if hmac.Equal([]byte(requestStats.Hash), []byte(fmt.Sprintf("%x", hash.Sum(nil)))) {
				newStat := types.Stats{
					Type:  statsType,
					Value: *statsValue,
				}
				err = h.Storage.UpdateGaugeStats(id, newStat)
				if err != nil {
					http.Error(w, "can't save stat", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
			} else {
				http.Error(w, "hash doesn't equal", http.StatusBadRequest)
			}
		}
		newStat := types.Stats{
			Type:  statsType,
			Value: *statsValue,
		}
		err = h.Storage.UpdateGaugeStats(id, newStat)
		if err != nil {
			http.Error(w, "can't save stat", http.StatusInternalServerError)
			return
		}
		//w.WriteHeader(http.StatusOK)
	case "counter":
		statsValue := requestStats.Delta
		// todo create func
		if h.Cfg.Key != "" {
			hash := hmac.New(sha256.New, []byte(h.Cfg.Key))
			hash.Write([]byte(fmt.Sprintf("%s:counter:%d", id, *statsValue)))
			if hmac.Equal([]byte(requestStats.Hash), []byte(fmt.Sprintf("%x", hash.Sum(nil)))) {
				newStat := types.Stats{
					Type:  statsType,
					Value: float64(*statsValue),
				}
				err = h.Storage.UpdateGaugeStats(id, newStat)
				if err != nil {
					http.Error(w, "can't save stat", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
			} else {
				http.Error(w, "hash doesn't equal", http.StatusBadRequest)
				return
			}
		}
		newStat := types.Stats{
			Type:  statsType,
			Value: float64(*statsValue),
		}
		err = h.Storage.UpdateCounterStats(id, newStat)
		if err != nil {
			http.Error(w, "can't save stat", http.StatusInternalServerError)
			return
		}
		//w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "unknown type", http.StatusNotImplemented)
		return
	}
}

// PostJSONStats ...
func (h Handler) PostJSONStats(w http.ResponseWriter, r *http.Request) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	var requestStats []types.Metrics
	err = json.Unmarshal(body, &requestStats)
	if err != nil {
		http.Error(w, "can't decode input json", http.StatusBadRequest)
		return
	}

	for _, requestStat := range requestStats {
		id := requestStat.ID
		statsType := requestStat.MType

		switch statsType {
		case "gauge":
			statsValue := requestStat.Value
			// todo create func
			if h.Cfg.Key != "" {
				hash := hmac.New(sha256.New, []byte(h.Cfg.Key))
				hash.Write([]byte(fmt.Sprintf("%s:gauge:%f", id, *statsValue)))
				if hmac.Equal([]byte(requestStat.Hash), []byte(fmt.Sprintf("%x", hash.Sum(nil)))) {
					newStat := types.Stats{
						Type:  statsType,
						Value: *statsValue,
					}
					err = h.Storage.UpdateGaugeStats(id, newStat)
					if err != nil {
						http.Error(w, "can't save stat", http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusOK)
				} else {
					http.Error(w, "hash doesn't equal", http.StatusBadRequest)
				}
			}
			newStat := types.Stats{
				Type:  statsType,
				Value: *statsValue,
			}
			err = h.Storage.UpdateGaugeStats(id, newStat)
			if err != nil {
				http.Error(w, "can't save stat", http.StatusInternalServerError)
				return
			}
			//w.WriteHeader(http.StatusOK)
		case "counter":
			statsValue := requestStat.Delta
			// todo create func
			if h.Cfg.Key != "" {
				hash := hmac.New(sha256.New, []byte(h.Cfg.Key))
				hash.Write([]byte(fmt.Sprintf("%s:counter:%d", id, *statsValue)))
				if hmac.Equal([]byte(requestStat.Hash), []byte(fmt.Sprintf("%x", hash.Sum(nil)))) {
					newStat := types.Stats{
						Type:  statsType,
						Value: float64(*statsValue),
					}
					err = h.Storage.UpdateGaugeStats(id, newStat)
					if err != nil {
						http.Error(w, "can't save stat", http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusOK)
				} else {
					http.Error(w, "hash doesn't equal", http.StatusBadRequest)
					return
				}
			}
			newStat := types.Stats{
				Type:  statsType,
				Value: float64(*statsValue),
			}
			err = h.Storage.UpdateCounterStats(id, newStat)
			if err != nil {
				http.Error(w, "can't save stat", http.StatusInternalServerError)
				return
			}
			//w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "unknown type", http.StatusNotImplemented)
			return
		}
	}
}

// PostURLStat ...
func (h Handler) PostURLStat(w http.ResponseWriter, r *http.Request) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	statsType := mux.Vars(r)["type"]
	switch statsType {
	case "gauge":
		statsValue, err := strconv.ParseFloat(mux.Vars(r)["value"], 64)
		if err != nil {
			http.Error(w, error.Error(err), http.StatusBadRequest)
			return
		}
		id := mux.Vars(r)["id"]
		newStat := types.Stats{
			Type:  statsType,
			Value: statsValue,
		}
		err = h.Storage.UpdateGaugeStats(id, newStat)
		if err != nil {
			http.Error(w, "can't save stat", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	case "counter":
		statsValue, err := strconv.ParseInt(mux.Vars(r)["value"], 10, 64)
		if err != nil {
			http.Error(w, error.Error(err), http.StatusBadRequest)
			return
		}
		id := mux.Vars(r)["id"]
		newStat := types.Stats{
			Type:  statsType,
			Value: float64(statsValue),
		}
		err = h.Storage.UpdateCounterStats(id, newStat)
		if err != nil {
			http.Error(w, "can't save stat", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "unknown type", http.StatusNotImplemented)
		return
	}
}

// Ping db
func (h Handler) Ping(w http.ResponseWriter, _ *http.Request) {
	err := h.Storage.Ping()
	if err != nil {
		http.Error(w, fmt.Sprintf("can't connect to db: %v", err), http.StatusInternalServerError)
	}
}
