package server

import (
	"database/sql"
	"fmt"
	"log"
	"metrics-alerts/internal/common"
	"os"
)

type MetricStorage interface {
	GetMaxID(metricType string) (int, error)
	InsertMetric(metricType, name string, value interface{}, id int) error
	UpdateMetricByID(metricType string, value interface{}, id int) error
	GetMetricIDByName(metricType, name string) (int, error)
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

func (s *SQLMetricStorage) InsertMetric(metricType, name string, value interface{}, id int) error {
	query := fmt.Sprintf("INSERT INTO %s VALUES ($1, $2, $3)", common.TableNames[metricType])
	_, err := s.DB.Exec(query, id+1, name, value)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLMetricStorage) UpdateMetricByID(metricType string, value interface{}, id int) error {
	query := fmt.Sprintf("UPDATE %s SET value = $1 WHERE id = $2", common.TableNames[metricType])
	_, err := s.DB.Exec(query, value, id)
	return err
}

func (s *SQLMetricStorage) GetMetricIDByName(metricType, name string) (int, error) {
	query := fmt.Sprintf("SELECT id FROM %s WHERE name = $1", common.TableNames[metricType])
	var id int
	err := s.DB.QueryRow(query, name).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return id, err
}
