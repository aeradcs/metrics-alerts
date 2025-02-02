package internal

import (
	"fmt"
	"net/http"
)

type Handler struct {
	Database *Database
}

func (h *Handler) PostMetric(w http.ResponseWriter, req *http.Request) {
	rows, err := h.Database.DB.Query("SELECT * FROM metrics")
	if err != nil {
		fmt.Println("Error fetching metrics:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var a string
		var b string
		var c float64
		rows.Scan(&a, &b, &c)
		fmt.Printf("post: %s %s %f\n", a, b, c)
	}
	w.Write([]byte("abccccccccc"))
}

func (h *Handler) GetMetric(w http.ResponseWriter, req *http.Request) {
	rows, err := h.Database.DB.Query("SELECT * FROM metrics")
	if err != nil {
		fmt.Println("Error fetching metrics:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var a string
		var b string
		var c float64
		rows.Scan(&a, &b, &c)
		fmt.Printf("get: %s %s %f\n", a, b, c)
	}
	w.Write([]byte("heheheheh"))
}
