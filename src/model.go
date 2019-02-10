package main

// InstructorData ...
type InstructorData struct {
	UbcID          int    `json:"ubcid"`
	Name           string `json:"name"`
	Difficulty     int    `json:"difficulty"`
	Overall        int    `json:"overall"`
	WouldTakeAgain string `json:"would_take_gain"`
}

// Instructor ...
type Instructor map[int]InstructorData

// Section ...
type Section map[string]int

// Course ...
type Course map[string]Section

// Department ...
type Department map[string]Course
