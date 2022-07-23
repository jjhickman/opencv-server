package info

import (
	"log"
	"net/http"

	"gocv.io/x/gocv"
)

func Info(w http.ResponseWriter, r *http.Request) {
	log.Print(gocv.OpenCVVersion())
}
