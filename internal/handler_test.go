package internal

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var MaxIDGauge = 0
var MaxIDCounter = 0

type mockDB struct{}

func (m *mockDB) GetMaxID(metricType string) (int, error) {
	if metricType == Gauge {
		return MaxIDGauge, nil
	} else if metricType == Counter {
		return MaxIDCounter, nil
	}
	return 0, nil
}

func (m *mockDB) InsertMetric(metricType, name string, value interface{}, id int) error {
	return nil
}

func (m *mockDB) UpdateMetricByID(metricType string, value interface{}, id int) error {
	return nil
}

func (m *mockDB) GetMetricIDByName(metricType, name string) (int, error) {
	if metricType == Gauge {
		return MaxIDGauge, nil
	} else if metricType == Counter {
		return MaxIDCounter, nil
	}
	return 0, nil
}

func TestShortenUrlAPI(t *testing.T) {
	tests := []struct {
		name         string
		handlerFunc  func(handler *Handler) http.HandlerFunc
		method       string
		urlParam     string
		url          string
		body         string
		responseCode int
		responseBody string
		createdObj   bool
		metricType   string
	}{
		{
			name:         "Create Gauge Success",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.UpdateMetric },
			method:       http.MethodPost,
			urlParam:     "/{metric_type}/{metric_name}/{metric_value}",
			url:          "/gauge/first/666.66",
			body:         "",
			responseCode: http.StatusOK,
			responseBody: "Successfully updated metric first",
			createdObj:   true,
			metricType:   Gauge,
		},
		{
			name:         "Update Gauge Success",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.UpdateMetric },
			method:       http.MethodPost,
			urlParam:     "/{metric_type}/{metric_name}/{metric_value}",
			url:          "/gauge/first/999.66",
			body:         "",
			responseCode: http.StatusOK,
			responseBody: "Successfully updated metric first",
			createdObj:   false,
			metricType:   Gauge,
		},
		{
			name:         "Update Gauge Error Value Is Not Float",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.UpdateMetric },
			method:       http.MethodPost,
			urlParam:     "/{metric_type}/{metric_name}/{metric_value}",
			url:          "/gauge/first/h",
			body:         "",
			responseCode: http.StatusBadRequest,
			responseBody: "Metric value is not a valid float\n",
			createdObj:   false,
			metricType:   Gauge,
		},
		{
			name:         "Update Gauge Error Wrong Method",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.UpdateMetric },
			method:       http.MethodPut,
			urlParam:     "/{metric_type}/{metric_name}/{metric_value}",
			url:          "/gauge/first/h",
			body:         "",
			responseCode: http.StatusMethodNotAllowed,
			responseBody: "Only POST requests are allowed!\n",
			createdObj:   false,
			metricType:   Gauge,
		},
		{
			name:         "Update Gauge Error Wrong Metric Type",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.UpdateMetric },
			method:       http.MethodPost,
			urlParam:     "/{metric_type}/{metric_name}/{metric_value}",
			url:          "/wrong/first/h",
			body:         "",
			responseCode: http.StatusBadRequest,
			responseBody: fmt.Sprintf("Metric type is invalid, possible types are: %s, provided type is: wrong\n", GetAllMetricTypesStr()),
			createdObj:   false,
			metricType:   Gauge,
		},
		// -------------------------------------------------------------
		{
			name:         "Create Counter Success",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.UpdateMetric },
			method:       http.MethodPost,
			urlParam:     "/{metric_type}/{metric_name}/{metric_value}",
			url:          "/counter/first/666",
			body:         "",
			responseCode: http.StatusOK,
			responseBody: "Successfully inserted metric first",
			createdObj:   true,
			metricType:   Counter,
		},
		{
			name:         "Create Counter Success 1",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.UpdateMetric },
			method:       http.MethodPost,
			urlParam:     "/{metric_type}/{metric_name}/{metric_value}",
			url:          "/counter/first/999",
			body:         "",
			responseCode: http.StatusOK,
			responseBody: "Successfully inserted metric first",
			createdObj:   true,
			metricType:   Counter,
		},
		{
			name:         "Update Counter Error Value Is Not Int",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.UpdateMetric },
			method:       http.MethodPost,
			urlParam:     "/{metric_type}/{metric_name}/{metric_value}",
			url:          "/counter/first/h",
			body:         "",
			responseCode: http.StatusBadRequest,
			responseBody: "Metric value is not a valid int\n",
			createdObj:   false,
			metricType:   Counter,
		},
		{
			name:         "Update Counter Error Value Is Not Int 1",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.UpdateMetric },
			method:       http.MethodPost,
			urlParam:     "/{metric_type}/{metric_name}/{metric_value}",
			url:          "/counter/first/77.77",
			body:         "",
			responseCode: http.StatusBadRequest,
			responseBody: "Metric value is not a valid int\n",
			createdObj:   false,
			metricType:   Counter,
		},
		{
			name:         "Update Counter Error Wrong Method",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.UpdateMetric },
			method:       http.MethodPut,
			urlParam:     "/{metric_type}/{metric_name}/{metric_value}",
			url:          "/counter/first/h",
			body:         "",
			responseCode: http.StatusMethodNotAllowed,
			responseBody: "Only POST requests are allowed!\n",
			createdObj:   false,
			metricType:   Counter,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockDB := &mockDB{}
			handler := &Handler{Storage: mockDB}
			router := mux.NewRouter()
			router.HandleFunc(test.urlParam, test.handlerFunc(handler)).Methods(test.method)

			req := httptest.NewRequest(test.method, test.url, strings.NewReader(test.body))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}
			assert.Equal(t, resp.StatusCode, test.responseCode)
			assert.Equal(t, string(body), test.responseBody)
			if test.metricType == Gauge && test.createdObj {
				MaxIDGauge += 1
			}
			if test.metricType == Counter && test.createdObj {
				MaxIDCounter += 1
			}
		})
	}
}
