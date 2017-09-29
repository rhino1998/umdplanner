package testudo

import (
	"encoding/json"
	"io"
	"sync"
)

func NewStore() ClassStore {
	return &sliceStore{classes: make([]*Class, 0, 0)}
}

func LoadStore(r io.Reader) (ClassStore, error) {
	dec := json.NewDecoder(r)
	var store []*Class

	err := dec.Decode(&store)

	return &sliceStore{classes: store}, err
}

type sliceStore struct {
	lock    sync.RWMutex
	classes []*Class
}

func (cs *sliceStore) Dump(w io.Writer) error {
	cs.lock.RLock()
	enc := json.NewEncoder(w)
	err := enc.Encode(cs.classes)

	cs.lock.RUnlock()

	return err
}

func (cs *sliceStore) Store(c *Class) error {
	cs.lock.Lock()
	cs.classes = append(cs.classes, c)
	cs.lock.Unlock()
	return nil
}

func (cs *sliceStore) QueryAll() Query {
	return &allQuery{
		eval: func() <-chan *Class {
			out := make(chan *Class)
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
	eval func() <-chan *Class
}

func (q *allQuery) Evaluate() <-chan *Class {
	return q.eval()
}
