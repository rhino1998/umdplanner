package class

import (
	"context"
	"regexp"

	"github.com/rhino1998/umdplanner/testudo/duration"
	"github.com/rhino1998/umdplanner/testudo/section"
)

var MatchCode = regexp.MustCompile("[a-zA-Z]{3,4}[0-9]{3}[a-zA-Z]?")

type Class struct {
	Title   string   `json:"title"`
	Code    string   `json:"code"`
	Credits int      `json:"credits"`
	Prereqs []*Class `json:"-"`

	Description  string `json:"description"`
	Prerequisite string `json:"prerequisite"`
	Restriction  string `json:"restriction"`

	GenEd    GenEd              `json:"gen_ed"`
	Sections []*section.Section `json:"sections"`
}

func (c *Class) HasGenEd(ge GenEd) bool {
	return c.GenEd&ge != 0
}

func (c *Class) QueryAll() section.Query {
	return &sectionQuery{
		eval: func(ctx context.Context) <-chan *section.Section {
			out := make(chan *section.Section)
			go func() {
				for _, s := range c.Sections {
					select {
					case out <- s:
					case <-ctx.Done():
						break
					}
				}
				close(out)
			}()
			return out
		},
	}
}

type sectionQuery struct {
	eval func(context.Context) <-chan *section.Section
}

func (q *sectionQuery) Evaluate(ctx context.Context) <-chan *section.Section {
	return q.eval(ctx)
}

func (c *Class) Conflicts(o duration.Conflicter) bool {
	if c == nil || o == nil {
		return false
	}
	switch oT := o.(type) {
	case *Class:
		for _, oS := range oT.Sections {
			if !c.Conflicts(oS) {
				return false
			}
		}
	default:
		for _, cS := range c.Sections {
			if !o.Conflicts(cS) {
				return false
			}
		}
	}
	return true

}
