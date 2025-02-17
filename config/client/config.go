package client

import "flag"

var (
	BaseUrl        = "http://localhost"
	Port           = flag.String("a", "8080", "port")
	ReportInterval = flag.String("r", "10", "interval to send stats on server")
	PollInterval   = flag.String("p", "2", "interval to update stats on client")
)
