package ubcrmp

import (
	"fmt"
)

// Start ...
func Start() {
	fmt.Println("=== ubc-rmp-srcaper start ===")
	// buildCourseJSON()
	QueryRMP()
	fmt.Println("=== ubc-rmp-srcaper end   ===")
}
