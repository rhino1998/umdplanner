package duration

import "time"

type Duration struct {
	Start time.Time
	End   time.Time
}

type Conflicter interface {
	Conflicts(Conflicter) bool
}

func (d Duration) Conflicts(o Conflicter) bool {
	do, ok := o.(*Duration)
	if ok {
		return (d.Start.Before(do.Start) && d.End.After(do.Start)) ||
			(do.Start.Before(d.Start) && do.End.After(d.Start))
	}
	return o.Conflicts(d)
}
