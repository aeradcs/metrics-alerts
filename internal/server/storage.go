package server

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"metrics-alerts/internal/common"
	"os"
)

type MetricStorage interface {
	GetMaxID(metricType string) (int, error)
	InsertMetric(*common.Metric) error
	UpdateMetric(*common.Metric) error
	GetMetric(metricType, name string) (*common.Metric, error)
	GetAllMetrics() ([]*common.Metric, error)
}

type SQLMetricStorage struct {
	DB *sql.DB
}

func NewSQLMetricStorage() *SQLMetricStorage {
	connStr := os.Getenv("DB_SERVER_CONN_STR")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Database connected")
	return &SQLMetricStorage{DB: db}
}

func (s *SQLMetricStorage) Close() {
	s.DB.Close()
	fmt.Println("Database closed")
}

func (s *SQLMetricStorage) GetMaxID(metricType string) (int, error) {
	query := fmt.Sprintf("SELECT coalesce(max(id), 0) FROM %s", common.TableNames[metricType])
	var maxID int
	err := s.DB.QueryRow(query).Scan(&maxID)
	if err != nil {
		return 0, err
	}
	return maxID, nil
}

func (s *SQLMetricStorage) InsertMetric(metric *common.Metric) error {
	query := fmt.Sprintf("INSERT INTO %s VALUES ($1, $2, $3)", common.TableNames[metric.MetricType])
	_, err := s.DB.Exec(query, metric.ID, metric.Name, metric.Value)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLMetricStorage) UpdateMetric(metric *common.Metric) error {
	query := fmt.Sprintf("UPDATE %s SET value = $1 WHERE id = $2", common.TableNames[metric.MetricType])
	_, err := s.DB.Exec(query, metric.Value, metric.ID)
	return err
}

func (s *SQLMetricStorage) GetMetric(metricType, name string) (*common.Metric, error) {
	query := fmt.Sprintf("SELECT id, value FROM %s WHERE name = $1", common.TableNames[metricType])
	var id int
	var value interface{}
	err := s.DB.QueryRow(query, name).Scan(&id, &value)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return common.NewMetric(id, name, metricType, value), err
}

func (s *SQLMetricStorage) GetAllMetrics() ([]*common.Metric, error) {
	var metrics []*common.Metric
	for key := range common.TableNames {
		query := fmt.Sprintf("SELECT * FROM %s", common.TableNames[key])
		rows, err := s.DB.Query(query)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			var name string
			var value interface{}
			err := rows.Scan(&id, &name, &value)
			if err != nil {
				return nil, err
			}
			m := &common.Metric{
				MetricType: key,
				ID:         id,
				Name:       name,
				Value:      value,
			}
			metrics = append(metrics, m)
		}
	}
	return metrics, nil
}
