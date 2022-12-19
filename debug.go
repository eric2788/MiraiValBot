package main

import (
	"net/http"
	_ "net/http/pprof"

	log "github.com/sirupsen/logrus"
)

func debugServe() {
	if err := http.ListenAndServe("0.0.0.0:45678", http.DefaultServeMux); err != nil {
		log.Fatal(err)
	}
}
