package query

import "github.com/rhino1998/umdplanner/testudo/course"

type genEd struct {
	parent Query
	gened  course.GenEd
}

type Query interface {
	Evaluate() <-chan *course.Class
}

func (q *genEd) Evaluate() <-chan *course.Class {
	out := make(chan *course.Class)
	go func() {
		for class := range q.parent.Evaluate() {
			if class.HasGenEd(q.gened) {
				out <- class
			}
		}
		close(out)
	}()
	return out
}

type excludeTimes struct {
	parent    Query
	durations []course.Duration
}

func (q *excludeTimes) Evaluate() <-chan *course.Class {
	out := make(chan *course.Class)
	go func() {
		for class := range q.parent.Evaluate() {
			conflict := false
			for _, section := range class.Sections {
				if section.Conflicts(q.durations...) {
					conflict = true
				}
			}
			if !conflict {
				out <- class
			}
		}
		close(out)
	}()
	return out
}

func WithGenEd(q Query, ge course.GenEd) Query {
	return &genEd{parent: q, gened: ge}
}

func WithExcludedTimes(q Query, ds ...course.Duration) Query {
	return &excludeTimes{parent: q, durations: ds}
}
