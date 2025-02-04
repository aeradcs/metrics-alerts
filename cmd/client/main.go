package main

import (
	"bytes"
	"fmt"
	"io"
	"metrics-alerts/config/client"
	"metrics-alerts/internal/common"
	"net/http"
)

func main() {
	var metricType = common.Gauge
	var metricName = "abc"
	var metricValue = 10

	httpClient := &http.Client{}
	url := fmt.Sprintf("%s/update/%s/%s/%d", client.BaseUrl, metricType, metricName, metricValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(nil))
	if err != nil {
		fmt.Printf("Error while creating request: %v\n", err)
	}
	fmt.Printf("Sending request to the server: %v\n", req)
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("Error while sending request: %v -- %v\n", req, err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Got response from the server: %s\n", string(body))
}
