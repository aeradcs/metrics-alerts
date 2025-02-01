package metric

import "metrics-alerts/internal/domain/metric"

type metricStorage struct {
}

func (m *metricStorage) Save(metric *metric.Metric) *metric.Metric {
	return nil
}

func (m *metricStorage) Update(metric *metric.Metric) *metric.Metric {
	return nil
}

func (m *metricStorage) Get(name string) *metric.Metric {
	return nil
}

func (m *metricStorage) GetAll(limit, offset int) []*metric.Metric {
	return nil
}
