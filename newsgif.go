package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const itn_url = "https://en.wikipedia.org/w/api.php?action=query&prop=revisions&titles=Template:In_the_news&rvslots=*&rvprop=content&formatversion=2&format=json"

func main() {
	fmt.Println(itn_url)
	resp, err := http.Get(itn_url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
