package main

import (
	"encoding/json"
	"net/http"

	"github.com/rhino1998/umdplanner/testudo"
)

type server struct {
	store testudo.ClassStore
}

type courseQuery struct {
}

func (s *server) QueryCourses(w http.ResponseWriter, r *http.Request) {
	dec, err := json.NewDecoder(r.Body)

	dec.Decode()
}
