package ubcrmp

import (
	"fmt"
)

// Start ...
func Start() {
	fmt.Println("=== ubc-rmp-srcaper start ===")

	courseToInstrFileName := "courseToinstrID.json"
	instrToRatingFileName := "instrIDToRating.json"

	buildCourseJSON(courseToInstrFileName, instrToRatingFileName)
	QueryRMP(instrToRatingFileName)

	fmt.Println("=== ubc-rmp-srcaper end   ===")
}
