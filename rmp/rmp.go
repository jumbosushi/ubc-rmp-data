package rmp

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	// "github.com/jumbosushi/ubc-rmp-scraper/model"
)

// MakeRequest ..
func MakeRequest() {
	txt, err := ioutil.ReadFile("./rmp/rmpquery.txt")

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
