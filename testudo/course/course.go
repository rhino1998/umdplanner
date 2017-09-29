package course

import (
	"regexp"
	"time"
)

var MatchCode = regexp.MustCompile("[a-zA-Z]{3,4}[0-9]{3}[a-zA-Z]?")

type Class struct {
	Title   string   `json:"title"`
	Code    string   `json:"code"`
	Credits int      `json:"credits"`
	Prereqs []*Class `json:"prereqs"`

	Description  string `json:"description"`
	Prerequisite string `json:"prerequisite"`
	Restriction  string `json:"restriction"`

	GenEd    GenEd     `json:"gen_ed"`
	Sections []Section `json:"sections"`
}

func (c *Class) HasGenEd(ge GenEd) bool {
	return c.GenEd&ge != 0
}

func (c *Class) Conflicts(o *Class) bool {
	if c == nil || o == nil {
		return false
	}
	for _, cS := range c.Sections {
		for _, oS := range o.Sections {
			durations := make([]Duration, len(oS.Times))
			for i, t := range oS.Times {
				durations[i] = t.Duration
			}
			if cS.Conflicts(durations...) {
				return true
			}
		}
	}
	return false

}

type Section struct {
	Code      string `json:"code"`
	Times     []Time `json:"times"`
	Professor string `json:"professor"`
}

func (s *Section) Conflicts(ds ...Duration) bool {
	for _, sT := range s.Times {
		for _, d := range ds {

			if (sT.Start.Before(d.Start) && sT.End.After(d.Start)) ||
				(d.Start.Before(sT.Start) && d.End.After(sT.Start)) {
				return true
			}
		}
	}
	return false
}

type Time struct {
	Duration
	Room string
}

type Duration struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
