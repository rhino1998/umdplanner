package testudo

import (
	"context"
	"fmt"
	"io"
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rhino1998/umdplanner/testudo/course"
	"github.com/rhino1998/umdplanner/testudo/query"
)

const url = "https://ntst.umd.edu/soc/"

type ClassStore interface {
	Set(*course.Class) error
	Get(string) (*course.Class, error)
	QueryAll() query.Query
	Dump(io.Writer) error
}

//ScrapeAll scrapes all the classes from testudo.umd.edu
func ScrapeAll(url string, cs ClassStore) error {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	sema := make(chan struct{}, runtime.NumCPU())
	doc.Find(".course-prefix a.clearfix").Each(func(i int, s *goquery.Selection) {
		part, _ := s.Attr("href")
		wg.Add(1)
		sema <- struct{}{}
		go func() {
			ScrapeDepartment(url+"/"+part, cs)
			<-sema
			wg.Done()
		}()
	})
	wg.Wait()
	return nil
}

//ScrapeDepartment scrapes a whole department list of classes from
//testudo.umd.edu
func ScrapeDepartment(url string, cs ClassStore) error {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	doc.Find(".courses-container div.course").Each(func(_ int, s *goquery.Selection) {
		code, ok := s.Attr("id")
		if !ok {
			return
		}
		wg.Add(1)
		go func() {
			class, err := ScrapeClass(url + "/" + code)
			if err != nil {
				log.Println(err)
			}
			cs.Set(class)
			wg.Done()
		}()
	})
	wg.Wait()

	return nil

}

//ScrapeClass scrapes a class from testudo.umd.edu schedule of classes
func ScrapeClass(url string) (*course.Class, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}
	s := doc.Find(".courses-container div.course").First()
	code, ok := s.Attr("id")
	if !ok {
		return nil, fmt.Errorf("No class code found: %q", url)
	}

	fmt.Println(code)

	title := s.Find(".course-title").First().Text()
	credits, err := strconv.Atoi(s.Find(".course-min-credits").First().Text())
	if err != nil {
		log.Printf("%s: %q\n", code, err)
	}

	genedFields := s.Find(".course-subcategory")
	genedCodes := make([]string, genedFields.Length())

	i := 0
	genedFields.Each(func(_ int, s *goquery.Selection) {
		genedCodes[i] = strings.TrimSpace(s.Text())
		i++
	})
	geneds := course.ParseGenEd(genedCodes)

	sectionFields := s.Find(".section")
	sections := make([]course.Section, sectionFields.Length())
	i = 0
	sectionFields.Each(func(_ int, s *goquery.Selection) {
		timesFields := s.Find(".row")
		times := make([]course.Time, 0, timesFields.Length()*7)

		j := 0
		timesFields.Each(func(_ int, s *goquery.Selection) {
			days := parseDays(s.Find(".section-day-time-group span.section-days").Text())
			if len(days) == 0 {
				return
			}

			start, err := time.Parse("3:04pm", s.Find(".class-start-time").Text())
			if err != nil {
				fmt.Println(err)
				return
			}
			end, err := time.Parse("3:04pm", s.Find(".class-end-time").Text())
			if err != nil {
				return
			}

			for _, day := range days {
				times = append(times, course.Time{
					Room: fmt.Sprintf(
						"%s %s",
						s.Find(".building-code").Text(),
						s.Find(".class-room").Text(),
					),
					Duration: course.Duration{
						Start: start.AddDate(0, 0, int(day)),
						End:   end.AddDate(0, 0, int(day)),
					},
				})
				j++
			}
		})

		sections[i] = course.Section{
			Times:     times,
			Code:      strings.TrimSpace(s.Find(".section-id").Text()),
			Professor: s.Find(".section-instructor").First().Text(),
		}
		i++
	})

	class := &course.Class{
		Code:        code,
		Title:       title,
		Credits:     credits,
		GenEd:       geneds,
		Prereqs:     []*course.Class{},
		Description: s.Find(".approved-course-text").Last().Text(),
		Prerequisite: strings.Replace(s.Find(".approved-course-text div strong").
			FilterFunction(func(_ int, s *goquery.Selection) bool {
				return s.Text() == "Prerequisite:"
			}).Parent().First().Text(), "Prerequisite: ", "", -1),
		Restriction: strings.Replace(s.Find(".approved-course-text div strong").
			FilterFunction(func(_ int, s *goquery.Selection) bool {
				return s.Text() == "Restriction:"
			}).Parent().First().Text(), "Restriction: ", "", -1),
		Sections: sections,
	}

	return class, nil

}

func linkClasses(store ClassStore) {
	ch := store.QueryAll().Evaluate(context.Background())
	for class := range ch {
		reqs := course.MatchCode.FindAllString(class.Prerequisite, -1)
		for _, req := range reqs {
			oClass, err := store.Get(req)
			if err != nil {
				continue
			}
			class.Prereqs = append(class.Prereqs, oClass)
		}
	}
}
