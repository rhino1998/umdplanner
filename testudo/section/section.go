package section

import "github.com/rhino1998/umdplanner/testudo/duration"

type Section struct {
	Code      string     `json:"code"`
	Meetings  []*Meeting `json:"meetings"`
	Professor string     `json:"professor"`
}

func (s *Section) Conflicts(o duration.Conflicter) bool {
	for _, sM := range s.Meetings {
		if o.Conflicts(sM) {
			return true
		}
	}
	return false
}

type Meeting struct {
	Room     string
	Building string
	duration.Duration
}
