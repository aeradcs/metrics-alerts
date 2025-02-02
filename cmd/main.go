package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"metrics-alerts/config"
	"metrics-alerts/internal"
	"net/http"
	"os"
)

func main() {
	// args
	flag.Parse()
	fmt.Printf("Parsed args : a = %s\n", *config.Port)

	// env
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	// db
	connStr := os.Getenv("DB_CONN_STR")
	if connStr == "" {
		panic("DB_CONN_STR environment variable is not set")
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// server
	router := mux.NewRouter()
	router.HandleFunc(`/`, internal.PostMetric)
	router.HandleFunc(`/{name}`, internal.GetMetric)

	fmt.Printf("Server is running on :%s\n", *config.Port)
	if err := http.ListenAndServe(":"+*config.Port, router); err != nil {
		panic(err)
	}
}
