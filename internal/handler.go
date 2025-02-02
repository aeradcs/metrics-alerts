package internal

import (
	"net/http"
)

func PostMetric(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("abccccccccc"))
}

func GetMetric(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("heheheheh"))
}
