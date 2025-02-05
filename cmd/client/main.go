package main

import (
	"math/rand"
	clientrequests "metrics-alerts/internal/client"
	"metrics-alerts/internal/common"
	"net/http"
	"reflect"
	"runtime"
)

func main() {
	httpClient := &http.Client{}

	memStats := clientrequests.ExtendedMemStats{}
	runtime.ReadMemStats(&memStats.MemStats)
	v := reflect.ValueOf(memStats.MemStats)
	t := reflect.TypeOf(memStats.MemStats)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if common.MemStatsFields[field.Name] {
			_, _ = clientrequests.SendUpdateRequest(httpClient, common.Gauge, field.Name, v.Field(i).Interface())
			memStats.PollCount++
			_, _ = clientrequests.SendUpdateRequest(httpClient, common.Counter, "PollCount", memStats.PollCount)
			memStats.RandomValue = rand.Float64() * float64(rand.Int31n(1000))
			_, _ = clientrequests.SendUpdateRequest(httpClient, common.Gauge, "RandomValue", memStats.RandomValue)
		}
	}

}
