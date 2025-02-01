package metric

import (
	"context"
)

type Service interface {
	Save(ctx context.Context, dto *SaveMetricDTO) *Metric
	Update(ctx context.Context, dto *UpdateMetricDTO) *Metric
	Get(ctx context.Context, name string) *Metric
	GetAll(ctx context.Context, limit, offset int) []*Metric
}

type service struct {
	storage Storage
}

func NewService(storage Storage) Service {
	return &service{
		storage: storage,
	}
}

func (s *service) Save(ctx context.Context, dto *SaveMetricDTO) *Metric {
	return nil
}

func (s *service) Update(ctx context.Context, dto *UpdateMetricDTO) *Metric {
	return nil
}

func (s *service) Get(ctx context.Context, name string) *Metric {
	return s.storage.Get(name)
}

func (s *service) GetAll(ctx context.Context, offset, limit int) []*Metric {
	return s.storage.GetAll(offset, limit)
}
