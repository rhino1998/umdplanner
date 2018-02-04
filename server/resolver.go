package main

import (
	"context"
	"unsafe"

	graphql "github.com/neelance/graphql-go"
	"github.com/rhino1998/umdplanner/testudo"
	"github.com/rhino1998/umdplanner/testudo/class"
	"github.com/rhino1998/umdplanner/testudo/duration"
	"github.com/rhino1998/umdplanner/testudo/section"
)

type Resolver struct {
	store testudo.ClassStore
}

func (r *Resolver) Classes(ctx context.Context, args struct {
	MinCredits *int32
	MaxCredits *int32
	GenEds     *[]string
}) ([]*classResolver, error) {
	var l []*classResolver
	q := r.store.QueryAll()
	if args.GenEds != nil && len(*args.GenEds) > 0 {
		q = class.WithGenEd(q, class.ParseGenEd(*args.GenEds))
	}
	if args.MinCredits != nil {
		q = class.WithMinCredits(q, int(*args.MinCredits))
	}
	if args.MaxCredits != nil {
		q = class.WithMaxCredits(q, int(*args.MaxCredits))
	}
	for c := range q.Evaluate(ctx) {
		l = append(l, &classResolver{c})
	}
	return l, nil

}

func (r *Resolver) Class(ctx context.Context, args struct{ Code string }) (*classResolver, error) {
	c, err := r.store.Get(args.Code)
	return &classResolver{c}, err
}

type classResolver struct {
	class *class.Class
}

func (c *classResolver) Code() string {
	return c.class.Code
}

func (c *classResolver) Title() string {
	return c.class.Title
}

func (c *classResolver) Description() *string {
	if c.class.Description == "" {
		return nil
	}
	return &c.class.Description
}

func (c *classResolver) Prerequisite() *string {
	if c.class.Prerequisite == "" {
		return nil
	}
	return &c.class.Prerequisite
}

func (c *classResolver) Restriction() *string {
	if c.class.Restriction == "" {
		return nil
	}
	return &c.class.Restriction
}

func (c *classResolver) GenEds() []string {
	var genEds []string
	ge := c.class.GenEd
	for i := 0; i < int(unsafe.Sizeof(ge)*8); i++ {
		s := (class.GenEd(ge & (1 << uint(i)))).String()
		if s != "" {
			genEds = append(genEds, s)
		}
	}
	return genEds
}

func (c *classResolver) Credits() int32 {
	return int32(c.class.Credits)
}

func (c *classResolver) Sections(ctx context.Context) []*sectionResolver {
	var l []*sectionResolver
	q := c.class.QueryAll()
	for s := range q.Evaluate(ctx) {
		l = append(l, &sectionResolver{s})
	}
	return l

}

type sectionResolver struct {
	section *section.Section
}

func (s *sectionResolver) Code() string {
	return s.section.Code
}

func (s *sectionResolver) Meetings() []*meetingResolver {
	var l []*meetingResolver
	for _, m := range s.section.Meetings {
		l = append(l, &meetingResolver{m})
	}
	return l
}

type meetingResolver struct {
	meeting *section.Meeting
}

func (m *meetingResolver) Building() string {
	return m.meeting.Building
}

func (m *meetingResolver) Room() string {
	return m.meeting.Room
}

func (m *meetingResolver) Duration() *durationResolver {
	return &durationResolver{m.meeting.Duration}
}

type durationResolver struct {
	duration duration.Duration
}

func (d *durationResolver) Start() graphql.Time {
	return graphql.Time{
		Time: d.duration.Start,
	}
}
func (d *durationResolver) End() graphql.Time {
	return graphql.Time{
		Time: d.duration.End,
	}
}
