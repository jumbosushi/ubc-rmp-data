package main

import (
	_ "encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	_ "github.com/gocolly/colly/debug"
)

// Instructor ...
type Instructor struct {
	Name           string `json:"name"`
	Difficulty     int    `json:"difficulty"`
	Overall        int    `json:"overall"`
	WouldTakeAgain string `json:"would_take_gain"`
}

// How to convert map to json
// https://stackoverflow.com/questions/24652775/convert-go-map-to-json

// Section ...
type Section map[string]Instructor

// Course ...
type Course map[string]Section

// Department ...
type Department map[string]Course

func getURLParam(r *colly.Request, param string) string {
	URL := r.URL.String()
	URLSplit := strings.Split(URL, param+"=")
	return URLSplit[1]
}

func getFullURL(path string) string {
	coursesPrefix := "https://courses.students.ubc.ca"
	return coursesPrefix + path
}

func main() {
	// cpsc := make(Department)
	// cpsc["CPSC"] = make(Course)
	// cpsc["CPSC"]["110"] = make(Section)
	// cpsc["CPSC"]["110"]["001"] = Instructor{}

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
	departmentCollector := c.Clone()
	courseCollector := c.Clone()
	// sectionCollector := c.Clone()

	// =======================
	// All courses page callbacks

	ubcCourseInfo := make(Department)
	coursesURL := "https://courses.students.ubc.ca/cs/courseschedule?pname=subjarea&tname=subj-all-departments"

	departmentPath := "/cs/courseschedule?pname=subjarea&tname=subj-department"
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// If department link
		if strings.HasPrefix(e.Attr("href"), departmentPath) {
			deptURL := getFullURL(e.Attr("href"))
			departmentCollector.Visit(deptURL)
		}
	})

	// =======================
	// departmentCollector callbacks

	var curDepartment string
	coursePath := "/cs/courseschedule?pname=subjarea&tname=subj-course"

	departmentCollector.OnRequest(func(r *colly.Request) {
		curDepartment = getURLParam(r, "dept")
	})

	departmentCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// If course link
		if strings.HasPrefix(e.Attr("href"), coursePath) {
			fmt.Println(curDepartment)
			ubcCourseInfo[curDepartment] = make(Course)
			courseURL := getFullURL(e.Attr("href"))
			courseCollector.Visit(courseURL)
		}
	})

	// =======================
	// courseCollector callbacks

	var curCourse string

	courseCollector.OnRequest(func(r *colly.Request) {
		curCourse = getURLParam(r, "course")
	})

	courseCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// If course link
		ubcCourseInfo[curDepartment][curCourse] = make(Section)
		if strings.HasPrefix(e.Attr("href"), coursePath) {
			fmt.Println(curCourse)
		}
	})

	// =======================
	// sectionCollector callbacks

	// var curSection string

	// sectionCollector.OnRequest(func(r *colly.Request) {
	// curSection = getURLParam(r, "dept")
	// })

	// courseCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// if strings.HasPrefix(e.Attr("href"), coursePath) {
	// fmt.Println(curDepartment)
	// }
	// })

	// var allowedCourseActivity = [...]string{
	// "Lecture",
	// "Lecture-Seminar",
	// "Lecture-Laboratory",
	// "Distance Education",
	// "Laboratory",
	// "Practicum",
	// "Seminar",
	// }

	c.Visit(coursesURL)

}
