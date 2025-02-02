package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"metrics-alerts/config"
	"metrics-alerts/internal"
	"net/http"
)

func main() {
	// args
	flag.Parse()
	fmt.Printf("Parsed args : a = %s\n", *config.Port)

	// env
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env", err)
	}

	// db & server
	db := internal.NewSQLMetricStorage()
	defer db.Close()
	handler := internal.Handler{Storage: db}
	router := mux.NewRouter()
	router.HandleFunc(`/update/{metric_type}/{metric_name}/{metric_value}`, handler.UpdateMetric)

	fmt.Printf("Server is running on :%s\n\n\n", *config.Port)
	if err := http.ListenAndServe(":"+*config.Port, router); err != nil {
		panic(err)
	}
}
