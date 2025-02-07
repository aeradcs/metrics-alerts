package common

import (
	"sort"
	"strings"
)

var Gauge = "gauge"
var Counter = "counter"

var TableNames = map[string]string{
	Gauge:   "gauge",
	Counter: "counter",
}

type Metric struct {
	ID         int
	Name       string
	MetricType string
	Value      interface{}
}

func NewMetric(id int, name, metricType string, value interface{}) *Metric {
	if IsValidValue(value) && IsValidMetricType(metricType) {
		return &Metric{
			ID:         id,
			Name:       name,
			MetricType: metricType,
			Value:      value,
		}
	}
	return nil
}

func IsValidValue(input interface{}) bool {
	switch input.(type) {
	case int:
		return true
	case int64:
		return true
	case float64:
		return true
	default:
		return false
	}
}

func IsValidMetricType(input string) bool {
	for key := range TableNames {
		if key == input {
			return true
		}
	}
	return false
}

func GetAllMetricTypesStr() string {
	keys := make([]string, 0, len(TableNames))
	for key := range TableNames {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
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
