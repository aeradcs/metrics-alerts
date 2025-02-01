package metric

type Storage interface {
	Save(metric *Metric) *Metric
	Update(metric *Metric) *Metric
	Get(name string) *Metric
	GetAll(limit, offset int) []*Metric
}
