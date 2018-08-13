package main

import (
	"encoding/json"
	"fmt"
	"github.com/oshankkumar/prometheus-exporter/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type StringSvc struct {
	S string `json:"s"`
}

type Resp struct {
	S string `json:"s"`
}

type ErrResp struct {
	Error string `json:"error"`
}

type CountResp struct {
	Count int `json:"count"`
}

func reverse(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	var req = &StringSvc{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		SendError(w, err)
		return
	}
	str := []rune(req.S)
	for lindex, rindex := 0, len(str)-1; lindex < rindex; lindex, rindex = lindex+1, rindex-1 {
		str[lindex], str[rindex] = str[rindex], str[lindex]
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&Resp{string(str)})
}

func count(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	var req = &StringSvc{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		SendError(w,err)
		return
	}
	substr := r.URL.Query().Get("substr")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CountResp{strings.Count(req.S, substr)})
}

func SendError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(&ErrResp{err.Error()})
	return
}

func main() {
	http.HandleFunc("/count", count)
	http.HandleFunc("/reverse", reverse)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("listening on 8080")
	http.ListenAndServe(":8080", middleware.Prometheus(http.DefaultServeMux))
}
