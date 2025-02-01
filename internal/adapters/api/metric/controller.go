package metric

import (
	"github.com/julienschmidt/httprouter"
	"metrics-alerts/internal/adapters/api"
	"metrics-alerts/internal/domain/metric"
	"net/http"
)

const (
	metricsURL = "/metrics"
	metricURL  = "/metrics/:metric_name"
)

type handler struct {
	service metric.Service
}

func NewHandler(service metric.Service) api.Handler {
	return &handler{service: service}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(metricsURL, h.GetAll)
}

func (h *handler) GetAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	//h.GetAll(context.Background(), 0, 0)
	writer.Write([]byte("hehehe"))
	writer.WriteHeader(http.StatusOK)
}
