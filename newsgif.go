package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//const itnUrl = "https://en.wikipedia.org/w/api.php?action=query&prop=revisions&titles=Template:In_the_news&rvslots=*&rvprop=content&formatversion=2&format=json"
const itnUrl = "https://en.wikipedia.org/w/api.php?action=query&titles=Template:In_the_news&formatversion=2&prop=extracts&exintro&explaintext&format=json"
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

/*
func valueOf(v interface{}, key string) string {
	data := v.(map[string]interface{})

	for k, v := range data {
		switch v := v.(type) {
		case string:
			if k == key {
				return v
			}
			fmt.Println(k, v, "(string)")
		case float64:
			fmt.Println(k, v, "(float64)")
		case []interface{}:
			fmt.Println(k, "(array):")
			for i, u := range v {
				fmt.Println("    ", i, u)
			}
		case map[string]interface{}:
			fmt.Println(k, "(map):")
		default:
			fmt.Println(k, v, "(unknown)")
		}
	}

	return "foo"
}
*/

func main() {

	if verbose >= Debug {
		fmt.Println(itnUrl)
	}
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
	if verbose >= Debug {
		fmt.Println(bodyJson)
	}

	if verbose >= Trace {
		bodyJsonPretty, err := json.MarshalIndent(bodyJson, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(bodyJsonPretty))
	}

	//	var content = bodyJson["query"].(map[string]interface{})["pages"].([]interface{})[0].(map[string]interface{})["revisions"].([]interface{})[0].(map[string]interface{})["slots"].(map[string]interface{})["main"].(map[string]interface{})["extract"].(string)
	content, _ := bodyJson["query"].(map[string]interface{})["pages"].([]interface{})[0].(map[string]interface{})["extract"].(string)

	// TODO: needs error handling
	//fmt.Println(ok, content)

	newsLines := strings.Split(strings.TrimSuffix(content, "\n"), "\n")
	for _, v := range newsLines {
		fmt.Println(v, "\n")
	}

	fmt.Println(len(newsLines), cap(newsLines), newsLines[2])
}
