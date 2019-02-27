package ubcrmp

// InstructorData ...
type InstructorData struct {
	UbcID          int     `json:"ubcid"`
	RmpID          int     `json:"rmpid"`
	Name           string  `json:"name"`
	Difficulty     float64 `json:"difficulty"`
	Overall        float64 `json:"overall"`
	WouldTakeAgain string  `json:"would_take_again"`
}

// Instructor ...
type Instructor map[int]InstructorData

// Section ...
type Section map[string]int

// Course ...
type Course map[string]Section

// Department ...
type Department map[string]Course
