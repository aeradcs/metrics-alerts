package client

import (
	"bytes"
	"fmt"
	"io"
	"metrics-alerts/config/client"
	"net/http"
)

func SendUpdateRequest(httpClient *http.Client, metricType, metricName string, metricValue interface{}) (string, error) {
	url := fmt.Sprintf("%s/update/%s/%s/%v", client.BaseUrl, metricType, metricName, metricValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(nil))
	if err != nil {
		fmt.Printf("Error while creating request: %v\n", err)
		return "", err
	}
	fmt.Printf("Sending request to the server: %v %v\n", req.Method, req.URL)
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("Error while sending request: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Got response from the server: %s\n", string(body))
	return string(body), nil
}
