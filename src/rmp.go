package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// MakeRequest ..
func MakeRequest() {
	// Stored in separate file as it's > 500 characters long
	txt, err := ioutil.ReadFile("rmpquery.txt")

	if err != nil {
		log.Fatalln(err)
	}

	origRmpQuery := string(txt)
	rmpQuery := strings.Replace(origRmpQuery, "NAME", "WOLFMAN", 1)
	log.Print(rmpQuery)

	resp, err := http.Get(rmpQuery)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
}
