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
}

// NewHandler returns a pointer to Handler
func NewHandler(s storage.Store) *Handler {
	return &Handler{
		Storage: &s,
	}
}

// SetupRouters sets up all the routes for server
func (h *Handler) SetupRouters() {
	h.Router = mux.NewRouter()
	h.Router.HandleFunc("/{id}", h.GetStatsByID).Methods(http.MethodGet)
	h.Router.HandleFunc("/", h.GetAllStats).Methods(http.MethodGet)
	h.Router.HandleFunc("/update/", h.PostJSONStat).Methods(http.MethodPost)
	h.Router.HandleFunc("/value/", h.GetStatsByTypeJSON).Methods(http.MethodPost)
	h.Router.HandleFunc("/update/{type}/{id}/{value}", h.PostURLStat).Methods(http.MethodPost)
	h.Router.HandleFunc("/value/{type}/{id}", h.GetStatsByType).Methods(http.MethodGet)
}

//GetStatsByID handler that return types.Stats by ID
func (h *Handler) GetStatsByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	stat, err := h.Storage.GetStatsByID(id)
	if err != nil {
		http.Error(w, "Can't get stat by this ID", http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(stat); err != nil {
		http.Error(w, "unable to marshal the struct", http.StatusBadRequest)
		return
	}

	//_, err = fmt.Fprintf(w, "%+v", stat)
	//if err != nil {
	//	http.Error(w, error.Error(err), http.StatusInternalServerError)
	//}
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

	stat, err := h.Storage.GetStatsByID(id)
	if err != nil {
		http.Error(w, "can't get stat by this ID", http.StatusNotFound)
		return
	}

	var responseStats types.Metrics
	responseStats.ID = id
	responseStats.MType = statsType
	if statsType == "gauge" {
		responseStats.Value = &stat.Value
	} else {
		delta := int64(stat.Value)
		responseStats.Delta = &delta
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

// PostJSONStat handler that save json request to storage.Store
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
	if statsType == "gauge" {
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
	} else if statsType == "counter" {
		statsValue := requestStats.Delta
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
		http.Error(w, "unknown type", http.StatusNotImplemented)
		return
	}
	//_, err = fmt.Fprintf(w, "%+v", newStat)
	//if err != nil {
	//	return
	//}
}

// PostURLStat ...
func (h Handler) PostURLStat(w http.ResponseWriter, r *http.Request) {
	statsType := mux.Vars(r)["type"]
	if statsType != "gauge" && statsType != "counter" {
		http.Error(w, "unknown type", http.StatusNotImplemented)
		return
	}
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

	if statsType == "gauge" {
		err = h.Storage.UpdateGaugeStats(id, newStat)
		if err != nil {
			http.Error(w, "can't save stat", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		err = h.Storage.UpdateCounterStats(id, newStat)
		if err != nil {
			http.Error(w, "can't save stat", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	}
	_, err = fmt.Fprintf(w, "%+v", newStat)
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
	stat, err := h.Storage.GetStatsByID(id)
	if err != nil {
		http.Error(w, "can't get stat by this ID", http.StatusNotFound)
		return
	}
	_, err = fmt.Fprintf(w, "%+v", stat.Value)
	if err != nil {
		return
	}
}
