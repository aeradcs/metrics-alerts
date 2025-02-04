package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"metrics-alerts/config/server"
	server2 "metrics-alerts/internal/server"
	"net/http"
)

func main() {
	// args
	flag.Parse()
	fmt.Printf("Parsed args : a = %s\n", *server.Port)

	// env
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env", err)
	}

	// db & server
	db := server2.NewSQLMetricStorage()
	defer db.Close()
	handler := server2.Handler{Storage: db}
	router := mux.NewRouter()
	router.HandleFunc(`/update/{metric_type}/{metric_name}/{metric_value}`, handler.UpdateMetric)

	fmt.Printf("Server is running on :%s\n\n\n", *server.Port)
	if err := http.ListenAndServe(":"+*server.Port, router); err != nil {
		panic(err)
	}
}
