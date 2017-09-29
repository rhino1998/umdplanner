package testudo

func QueryWithGenEd(q Query, ge GenEd) Query {

	return &queryGenEd{parent: q, gened: ge}
}

func QueryWithExcludedTimes(q Query, ds ...Duration) Query {
	return &queryExcludeTimes{parent: q, durations: ds}
}

type Query interface {
	Evaluate() <-chan *Class
}

type queryGenEd struct {
	parent Query
	gened  GenEd
}

func (q *queryGenEd) Evaluate() <-chan *Class {
	out := make(chan *Class)
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

type queryExcludeTimes struct {
	parent    Query
	durations []Duration
}

func (q *queryExcludeTimes) Evaluate() <-chan *Class {
	out := make(chan *Class)
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
