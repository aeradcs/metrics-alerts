package metric

type SaveMetricDTO struct {
	MetricType string
	Name       string
	Value      interface{}
}

type UpdateMetricDTO struct {
	MetricType string
	Name       string
	Value      interface{}
}
