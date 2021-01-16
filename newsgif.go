package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

const itnUrl = "https://en.wikipedia.org/w/api.php?action=query&prop=revisions&titles=Template:In_the_news&rvslots=*&rvprop=content&formatversion=2&format=json"
const (
	Error   = 1
	Warning = 2
	Info    = 3
	Debug   = 4
	Trace   = 5
)

var verbose int

func init() {
	flag.IntVar(&verbose, "v", 1, "Verbosity level")
	flag.Parse()

	if verbose < Error {
		verbose = Error
	} else if verbose > Trace {
		verbose = Trace
	}
}

func main() {

	fmt.Println(itnUrl)
	resp, err := http.Get(itnUrl)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	bodyRaw, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var bodyJson map[string]interface{}
	json.Unmarshal([]byte(bodyRaw), &bodyJson)

	if verbose >= Trace {
		bodyJsonPretty, err := json.MarshalIndent(bodyJson, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(bodyJsonPretty))
	}

	for k, v := range bodyJson["query"].(map[string]interface{}) {
		fmt.Println("xx", k, "yy", v, "zz")
		if k == "pages" {
			for k2, v2 := range v.([]interface{}) {
				fmt.Println("xxx", k2, "yyy", v2, "zzz")
			}
		}
	}
}
