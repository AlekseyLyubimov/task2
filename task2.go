package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"time"
)

func main() {
	addr := "localhost:8080"

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/hello", HelloHandler)
	mux.HandleFunc("/v1/encoding", AddEncoderToContext(EncodingRequestHandler))

	log.Printf("server is listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, LoggerMiddleware(mux)))
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func EncodingRequestHandler(w http.ResponseWriter, r *http.Request) {

	var employee = EmployeeInfo{123456, "Abbaddon", "Lupercal"}

	var encoder = r.Context().Value(requestDefinedEncoderKey)

	xmlEnc, ok := encoder.(*xml.Encoder)
	if ok {
		xmlEnc.Encode(employee)
	}

	jsonEnc, ok := encoder.(*json.Encoder)
	if ok {
		jsonEnc.Encode(employee)
	}

}

func LoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Printf("Logger middleware says: %s %s %v", r.Method, r.URL.Path, time.Now().Format(time.StampMilli))
	}
}

type encoderKey string

const requestDefinedEncoderKey encoderKey = "encoder"

func AddEncoderToContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Header.Get("content-type") {
		case "xml":
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestDefinedEncoderKey, xml.NewEncoder(w))))
		case "json":
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestDefinedEncoderKey, json.NewEncoder(w))))
		default:
			w.WriteHeader(400)
		}
	}

}

type EmployeeInfo struct {
	EmployeeId int
	FirstName  string
	LastName   string
}
