package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"github.com/kiraplenkin/go-musthave-devops/internal/validator"
	"net/http"
	"strconv"
)

// Handler stores pointers to service
type Handler struct {
	Router  *mux.Router
	Service *storage.Store
}

// NewHandler returns a pointer to Handler
func NewHandler(s storage.Store) *Handler {
	return &Handler{
		Service: &s,
	}
}

// SetupRouters sets up all the routes for server
func (h *Handler) SetupRouters() {
	h.Router = mux.NewRouter()

	h.Router.HandleFunc("/api/stat/{id}", h.GetStatsByID).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/stat/", h.GetAllStats).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/stat/", h.PostStat).Methods(http.MethodPost)

	h.Router.HandleFunc("/api/health/", h.CheckHealth).Methods(http.MethodGet)
}

//GetStatsByID handler that return types.Stats by ID
func (h *Handler) GetStatsByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse uint from id", http.StatusBadRequest)
		return
	}

	stat, err := h.Service.GetStatsByID(uint(i))
	if err != nil {
		http.Error(w, "Can't get stat by this ID", http.StatusBadRequest)
		return
	}

	_, err = fmt.Fprintf(w, "%+v", stat)
	if err != nil {
		http.Error(w, error.Error(err), http.StatusInternalServerError)
	}
}

//GetAllStats handler that return all values from storage.Store
func (h Handler) GetAllStats(w http.ResponseWriter, _ *http.Request) {
	allStats, err := h.Service.GetAllStats()
	if err != nil {
		http.Error(w, "Can't get all stats", http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprintf(w, "%v", *allStats)
	if err != nil {
		return
	}
}

//PostStat handler that save types.Stats to storage.Store
func (h Handler) PostStat(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "can't parse form", http.StatusInternalServerError)
	}

	id := r.Form.Get("id")
	alloc := r.Form.Get("Alloc")
	totalAlloc := r.Form.Get("TotalAlloc")
	sys := r.Form.Get("Sys")
	mallocs := r.Form.Get("Mallocs")
	frees := r.Form.Get("Frees")
	liveObjects := r.Form.Get("LiveObjects")
	pauseTotalNs := r.Form.Get("PauseTotalNs")
	numGC := r.Form.Get("NumGC")
	numGoroutine := r.Form.Get("NumGoroutine")

	err = validator.Require(
		id, alloc, totalAlloc, sys, mallocs, frees, liveObjects, pauseTotalNs, numGC, numGoroutine,
	)
	if err != nil {
		http.Error(w, error.Error(err), http.StatusBadRequest)
		return
	}

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse uint from id", http.StatusBadRequest)
		return
	}

	newStat := types.Stats{
		Alloc:        validator.Transform(alloc),
		TotalAlloc:   validator.Transform(totalAlloc),
		Sys:          validator.Transform(sys),
		Mallocs:      validator.Transform(mallocs),
		Frees:        validator.Transform(frees),
		LiveObjects:  validator.Transform(liveObjects),
		PauseTotalNs: validator.Transform(pauseTotalNs),
		NumGC:        validator.Transform(numGC),
		NumGoroutine: validator.Transform(numGoroutine),
	}
	err = h.Service.SaveStats(uint(i), newStat)
	if err != nil {
		http.Error(w, "Can't save stat", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = fmt.Fprintf(w, "%+v", newStat)
	if err != nil {
		return
	}
	fmt.Printf("%+v", newStat)
}

// CheckHealth handler to check health
func (h Handler) CheckHealth(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "alive!")
	if err != nil {
		return
	}
}
