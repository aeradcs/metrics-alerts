package main

import (
	"flag"
	"fmt"
	"metrics-alerts/config/client"
	clientrequests "metrics-alerts/internal/client"
	"net/http"
	"strconv"
)

func main() {
	// args
	flag.Parse()
	fmt.Printf("Parsed args : a = %s, p = %s, r = %s\n", *client.Port, *client.PollInterval, *client.ReportInterval)

	httpClient := &http.Client{}
	memStats := clientrequests.ExtendedMemStats{}

	pollInterval, err := strconv.Atoi(*client.PollInterval)
	if err != nil {
		fmt.Printf("Error parsing poll interval: %s\n", err)
	}
	reportInterval, err := strconv.Atoi(*client.ReportInterval)
	if err != nil {
		fmt.Printf("Error parsing report interval: %s\n", err)
	}
	go clientrequests.GatherStats(&memStats, int64(pollInterval))
	go clientrequests.SendStats(&memStats, httpClient, int64(reportInterval))

	select {}
}
