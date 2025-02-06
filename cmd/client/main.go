package main

import (
	clientrequests "metrics-alerts/internal/client"
	"net/http"
)

func main() {
	httpClient := &http.Client{}
	memStats := clientrequests.ExtendedMemStats{}

	go clientrequests.GatherStats(&memStats)
	go clientrequests.SendStats(&memStats, httpClient)

	select {}
}
