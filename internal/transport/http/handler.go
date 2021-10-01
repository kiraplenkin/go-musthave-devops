package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
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
	//h.Router.HandleFunc("/update", h.PostJsonStat).Methods(http.MethodPost)
	h.Router.HandleFunc("/update/{type}/{id}/{value}", h.PostURLStat).Methods(http.MethodPost)
	h.Router.HandleFunc("/value/{type}/{id}", h.GetStatsByType).Methods(http.MethodGet)

	h.Router.HandleFunc("/health/", h.CheckHealth).Methods(http.MethodGet)
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

////PostJsonStat handler that save json request to storage.Store
//func (h Handler) PostJsonStat(w http.ResponseWriter, r *http.Request) {
//	body, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		http.Error(w, "can't read body", http.StatusBadRequest)
//		return
//	}
//	//var requestStats types.RequestStats
//	var requestStats []types.Metric
//
//	err = json.Unmarshal(body, &requestStats)
//	if err != nil {
//		fmt.Println(err)
//		http.Error(w, "can't decode input json", http.StatusBadRequest)
//		return
//	}
//
//	id := requestStats[0].ID
//	//err = validator.Require(
//	//	requestStats.ID,
//	//	requestStats.TotalAlloc,
//	//	requestStats.Sys,
//	//	requestStats.Mallocs,
//	//	requestStats.Frees,
//	//	requestStats.LiveObjects,
//	//	requestStats.NumGoroutine,
//	//)
//	//if err != nil {
//	//	http.Error(w, "can't save stat", http.StatusBadRequest)
//	//	return
//	//}
//
//	//newStat := types.Stats{
//	//	TotalAlloc:   requestStats.TotalAlloc,
//	//	Sys:          requestStats.Sys,
//	//	Mallocs:      requestStats.Mallocs,
//	//	Frees:        requestStats.Frees,
//	//	LiveObjects:  requestStats.LiveObjects,
//	//	NumGoroutine: requestStats.NumGoroutine,
//	//}
//
//	err = h.Storage.SaveStats(id, requestStats)
//	if err != nil {
//		http.Error(w, "can't save stat", http.StatusInternalServerError)
//		return
//	}
//	w.WriteHeader(http.StatusCreated)
//	_, err = fmt.Fprintf(w, "%+v", requestStats)
//	if err != nil {
//		return
//	}
//}

func existMetric(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// PostURLStat ...
func (h Handler) PostURLStat(w http.ResponseWriter, r *http.Request) {
	statsType := mux.Vars(r)["type"]
	if statsType != "gauge" && statsType != "counter" {
		http.Error(w, "Can't save stat", http.StatusNotImplemented)
		return
	}
	statsValue, err := strconv.ParseFloat(mux.Vars(r)["value"], 64)
	if err != nil {
		http.Error(w, error.Error(err), http.StatusBadRequest)
		return
	}
	id := mux.Vars(r)["id"]
	if !existMetric(id, types.Metrics) {
		http.Error(w, "unknown metric", http.StatusOK)
		return
	}

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
		w.WriteHeader(http.StatusCreated)
		_, err = fmt.Fprintf(w, "%+v", newStat)
		if err != nil {
			return
		}
	} else {
		err = h.Storage.UpdateCounterStats(id, newStat)
		if err != nil {
			http.Error(w, "can't save stat", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, err = fmt.Fprintf(w, "%+v", newStat)
		if err != nil {
			return
		}
	}
}

// CheckHealth handler to check health
func (h Handler) CheckHealth(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "alive!")
	if err != nil {
		return
	}
}

// GetStatsByType ...
func (h Handler) GetStatsByType(w http.ResponseWriter, r *http.Request) {
	statsType := mux.Vars(r)["type"]
	if statsType != "gauge" && statsType != "counter" {
		http.Error(w, "Can't save stat", http.StatusNotImplemented)
		return
	}

	id := mux.Vars(r)["id"]
	stat, err := h.Storage.GetStatsByID(id)
	if err != nil {
		http.Error(w, "Can't get stat by this ID", http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(stat); err != nil {
		http.Error(w, "unable to marshal the struct", http.StatusOK)
		return
	}
}
