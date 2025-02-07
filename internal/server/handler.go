package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"metrics-alerts/internal/common"
	"net/http"
	"strconv"
)

type Handler struct {
	Storage MetricStorage
}

func (h *Handler) UpdateMetric(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Received request : %s %s\n", req.Method, req.URL)
	if req.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	metricType := mux.Vars(req)["metric_type"]
	metricName := mux.Vars(req)["metric_name"]
	metricValue := mux.Vars(req)["metric_value"]
	if metricName == "" {
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}
	if metricType == "" || metricValue == "" {
		http.Error(w, "Metric type, name and value are required", http.StatusBadRequest)
		return
	}
	if !common.IsValidMetricType(metricType) {
		http.Error(w, fmt.Sprintf("Metric type is invalid, possible types are: %s, provided type is: %s", common.GetAllMetricTypesStr(), metricType), http.StatusBadRequest)
		return
	}
	var convertedMetricValue interface{}
	var err error
	if metricType == common.Gauge {
		convertedMetricValue, err = strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Metric value is not a valid float", http.StatusBadRequest)
			return
		}
	} else if metricType == common.Counter {
		convertedMetricValue, err = strconv.Atoi(metricValue)
		if err != nil {
			http.Error(w, "Metric value is not a valid int", http.StatusBadRequest)
			return
		}
	} else {
		convertedMetricValue, err = 0, nil
	}

	if metricType == common.Gauge {
		created, err := h.replaceValue(metricType, metricName, convertedMetricValue)
		if err != nil {
			http.Error(w, "Error occurred during updating metric "+err.Error(), http.StatusInternalServerError)
			return
		}
		if created {
			w.Write([]byte("Successfully inserted metric " + metricName))
		} else {
			w.Write([]byte("Successfully updated metric " + metricName))
		}

	} else if metricType == common.Counter {
		err := h.addValue(metricType, metricName, convertedMetricValue)
		if err != nil {
			http.Error(w, "Error occurred during inserting metric "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Successfully inserted metric " + metricName))
	}
}

func (h *Handler) addValue(metricType, name string, value interface{}) error {
	maxID, err := h.Storage.GetMaxID(metricType)
	if err != nil {
		return err
	}
	err = h.Storage.InsertMetric(common.NewMetric(maxID+1, name, metricType, value))
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) replaceValue(metricType, name string, value interface{}) (bool, error) {
	maxID, err := h.Storage.GetMaxID(metricType)
	if err != nil {
		return false, err
	}
	existingMetric, err := h.Storage.GetMetric(metricType, name)
	if err != nil {
		return false, err
	}
	if existingMetric == nil {
		h.Storage.InsertMetric(common.NewMetric(maxID+1, name, metricType, value))
		return true, nil
	}
	existingMetric.Value = value
	h.Storage.UpdateMetric(existingMetric)
	return false, nil
}

func (h *Handler) GetAllMetrics(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Received request : %s %s\n", req.Method, req.URL)
	if req.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return

	}
	metrics, err := h.Storage.GetAllMetrics()
	if err != nil {
		http.Error(w, "Error occurred during retrieving metrics "+err.Error(), http.StatusInternalServerError)
		return
	}
	for _, m := range metrics {
		w.Write([]byte(fmt.Sprintf("Metric of type %s: %s = %v\n", m.MetricType, m.Name, m.Value)))
	}
}

func (h *Handler) GetMetric(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Received request : %s %s\n", req.Method, req.URL)
	if req.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return

	}
	metricType := mux.Vars(req)["metric_type"]
	metricName := mux.Vars(req)["metric_name"]
	if metricType == "" || metricName == "" {
		http.Error(w, "Metric type and name are required", http.StatusBadRequest)
		return
	}
	if !common.IsValidMetricType(metricType) {
		http.Error(w, fmt.Sprintf("Metric type is invalid, possible types are: %s, provided type is: %s", common.GetAllMetricTypesStr(), metricType), http.StatusBadRequest)
		return
	}
	existingMetric, err := h.Storage.GetMetric(metricType, metricName)
	if existingMetric == nil {
		http.Error(w, fmt.Sprintf("Metric %s of type %s not found", metricName, metricType), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Error occurred during retrieving metric: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("Metric %s value is %v\n", existingMetric.Name, existingMetric.Value)))
}
