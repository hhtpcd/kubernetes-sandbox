package main

import (
	"log"
	"net/http"
	"runtime/debug"
)

var (
	// Version is the current version of the application.
	VERSION, GITSHA, BUILD_ID string
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", sendOk)
	mux.HandleFunc("/healthz", healthz)

	info, ok := debug.ReadBuildInfo()
	if !ok {
		log.Fatalln("failed to read build info")
	}

	log.Printf("%+v\n", info)

	log.Println("starting listener on port 8080")
	log.Fatalln(http.ListenAndServe(":8080", mux))
}

func sendOk(w http.ResponseWriter, r *http.Request) {
	log.Printf("received request from client: %s\n", r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
