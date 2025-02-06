package client

import (
	"fmt"
	"math/rand"
	"metrics-alerts/config/client"
	"metrics-alerts/internal/common"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

type ExtendedMemStats struct {
	runtime.MemStats
	PollCount   int64
	RandomValue float64
}

func SendStats(memStats *ExtendedMemStats, httpClient *http.Client) {
	for {
		v := reflect.ValueOf(memStats.MemStats)
		t := reflect.TypeOf(memStats.MemStats)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if common.MemStatsFields[field.Name] {
				_, _ = SendUpdateRequest(httpClient, common.Gauge, field.Name, v.Field(i).Interface())
				memStats.PollCount++
				_, _ = SendUpdateRequest(httpClient, common.Counter, "PollCount", memStats.PollCount)
				memStats.RandomValue = rand.Float64() * float64(rand.Int31n(1000))
				_, _ = SendUpdateRequest(httpClient, common.Gauge, "RandomValue", memStats.RandomValue)
			}
		}
		fmt.Println("\nSending stats done, going to sleep...\n")
		time.Sleep(time.Duration(client.ReportInterval) * time.Second)
	}
}

func GatherStats(memStats *ExtendedMemStats) {
	for {
		runtime.ReadMemStats(&memStats.MemStats)
		fmt.Println("\nGathering stats done, going to sleep...\n")
		time.Sleep(time.Duration(client.PollInterval) * time.Second)
	}
}
