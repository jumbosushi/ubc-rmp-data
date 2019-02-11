package ubcrmp

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func getRmpQuery(name string, ubcID int) string {
	baseRmpQuery := "https://www.ratemyprofessors.com/search.jsp?query=university+of+british+columbia+"
	// Add name
	// Query strings can't include ", " or "  "
	noComma := strings.Replace(name, ",", "", -1)
	noSpace := strings.Replace(noComma, " ", "+", -1)

	// Add ubcid to be passed to Colly
	baseUbcIDParam := "&ubcid="
	ubcIDStr := strconv.Itoa(ubcID)
	ubcIDParam := baseUbcIDParam + ubcIDStr
	return baseRmpQuery + noSpace + ubcIDParam
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

func getTID(path string) int {
	parts := strings.Split(path, "tid=")
	tid, err := strconv.Atoi(string(parts[1]))
	if err != nil {
		log.Fatal(err)
	}
	return tid
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

	// =======================
	// rmpSearchCollector callbacks

	rmpSearchCollector := c.Clone()
	// rmpStatsCollector := c.Clone()

	var curInstrUbcID int

	rmpSearchCollector.OnRequest(func(r *colly.Request) {
		q := r.URL.Query()
		curInstrUbcID, _ = strconv.Atoi(q.Get("ubcid"))
	})

	rmpSearchCollector.OnHTML("li.PROFESSOR > a[href]", func(e *colly.HTMLElement) {
		// If instr link exits, get tid (RmpID)
		instrLink := e.Attr("href")
		tid := getTID(instrLink)
		// Update RmpID
		rmpInstr := instrMap[curInstrUbcID]
		rmpInstr.RmpID = tid
		log.Println(curInstrUbcID)
		log.Println(rmpInstr.RmpID)
		instrMap[curInstrUbcID] = rmpInstr
	})

	for _, instrData := range instrMap {
		rmpQuery := getRmpQuery(instrData.Name, instrData.UbcID)
		log.Println(rmpQuery)
		rmpSearchCollector.Visit(rmpQuery)
		rmpSearchCollector.Wait()
	}
}
