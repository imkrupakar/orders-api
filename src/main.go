package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ready")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s up\n", os.Getenv("SERVICE_NAME"))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
