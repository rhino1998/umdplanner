package class

import (
	"context"

	"github.com/rhino1998/umdplanner/testudo/duration"
)

type Query interface {
	Evaluate(context.Context) <-chan *Class
}

func WithGenEd(q Query, ge GenEd) Query {
	return &genEd{parent: q, gened: ge}
}

func WithMinCredits(q Query, c int) Query {
	return &minCredit{parent: q, min: c}
}

func WithMaxCredits(q Query, c int) Query {
	return &maxCredit{parent: q, max: c}
}

func WithExcludedTimes(q Query, ds ...duration.Duration) Query {
	return &excludeTimes{parent: q, durations: ds}
}

type genEd struct {
	parent Query
	gened  GenEd
}

func (q *genEd) Evaluate(ctx context.Context) <-chan *Class {
	out := make(chan *Class)
	ch := q.parent.Evaluate(ctx)

	go func() {
		for c := range ch {
			if c.HasGenEd(q.gened) {
				select {
				case out <- c:
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
	durations []duration.Duration
}

func (q *excludeTimes) Evaluate(ctx context.Context) <-chan *Class {
	out := make(chan *Class)
	ch := q.parent.Evaluate(ctx)

	go func() {
		for c := range ch {
			conflict := false
			for _, d := range q.durations {
				if c.Conflicts(d) {
					conflict = true
				}
			}
			if !conflict {
				select {
				case out <- c:
				case <-ctx.Done():
					break
				}
			}
		}
		close(out)
	}()
	return out
}

type minCredit struct {
	parent Query
	min    int
}

func (q *minCredit) Evaluate(ctx context.Context) <-chan *Class {
	out := make(chan *Class)
	ch := q.parent.Evaluate(ctx)

	go func() {
		for c := range ch {
			if c.Credits >= q.min {
				select {
				case out <- c:
				case <-ctx.Done():
					break
				}
			}
		}
		close(out)
	}()
	return out
}

type maxCredit struct {
	parent Query
	max    int
}

func (q *maxCredit) Evaluate(ctx context.Context) <-chan *Class {
	out := make(chan *Class)
	ch := q.parent.Evaluate(ctx)

	go func() {
		for c := range ch {
			if c.Credits <= q.max {
				select {
				case out <- c:
				case <-ctx.Done():
					break
				}
			}
		}
		close(out)
	}()
	return out
}
