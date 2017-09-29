package testudo

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/rhino1998/umdplanner/testudo/course"
	"github.com/rhino1998/umdplanner/testudo/query"
)

//NewStore makes a new map-based store for classes
func NewStore() ClassStore {
	return &mapStore{classes: make(map[string]*course.Class)}
}

//LoadStore loads the store from reader as json
func LoadStore(r io.Reader) (ClassStore, error) {
	dec := json.NewDecoder(r)
	var classes map[string]*course.Class

	err := dec.Decode(&classes)

	cs := &mapStore{classes: classes}
	linkClasses(cs)
	return cs, err
}

type mapStore struct {
	lock    sync.RWMutex
	classes map[string]*course.Class
}

func (cs *mapStore) Dump(w io.Writer) error {
	cs.lock.RLock()
	enc := json.NewEncoder(w)
	err := enc.Encode(cs.classes)

	cs.lock.RUnlock()

	return err
}

func (cs *mapStore) Get(code string) (*course.Class, error) {
	cs.lock.RLock()
	class, ok := cs.classes[code]
	cs.lock.RUnlock()
	var err error
	if !ok {
		err = fmt.Errorf("Class not found: %q", code)
	}
	return class, err
}

func (cs *mapStore) Set(c *course.Class) error {
	cs.lock.Lock()
	cs.classes[c.Code] = c
	cs.lock.Unlock()
	return nil
}

func (cs *mapStore) QueryAll() query.Query {
	return &allQuery{
		eval: func() <-chan *course.Class {
			out := make(chan *course.Class)
			go func() {
				cs.lock.RLock()
				for _, class := range cs.classes {
					out <- class
				}
				cs.lock.RUnlock()
				close(out)
			}()
			return out
		},
	}
}

type allQuery struct {
	eval func() <-chan *course.Class
}

func (q *allQuery) Evaluate() <-chan *course.Class {
	return q.eval()
}
