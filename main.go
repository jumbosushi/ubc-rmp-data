package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	_ "github.com/gocolly/colly/debug"
)

var allowedCourseActivity = [...]string{
	"Lecture",
	"Lecture-Seminar",
	"Lecture-Laboratory",
	"Distance Education",
	"Laboratory",
	"Practicum",
	"Seminar",
}

func main() {
	coursesURL := "https://courses.students.ubc.ca/cs/courseschedule?pname=subjarea&tname=subj-all-departments"

	c := colly.NewCollector(
		// Allow crawling to be done in parallel / async
		// colly.Async(true),
		colly.UserAgent("UBC-RMP Bot"),
		// Attach a debugger to the collector
		// colly.Debugger(&debug.LogDebugger{}),
	)

	c.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob: "courses.students.ubc.ca/*",
		// Set delay between requests
		Delay: 1 * time.Second,
		// Add additional random delay
		RandomDelay: 1 * time.Second,
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	// Create another collectors
	// departmentCollector := c.Clone()
	// courseCollector := c.Clone()
	// sectionCollector := c.Clone()

	departmentURL := "/cs/courseschedule?pname=subjarea&tname=subj-department&dept="
	c.OnHTML("a", func(e *colly.HTMLElement) {
		if strings.HasPrefix(e.Attr("href"), departmentURL) {
			fmt.Println(e.Text)
			fmt.Println(e.Attr("href"))
		}
	})

	c.Visit(coursesURL)
}
