package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rhino1998/umdplanner/testudo"
)

type server struct {
	store testudo.ClassStore
}

func newServer(store testudo.ClassStore) (*server, error) {
	return &server{store: store}, nil
}

func (s *server) QueryCoursesGET(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	defer r.Body.Close()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	ch := s.store.QueryAll().Evaluate(ctx)
	for class := range ch {
		err := enc.Encode(class)
		fmt.Println(class.Code)
		if err != nil {
			send501(w, err)
			fmt.Println("fail")
			return
		}
	}
}

func (s *server) FlagCoursePOST(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var flaggedCourse struct {
		Code   string `json:"code"`
		Reason string `json:"reason"`
	}

	err := dec.Decode(&flaggedCourse)
	if err != nil {
		send501(w, err)
	}

	//TODO handle courseFlags
}

func (s *server) QueryCoursesPOST(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	//TODO figure out query format
	var q string
	err := dec.Decode(&q)
	if err != nil {
		send501(w, err)
		return
	}
}
