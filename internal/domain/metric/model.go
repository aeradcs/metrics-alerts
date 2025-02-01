package metric

import "errors"

type Metric struct {
	MetricType string
	Name       string
	Value      interface{}
}

const (
	Gauge   = "gauge"
	Counter = "counter"
)

func NewMetric(name string, metricType string, value interface{}) (*Metric, error) {
	if metricType != Gauge && metricType != Counter {
		return nil, errors.New("invalid metric type: " + string(metricType))
	}
	return &Metric{
		Name:       name,
		MetricType: metricType,
		Value:      value,
	}, nil
}
