package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	_ "github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	_ "github.com/gocolly/colly/debug"
)

// InstructorData ...
type InstructorData struct {
	Name           string `json:"name"`
	Difficulty     int    `json:"difficulty"`
	Overall        int    `json:"overall"`
	WouldTakeAgain string `json:"would_take_gain"`
}

// How to convert map to json
// https://stackoverflow.com/questions/24652775/convert-go-map-to-json

// Instructor ...
type Instructor map[string]InstructorData

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

func getSubjectPath(subj string) string {
	subjectPrefix := "/cs/courseschedule?pname=subjarea&tname=subj-"
	return subjectPrefix + subj
}

func isAllowedActivity(activity string) bool {
	var allowedSectionActivity = [...]string{
		"Lecture",
		"Lecture-Seminar",
		"Lecture-Laboratory",
		"Distance Education",
		"Laboratory",
		"Practicum",
		"Seminar",
	}

	for _, allowed := range allowedSectionActivity {
		if allowed == activity {
			return true
		}
	}
	return false
}

func writeData(ubcCourseInfo Department) {
	jsonString, err := json.Marshal(ubcCourseInfo)
	err = ioutil.WriteFile("/tmp/dat1.json", jsonString, 0644)
	checkIO(err)
}

func checkIO(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
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
	sectionCollector := c.Clone()

	// =======================
	// All courses page callbacks

	ubcCourseInfo := make(Department)
	coursesURL := "https://courses.students.ubc.ca/cs/courseschedule?pname=subjarea&tname=subj-all-departments"

	departmentPath := getSubjectPath("department")

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
	coursePath := getSubjectPath("course")

	departmentCollector.OnRequest(func(r *colly.Request) {
		curDepartment = getURLParam(r, "dept")
		ubcCourseInfo[curDepartment] = make(Course)
	})

	departmentCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// If course link
		if strings.HasPrefix(e.Attr("href"), coursePath) {
			fmt.Println(curDepartment)
			courseURL := getFullURL(e.Attr("href"))
			courseCollector.Visit(courseURL)
		}
	})

	// =======================
	// courseCollector callbacks

	var curCourse string
	// sectionPath := getSubjectPath("section")

	courseCollector.OnRequest(func(r *colly.Request) {
		curCourse = getURLParam(r, "course")
		ubcCourseInfo[curDepartment][curCourse] = make(Section)
	})

	// Class with that includes "section" (ex. section1, section2, etc)
	courseCollector.OnHTML("tr[class*=section]", func(tr *colly.HTMLElement) {
		sectionLink := tr.ChildAttr("td:nth-child(2) > a", "href")
		sectionActivity := tr.ChildText("td:nth-child(3)")

		if isAllowedActivity(sectionActivity) {
			sectionURL := getFullURL(sectionLink)
			sectionCollector.Visit(sectionURL)
		}
	})

	// =======================
	// sectionCollector callbacks

	var curSection string

	sectionCollector.OnRequest(func(r *colly.Request) {
		curSection = getURLParam(r, "section")
		ubcCourseInfo[curDepartment][curCourse][curSection] = make(Instructor)
	})

	courseCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if strings.HasPrefix(e.Attr("href"), coursePath) {
			writeData(ubcCourseInfo)
		}
	})

	c.Visit(coursesURL)
}
