package internal

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	Database *Database
}

func (h *Handler) UpdateMetric(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	metricType := mux.Vars(req)["metric_type"]
	metricName := mux.Vars(req)["metric_name"]
	metricValue := mux.Vars(req)["metric_value"]
	fmt.Printf("Received params from URL : type = %s, name = %s, value = %s\n", metricType, metricName, metricValue)
	if metricName == "" {
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}
	if metricType == "" || metricValue == "" {
		http.Error(w, "Metric type, name and value are required", http.StatusBadRequest)
		return
	}
	if !IsValidMetricType(metricType) {
		http.Error(w, fmt.Sprintf("Metric type is invalid, possible types are: %s", GetAllMetricTypesStr()), http.StatusBadRequest)
		return
	}

	if metricType == Gauge {
		err := ReplaceValue(h.Database.DB, metricType, metricName, metricValue)
		if err != nil {
			http.Error(w, "Error occurred during updating metric "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Successfully updated metric " + metricName))
	}
	if metricType == Counter {
		err := AddValue(h.Database.DB, metricType, metricName, metricValue)
		if err != nil {
			http.Error(w, "Error occurred during inserting metric "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Successfully inserted metric " + metricName))
	}
}

func AddValue(db *sql.DB, metricType, name, value string) error {
	query := fmt.Sprintf("SELECT max(id) FROM %s", TableNames[metricType])
	var maxID int
	err := db.QueryRow(query).Scan(&maxID)
	if err != nil {
		return err
	}

	query = fmt.Sprintf("INSERT INTO %s VALUES ($1, $2, $3)", TableNames[metricType])
	_, err = db.Exec(query, maxID+1, name, value)
	if err != nil {
		return err
	}
	return nil
}

func ReplaceValue(db *sql.DB, metricType, name, value string) error {
	query := fmt.Sprintf("SELECT max(id) FROM %s", TableNames[metricType])
	var maxID int
	err := db.QueryRow(query).Scan(&maxID)
	if err != nil {
		return err
	}

	query = fmt.Sprintf("SELECT * FROM %s WHERE name = $1", TableNames[metricType])
	rows, err := db.Query(query, name)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var metricName string
		var metricValue string
		if err := rows.Scan(&id, &metricName, &metricValue); err != nil {
			return err
		}
		fmt.Println("ID:", id, "Name:", metricName, "Value:", metricValue)

		query = fmt.Sprintf("UPDATE %s SET value = $1 WHERE id = $2", TableNames[metricType])
		_, err := db.Exec(query, value, id)
		if err != nil {
			return err
		}

		return nil
	}

	query = fmt.Sprintf("INSERT INTO %s VALUES ($1, $2, $3)", TableNames[metricType])
	_, err = db.Exec(query, maxID+1, name, value)
	if err != nil {
		return err
	}
	return nil
}
