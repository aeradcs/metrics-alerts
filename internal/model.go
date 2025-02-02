package internal

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
