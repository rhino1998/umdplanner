package section

import (
	"context"

	"github.com/rhino1998/umdplanner/testudo/duration"
)

type Query interface {
	Evaluate(context.Context) <-chan *Section
}

func WithExcludedTimes(q Query, ds ...duration.Duration) Query {
	return &excludeTimes{parent: q, durations: ds}
}

type excludeTimes struct {
	parent    Query
	durations []duration.Duration
}

func (q *excludeTimes) Evaluate(ctx context.Context) <-chan *Section {
	out := make(chan *Section)
	ch := q.parent.Evaluate(ctx)

	go func() {
		for s := range ch {
			conflict := false
			for _, d := range q.durations {
				if s.Conflicts(d) {
					conflict = true
				}
			}
			if !conflict {
				select {
				case out <- s:
				case <-ctx.Done():
					break
				}
			}
		}
		close(out)
	}()
	return out
}
