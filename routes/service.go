package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/anandkumarpatel/main/models"
	"github.com/gorilla/mux"
)

// AddServiceRoutes configures all routes for the Service model
func AddServiceRoutes(r *mux.Router, sm models.ServiceModel) {
	r.HandleFunc("/service", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := r.URL.Query()
		filter, err := createFilterFromQuery(params)
		if err != nil {
			sendErr(err, http.StatusBadRequest, w)
			return
		}

		services, err := sm.ListServiesWithFilter(r.Context(), filter)
		if err != nil {
			sendErr(err, http.StatusFailedDependency, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string][]*models.Service{
			"services": services,
		})

	}).Methods("GET")

	r.HandleFunc("/service/{name}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		name, exists := vars["name"]
		if !exists {
			sendErr(fmt.Errorf("nothing to lookup"), http.StatusBadRequest, w)
			return
		}
		service, err := sm.GetServiceByName(r.Context(), name)
		if err != nil {
			sendErr(err, http.StatusFailedDependency, w)
			return
		}
		if service == nil {
			sendErr(fmt.Errorf("service not found"), http.StatusNotFound, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]models.Service{
			"service": *service,
		})
	}).Methods("GET")
}

func createFilterFromQuery(params map[string][]string) (*models.ServiceFilter, error) {
	filter := models.ServiceFilter{}

	if v, ok := params["search"]; ok {
		if len(v) != 1 {
			return nil, fmt.Errorf("search param specified more then once: %v", v)
		}
		filter.Search = v[0]
	}

	if v, ok := params["sort"]; ok {
		if len(v) != 1 {
			return nil, fmt.Errorf("sort param specified more then once: %v", v)
		}
		filter.Sort = v[0]
	}

	if v, ok := params["page_size"]; ok {
		size, err := strconv.Atoi(v[0])
		if err != nil {
			return nil, err
		}
		filter.PageSize = size
	}

	if v, ok := params["page_number"]; ok {
		size, err := strconv.Atoi(v[0])
		if err != nil {
			return nil, err
		}
		filter.PageNumber = size
	}

	if err := filter.Validate(); err != nil {
		return nil, err
	}

	return &filter, nil
}

func sendErr(err error, code int, w http.ResponseWriter) {
	output := map[string]string{}
	w.WriteHeader(code)
	output["error"] = err.Error()
	json.NewEncoder(w).Encode(output)
}
