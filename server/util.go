package main

import (
	"fmt"
	"net/http"
)

func send501(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "%v", err)
}
