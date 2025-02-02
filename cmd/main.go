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
	db := internal.NewDatabase()
	defer db.Close()
	handler := internal.Handler{Database: db}
	router := mux.NewRouter()
	router.HandleFunc(`/`, handler.PostMetric)
	router.HandleFunc(`/{name}`, handler.GetMetric)

	fmt.Printf("Server is running on :%s\n", *config.Port)
	if err := http.ListenAndServe(":"+*config.Port, router); err != nil {
		panic(err)
	}
}
