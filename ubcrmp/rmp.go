package ubcrmp

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
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

func readJSON(instrToRatingFileName string) Instructor {
	instrJSON, _ := ioutil.ReadFile("data/" + instrToRatingFileName)
	instrData := make(Instructor)
	err := json.Unmarshal(instrJSON, &instrData)
	if err != nil {
		log.Fatal(err)
	}
	return instrData
}

func getTID(path string) int {
	parts := strings.Split(path, "tid=")
	tid, err := strconv.Atoi(string(parts[1]))
	if err != nil {
		log.Fatal(err)
	}
	return tid
}

func getFullRmpURL(path string) string {
	coursesPrefix := "https://www.ratemyprofessors.com"
	return coursesPrefix + path
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

// QueryRMP ..
func QueryRMP(allCoursesURL string, instrToRatingFileName string) {
	c := colly.NewCollector(
		// colly.Async(true),
		colly.UserAgent("UBC-RMP Bot"),
	)

	c.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainGlob: "www.ratemyprofessors.com/*",
		// Set delay between requests
		Delay: 1 * time.Second,
		// Add additional random delay
		RandomDelay: 1 * time.Second,
		// Parallelism: 2,
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	instrToRatingFileName = getTermFileName(allCoursesURL, instrToRatingFileName)
	instrMap := readJSON(instrToRatingFileName)

	rmpSearchCollector := c.Clone()
	rmpStatsCollector := c.Clone()

	// =======================
	// rmpSearchCollector callbacks

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
		instrMap[curInstrUbcID] = rmpInstr
		// Request stats page
		fullRmpURL := getFullRmpURL(instrLink)
		log.Println(fullRmpURL)
		rmpStatsCollector.Visit(fullRmpURL)
	})

	// =======================
	// rmpStatsCollector callbacks

	rmpStatsCollector.OnHTML(".left-breakdown", func(e *colly.HTMLElement) {
		quality := e.ChildText(".quality > div > .grade")
		difficulty := e.ChildText(".difficulty > .grade")
		wouldTakeAgain := e.ChildText(".takeAgain > .grade")
		// Update values
		rmpInstr := instrMap[curInstrUbcID]
		tempQ, _ := strconv.ParseFloat(quality, 32)
		rmpInstr.Overall = toFixed(tempQ, 1)
		tempD, _ := strconv.ParseFloat(difficulty, 32)
		rmpInstr.Difficulty = toFixed(tempD, 1)
		rmpInstr.WouldTakeAgain = wouldTakeAgain
		instrMap[curInstrUbcID] = rmpInstr
		log.Println(rmpInstr)
	})

	for _, instrData := range instrMap {
		rmpQuery := getRmpQuery(instrData.Name, instrData.UbcID)
		rmpSearchCollector.Visit(rmpQuery)
		writeJSON(instrMap, instrToRatingFileName)
	}
}
