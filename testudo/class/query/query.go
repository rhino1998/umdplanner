package query

import (
	"context"

	"github.com/rhino1998/umdplanner/testudo/course"
)

type genEd struct {
	parent Query
	gened  course.GenEd
}

type Query interface {
	Evaluate(context.Context) <-chan *course.Class
}

func (q *genEd) Evaluate(ctx context.Context) <-chan *course.Class {
	out := make(chan *course.Class)
	ch := q.parent.Evaluate(ctx)

	go func() {
		for class := range ch {
			if class.HasGenEd(q.gened) {
				select {
				case out <- class:
				case <-ctx.Done():
					break
				}
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

func (q *excludeTimes) Evaluate(ctx context.Context) <-chan *course.Class {
	out := make(chan *course.Class)
	ch := q.parent.Evaluate(ctx)

	go func() {
		for class := range ch {
			conflict := false
			for _, section := range class.Sections {
				if section.Conflicts(q.durations...) {
					conflict = true
				}
			}
			if !conflict {
				select {
				case out <- class:
				case <-ctx.Done():
					break
				}
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
