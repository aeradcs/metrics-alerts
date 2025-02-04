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
		err := h.replaceValue(metricType, metricName, convertedMetricValue)
		if err != nil {
			http.Error(w, "Error occurred during updating metric "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Successfully updated metric " + metricName))
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

	err = h.Storage.InsertMetric(metricType, name, value, maxID+1)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) replaceValue(metricType, name string, value interface{}) error {
	maxID, err := h.Storage.GetMaxID(metricType)
	if err != nil {
		return err
	}

	existingMetricID, err := h.Storage.GetMetricIDByName(metricType, name)
	if err != nil {
		return err
	}
	if existingMetricID == 0 {
		h.Storage.InsertMetric(metricType, name, value, maxID+1)
	} else {
		h.Storage.UpdateMetricByID(metricType, value, maxID)
	}
	return nil
}
