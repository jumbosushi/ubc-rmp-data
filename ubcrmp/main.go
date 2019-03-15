package ubcrmp

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func getPossibleURLs() []string {
	var result []string
	baseURL := "https://courses.students.ubc.ca/cs/courseschedule?tname=subj-all-departments"
	year := time.Now().Year()
	years := []string{strconv.Itoa(year), strconv.Itoa(year + 1), strconv.Itoa(year - 1)}
	terms := []string{"S", "W"}
	campuses := []string{"UBC", "UBCO"}

	for _, year := range years {
		for _, term := range terms {
			for _, campus := range campuses {
				urlParam := fmt.Sprintf("&sessyr=%s&sesscd=%s&campuscd=%s", year, term, campus)
				testURL := baseURL + urlParam
				resp, err := http.Get(testURL)
				if err != nil {
					log.Fatal(err)
				}
				if resp.StatusCode != 500 {
					result = append(result, testURL)
				}
			}
		}
	}
	return result
}

// Start ...
func Start() {
	fmt.Println("=== ubc-rmp-srcaper start ===")

	validURLs := getPossibleURLs()
	fmt.Printf("\n")
	for _, str := range validURLs {
		fmt.Printf("%s\n", str)
	}

	courseToInstrFileName := "courseToinstrID.json"
	instrToRatingFileName := "instrIDToRating.json"

	for _, allCoursesURL := range validURLs {
		buildCourseJSON(allCoursesURL, courseToInstrFileName, instrToRatingFileName)
		QueryRMP(allCoursesURL, instrToRatingFileName)
	}

	fmt.Println("=== ubc-rmp-srcaper end   ===")
}
