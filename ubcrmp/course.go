package ubcrmp

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// How to convert map to json
// https://stackoverflow.com/questions/24652775/convert-go-map-to-json

func getInstrID(path string) int {
	URL := getFullURL(path)
	u, err := url.Parse(URL)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	ubcidStr := q.Get("ubcid")
	ubcidInt, err := strconv.Atoi(ubcidStr)
	return ubcidInt
}

func getURLParam(r *colly.Request, param string) string {
	q := r.URL.Query()
	return q.Get(param)
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

func buildCourseJSON() {
	c := colly.NewCollector(
		colly.UserAgent("UBC-RMP Bot"),
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

	ubcInstrInfo := make(Instructor)
	ubcCourseInfo := make(Department)
	allCoursesURL := "https://courses.students.ubc.ca/cs/courseschedule?pname=subjarea&tname=subj-all-departments"

	departmentPath := getSubjectPath("department")

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// If department link
		deptLink := e.Attr("href")
		if strings.HasPrefix(deptLink, departmentPath) {
			deptURL := getFullURL(deptLink)
			departmentCollector.Visit(deptURL)
		}
	})

	// =======================
	// departmentCollector callbacks

	var curDepartment string
	coursePath := getSubjectPath("course")

	departmentCollector.OnRequest(func(r *colly.Request) {
		curDepartment = getURLParam(r, "dept")
		fmt.Printf("%s\n", curDepartment)
		ubcCourseInfo[curDepartment] = make(Course)
	})

	departmentCollector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// If course link
		courseLink := e.Attr("href")
		if strings.HasPrefix(courseLink, coursePath) {
			courseURL := getFullURL(courseLink)
			courseCollector.Visit(courseURL)
		}
	})

	// =======================
	// courseCollector callbacks

	var curCourse string
	sectionPath := getSubjectPath("section")

	courseCollector.OnRequest(func(r *colly.Request) {
		curCourse = getURLParam(r, "course")
		fmt.Printf("  %s\n", curCourse)
		ubcCourseInfo[curDepartment][curCourse] = make(Section)
	})

	// Class with that includes "section" (ex. section1, section2, etc)
	courseCollector.OnHTML("tr[class*=section]", func(tr *colly.HTMLElement) {
		sectionLink := tr.ChildAttr("td:nth-child(2) > a", "href")
		sectionActivity := tr.ChildText("td:nth-child(3)")

		if isAllowedActivity(sectionActivity) &&
			strings.HasPrefix(sectionLink, sectionPath) {
			sectionURL := getFullURL(sectionLink)
			sectionCollector.Visit(sectionURL)
		}
	})

	// =======================
	// sectionCollector callbacks

	var curSection string
	instrPath := "/cs/courseschedule?pname=inst"

	sectionCollector.OnRequest(func(r *colly.Request) {
		curSection = getURLParam(r, "section")
		fmt.Printf("    %s\n", curSection)
		// ubcCourseInfo[curDepartment][curCourse][curSection] = make(model.Instructor)
	})

	sectionCollector.OnHTML("td > a[href]", func(e *colly.HTMLElement) {
		if strings.HasPrefix(e.Attr("href"), instrPath) {
			instrName := e.Text
			instrUbcID := getInstrID(e.Attr("href"))

			ubcCourseInfo[curDepartment][curCourse][curSection] = instrUbcID
			writeJSON(ubcCourseInfo, "ubcrmpCourse.json")

			// If record already exists, skip
			if _, ok := ubcInstrInfo[instrUbcID]; ok {
				return
			}

			// Init InstructorData
			instrData := InstructorData{}
			instrData.Name = instrName
			instrData.UbcID = instrUbcID

			ubcInstrInfo[instrUbcID] = instrData
			writeJSON(ubcInstrInfo, "ubcrmpInstr.json")
		}
	})

	c.Visit(allCoursesURL)
}
