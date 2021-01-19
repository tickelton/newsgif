package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fogleman/gg"
	"image"
	//"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

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

func getHeadlines() []string {
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

	content, _ := bodyJson["query"].(map[string]interface{})["pages"].([]interface{})[0].(map[string]interface{})["extract"].(string)

	// TODO: needs error handling
	//fmt.Println(ok, content)

	newsLines := strings.Split(strings.TrimSuffix(content, "\n"), "\n")
	return newsLines
}

func main() {

	//	newsLines := getHeadlines()
	//	fmt.Println(len(newsLines), cap(newsLines), newsLines[2])

	//var width, height int = 380, 180
	var images []*image.Paletted
	var delays []int

	const S = 1024
	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	if err := dc.LoadFontFace("/usr/share/fonts/truetype/freefont/FreeSans.ttf", 96); err != nil {
		panic(err)
	}
	dc.SetRGB(0, 0, 0)
	s := "ONE DOES NOT SIMPLY"
	n := 6 // "stroke" size
	for dy := -n; dy <= n; dy++ {
		for dx := -n; dx <= n; dx++ {
			if dx*dx+dy*dy >= n*n {
				// give it rounded corners
				continue
			}
			x := S/2 + float64(dx)
			y := S/2 + float64(dy)
			dc.DrawStringAnchored(s, x, y, 0.5, 0.5)
		}
	}

	img1 := dc.Image()
	bounds := img1.Bounds()

	dst := image.NewPaletted(bounds, palette.Plan9)
	draw.Draw(dst, bounds, img1, bounds.Min, draw.Src)
	images = append(images, dst)
	delays = append(delays, 20)
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(s, S/2, S/2, 0.5, 0.5)

	img2 := dc.Image()
	bounds2 := img2.Bounds()

	dst2 := image.NewPaletted(bounds2, palette.Plan9)
	draw.Draw(dst2, bounds2, img2, bounds2.Min, draw.Src)
	images = append(images, dst2)
	delays = append(delays, 20)
	f, err := os.OpenFile("rgb.gif", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
	})

}
