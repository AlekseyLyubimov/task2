package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	addr := "localhost:8080"

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/hello", HelloHandler)
	mux.HandleFunc("/v1/time", CurrentTimeHandler)
	//wrap entire mux with logger middleware
	//wrappedMux := NewLogger(mux)

	log.Printf("server is listening at %s", addr)
	//use wrappedMux instead of mux as root handler
	log.Fatal(http.ListenAndServe(addr, Logogger(mux)))
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func CurrentTimeHandler(w http.ResponseWriter, r *http.Request) {
	curTime := time.Now().Format(time.Kitchen)
	w.Write([]byte(fmt.Sprintf("the current time is %v", curTime)))
}

func Logogger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Printf("Logger middleware says: %s %s %v", r.Method, r.URL.Path, time.Now().Format(time.StampMilli))
	}
}
