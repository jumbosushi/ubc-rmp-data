package model

// InstructorData ...
type InstructorData struct {
	Dept           string `json:"dept"`
	Lecture        string `json:"lecture"`
	Name           string `json:"name"`
	Difficulty     int    `json:"difficulty"`
	Overall        int    `json:"overall"`
	WouldTakeAgain string `json:"would_take_gain"`
}

// Instructor ...
type Instructor map[string]InstructorData

// Section ...
type Section map[string]string

// Course ...
type Course map[string]Section

// Department ...
type Department map[string]Course
