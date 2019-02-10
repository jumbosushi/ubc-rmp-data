package ubcrmp

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func getRmpQuery(name string) string {
	// Query strings can't include ", " or "  "
	noComma := strings.Replace(name, ",", "", -1)
	noSpace := strings.Replace(noComma, " ", "+", -1)

	origRmpQuery := "https://www.ratemyprofessors.com/search.jsp?query="
	return origRmpQuery + noSpace
}

func readJSON() (Department, Instructor) {
	courseJSON, _ := ioutil.ReadFile("data/ubcrmpCourse.json")
	courseData := make(Department)
	err := json.Unmarshal(courseJSON, &courseData)
	if err != nil {
		log.Fatal(err)
	}

	instrJSON, _ := ioutil.ReadFile("data/ubcrmpInstr.json")
	instrData := make(Instructor)
	err = json.Unmarshal(instrJSON, &instrData)
	if err != nil {
		log.Fatal(err)
	}
	return courseData, instrData
}

// QueryRMP ..
func QueryRMP() {
	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("UBC-RMP Bot"),
	)

	c.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob: "www.ratemyprofessors.com/*",
		// Set delay between requests
		Delay: 1 * time.Second,
		// Add additional random delay
		RandomDelay: 1 * time.Second,
		Parallelism: 2,
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	_, instrMap := readJSON()

	for _, instrData := range instrMap {
		rmpQuery := getRmpQuery(instrData.Name)
		c.Visit(rmpQuery)
	}
	c.Wait()
}
