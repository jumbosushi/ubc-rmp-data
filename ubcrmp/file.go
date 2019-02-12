package ubcrmp

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

func checkIO(e error) {
	if e != nil {
		panic(e)
	}
}
