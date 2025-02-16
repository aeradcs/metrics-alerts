package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io"
	"metrics-alerts/internal/common"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var MaxIDGauge = 0
var MaxIDCounter = 0

type mockDB struct{}

func (m *mockDB) InsertMetric(metric *common.Metric) error {
	return nil
}

func (m *mockDB) UpdateMetric(metric *common.Metric) error {
	return nil
}

func (m *mockDB) GetMetric(metricType, name string) (*common.Metric, error) {
	if metricType == common.Gauge {
		if MaxIDGauge == 0 || name == "not_exists" {
			return nil, nil
		}
		return common.NewMetric(MaxIDGauge, name, metricType, 20), nil
	} else if metricType == common.Counter {
		if MaxIDCounter == 0 || name == "not_exists" {
			return nil, nil
		}
		return common.NewMetric(MaxIDCounter, name, metricType, 10), nil
	}
	return common.NewMetric(0, name, metricType, 10), nil
}

func (m *mockDB) GetAllMetrics() ([]*common.Metric, error) {
	return []*common.Metric{
		common.NewMetric(1, "a", "counter", 10),
		common.NewMetric(2, "b", "counter", 20),
		common.NewMetric(1, "aaa", "gauge", 10),
		common.NewMetric(2, "aaa", "gauge", 20),
		common.NewMetric(3, "bbb", "gauge", 30),
	}, nil
}

func (m *mockDB) GetMaxID(metricType string) (int, error) {
	if metricType == common.Gauge {
		return MaxIDGauge, nil
	} else if metricType == common.Counter {
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
			responseBody: "Successfully inserted metric first",
			createdObj:   true,
			metricType:   common.Gauge,
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
			metricType:   common.Gauge,
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
			metricType:   common.Gauge,
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
			metricType:   common.Gauge,
		},
		{
			name:         "Update Gauge Error Wrong Metric Type",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.UpdateMetric },
			method:       http.MethodPost,
			urlParam:     "/{metric_type}/{metric_name}/{metric_value}",
			url:          "/wrong/first/h",
			body:         "",
			responseCode: http.StatusBadRequest,
			responseBody: fmt.Sprintf("Metric type is invalid, possible types are: %s, provided type is: wrong\n", common.GetAllMetricTypesStr()),
			createdObj:   false,
			metricType:   common.Gauge,
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
			metricType:   common.Counter,
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
			metricType:   common.Counter,
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
			metricType:   common.Counter,
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
			metricType:   common.Counter,
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
			metricType:   common.Counter,
		},
		// -------------------------------------------------------------
		{
			name:         "Get Counter Metric Success",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.GetMetric },
			method:       http.MethodGet,
			urlParam:     "/{metric_type}/{metric_name}",
			url:          "/counter/first",
			body:         "",
			responseCode: http.StatusOK,
			responseBody: "Metric first value is 10\n",
			createdObj:   false,
			metricType:   common.Counter,
		},
		{
			name:         "Get Counter Metric Not Found",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.GetMetric },
			method:       http.MethodGet,
			urlParam:     "/{metric_type}/{metric_name}",
			url:          "/counter/not_exists",
			body:         "",
			responseCode: http.StatusNotFound,
			responseBody: "Metric not_exists of type counter not found\n",
			createdObj:   false,
			metricType:   common.Counter,
		},
		{
			name:         "Get Counter Metric Wrong Type",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.GetMetric },
			method:       http.MethodGet,
			urlParam:     "/{metric_type}/{metric_name}",
			url:          "/wrong/first",
			body:         "",
			responseCode: http.StatusBadRequest,
			responseBody: "Metric type is invalid, possible types are: counter, gauge, provided type is: wrong\n",
			createdObj:   false,
			metricType:   common.Counter,
		},
		{
			name:         "Get Gauge Metric Success",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.GetMetric },
			method:       http.MethodGet,
			urlParam:     "/{metric_type}/{metric_name}",
			url:          "/gauge/first",
			body:         "",
			responseCode: http.StatusOK,
			responseBody: "Metric first value is 20\n",
			createdObj:   false,
			metricType:   common.Gauge,
		},
		{
			name:         "Get Gauge Metric Not Found",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.GetMetric },
			method:       http.MethodGet,
			urlParam:     "/{metric_type}/{metric_name}",
			url:          "/gauge/not_exists",
			body:         "",
			responseCode: http.StatusNotFound,
			responseBody: "Metric not_exists of type gauge not found\n",
			createdObj:   false,
			metricType:   common.Gauge,
		},
		{
			name:         "Get Gauge Metric Wrong Type",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.GetMetric },
			method:       http.MethodGet,
			urlParam:     "/{metric_type}/{metric_name}",
			url:          "/wrong/first",
			body:         "",
			responseCode: http.StatusBadRequest,
			responseBody: "Metric type is invalid, possible types are: counter, gauge, provided type is: wrong\n",
			createdObj:   false,
			metricType:   common.Gauge,
		},
		{
			name:         "Get Gauge Metric Wrong Method",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.GetMetric },
			method:       http.MethodPut,
			urlParam:     "/{metric_type}/{metric_name}",
			url:          "/wrong/first",
			body:         "",
			responseCode: http.StatusMethodNotAllowed,
			responseBody: "Only GET requests are allowed!\n",
			createdObj:   false,
			metricType:   common.Gauge,
		},
		// -------------------------------------------------------------
		{
			name:         "Get All Success",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.GetAllMetrics },
			method:       http.MethodGet,
			urlParam:     "/",
			url:          "/",
			body:         "",
			responseCode: http.StatusOK,
			responseBody: "Metric of type counter: a = 10\nMetric of type counter: b = 20\nMetric of type gauge: aaa = 10\nMetric of type gauge: aaa = 20\nMetric of type gauge: bbb = 30\n",
			createdObj:   false,
			metricType:   common.Counter,
		},
		{
			name:         "Get All Wrong Method",
			handlerFunc:  func(h *Handler) http.HandlerFunc { return h.GetAllMetrics },
			method:       http.MethodPut,
			urlParam:     "/",
			url:          "/",
			body:         "",
			responseCode: http.StatusMethodNotAllowed,
			responseBody: "Only GET requests are allowed!\n",
			createdObj:   false,
			metricType:   common.Counter,
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
			if test.metricType == common.Gauge && test.createdObj {
				MaxIDGauge += 1
			}
			if test.metricType == common.Counter && test.createdObj {
				MaxIDCounter += 1
			}
		})
	}
}
