package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

const itn_url = "https://en.wikipedia.org/w/api.php?action=query&prop=revisions&titles=Template:In_the_news&rvslots=*&rvprop=content&formatversion=2&format=json"
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

	fmt.Println(itn_url)
	resp, err := http.Get(itn_url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body_raw, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var body_json map[string]interface{}
	json.Unmarshal([]byte(body_raw), &body_json)

	if verbose >= Trace {
		body_json_pretty, err := json.MarshalIndent(body_json, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(body_json_pretty))
	}
}
