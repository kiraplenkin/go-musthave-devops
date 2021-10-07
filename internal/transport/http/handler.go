package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Handler stores pointers to service
type Handler struct {
	Router  *mux.Router
	Storage *storage.Store
	Cfg     *types.Config
}

// NewHandler returns a pointer to Handler
func NewHandler(s storage.Store, cfg types.Config) *Handler {
	return &Handler{
		Storage: &s,
		Cfg:     &cfg,
	}
}

// SetupRouters sets up all the routes for server
func (h *Handler) SetupRouters() {
	h.Router = mux.NewRouter()
	h.Router.HandleFunc("/", h.GetAllStats).Methods(http.MethodGet)
	h.Router.HandleFunc("/update/", h.PostJSONStat).Methods(http.MethodPost)
	h.Router.HandleFunc("/value/", h.GetStatsByTypeJSON).Methods(http.MethodPost)
	h.Router.HandleFunc("/update/{type}/{id}/{value}", h.PostURLStat).Methods(http.MethodPost)
	h.Router.HandleFunc("/value/{type}/{id}", h.GetStatsByType).Methods(http.MethodGet)
}

//GetAllStats handler that return all values from storage.Store
func (h Handler) GetAllStats(w http.ResponseWriter, _ *http.Request) {
	allStats, err := h.Storage.GetAllStats()
	if err != nil {
		http.Error(w, "Can't get all stats", http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprintf(w, "%v", *allStats)
	if err != nil {
		return
	}
}

// GetStatsByTypeJSON ...
func (h Handler) GetStatsByTypeJSON(w http.ResponseWriter, r *http.Request) {
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

	if statsType != "gauge" && statsType != "counter" {
		http.Error(w, "unknown type", http.StatusNotImplemented)
		return
	}

	var responseStats types.Metrics
	responseStats.ID = id
	responseStats.MType = statsType

	if statsType == "gauge" {
		stat, err := h.Storage.GetGaugeStatsByID(id)
		if err != nil {
			http.Error(w, "can't get gauge stat by this ID", http.StatusNotFound)
			return
		}
		responseStats.Value = &stat.Value
	} else {
		value, err := h.Storage.GetCounterStatsByID(id)
		if err != nil {
			http.Error(w, "can't get counter value by this ID", http.StatusNotFound)
		}
		responseStats.Delta = &value
	}

	resp, err := json.Marshal(responseStats)
	if err != nil {
		http.Error(w, "can't create JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		return
	}
}

// GetStatsByType ...
func (h Handler) GetStatsByType(w http.ResponseWriter, r *http.Request) {
	statsType := mux.Vars(r)["type"]
	id := mux.Vars(r)["id"]

	if statsType != "gauge" && statsType != "counter" {
		http.Error(w, "unknown type", http.StatusNotImplemented)
		return
	}

	if statsType == "gauge" {
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
	} else {
		value, err := h.Storage.GetCounterStatsByID(id)
		if err != nil {
			http.Error(w, "can't get counter value by this ID", http.StatusNotFound)
		}
		_, err = fmt.Fprintf(w, "%+v", value)
		if err != nil {
			return
		}
	}
}

// PostJSONStat ...
func (h Handler) PostJSONStat(w http.ResponseWriter, r *http.Request) {
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
	case "counter":
		statsValue := requestStats.Delta
		newStat := types.Stats{
			Type:  statsType,
			Value: float64(*statsValue),
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

// PostURLStat ...
func (h Handler) PostURLStat(w http.ResponseWriter, r *http.Request) {
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
	//_, err = fmt.Fprintf(w, "%+v", newStat)
	//if err != nil {
	//	return
	//}
}
