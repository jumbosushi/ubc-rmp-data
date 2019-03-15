package ubcrmp

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

func getFilePath(fname string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir + "/../data/" + fname
}

func writeJSON(class interface{}, fname string) {
	fpath := getFilePath(fname)
	jsonString, err := json.Marshal(class)
	err = ioutil.WriteFile(fpath, jsonString, 0644)
	checkIO(err)
}

func getTermFileName(termURL string, filename string) string {
	u, err := url.Parse(termURL)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	year := q.Get("sessyr")
	term := q.Get("sesscd")
	campus := q.Get("campuscd")

	return year + term + campus + filename
}

func checkIO(e error) {
	if e != nil {
		panic(e)
	}
}
