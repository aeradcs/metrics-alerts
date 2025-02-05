package common

import "strings"

var Gauge = "gauge"
var Counter = "counter"
var MetricTypes = []string{Gauge, Counter}
var TableNames = map[string]string{
	Gauge:   "gauge",
	Counter: "counter",
}

func IsValidMetricType(input string) bool {
	for _, valid := range MetricTypes {
		if valid == input {
			return true
		}
	}
	return false
}

func GetAllMetricTypesStr() string {
	return strings.Join(MetricTypes, ", ")
}

var MemStatsFields = map[string]bool{
	"Alloc":         true,
	"BuckHashSys":   true,
	"Frees":         true,
	"GCCPUFraction": true,
	"GCSys":         true,
	"HeapAlloc":     true,
	"HeapIdle":      true,
	"HeapInuse":     true,
	"HeapObjects":   true,
	"HeapReleased":  true,
	"HeapSys":       true,
	"LastGC":        true,
	"Lookups":       true,
	"MCacheInuse":   true,
	"MCacheSys":     true,
	"MSpanInuse":    true,
	"MSpanSys":      true,
	"Mallocs":       true,
	"NextGC":        true,
	"NumForcedGC":   true,
	"NumGC":         true,
	"OtherSys":      true,
	"PauseTotalNs":  true,
	"StackInuse":    true,
	"StackSys":      true,
	"Sys":           true,
	"TotalAlloc":    true,
}
