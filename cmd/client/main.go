package main

import (
	"fmt"
	clientrequests "metrics-alerts/internal/client"
	"metrics-alerts/internal/common"
	"net/http"
	"reflect"
	"runtime"
)

func main() {
	var metricType = common.Gauge
	httpClient := &http.Client{}

	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	v := reflect.ValueOf(memStats)
	t := reflect.TypeOf(memStats)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if common.MemStatsFields[field.Name] {
			fmt.Printf("Field %s -- Value %v\n", field.Name, v.Field(i).Interface())
			_, _ = clientrequests.SendUpdateRequest(httpClient, metricType, field.Name, v.Field(i).Interface())
		}
	}

}
