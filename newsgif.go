package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	//"image/color/palette"
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
const newsLineHeight = 24

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

/*
func drawNewsLine(idx int, text string, width int, height int, dc *gg.Context) {

	const stroke = 1
	topOffset := 80 + 24*idx
	//stringWidth, _ := fmt.Println(dc.MeasureString(text))

	dc.SetRGBA(0, 0, 0, 1)
	for dy := -stroke; dy <= stroke; dy++ {
		for dx := -stroke; dx <= stroke; dx++ {
			x := float64(width/2) + float64(dx)
			y := float64(topOffset) + float64(dy)
			dc.DrawStringAnchored(text, x, y, 0.5, 0.5)
		}
	}

	dc.SetRGBA(1, 1, 1, 1)
	dc.DrawStringAnchored(text, float64(width/2), float64(topOffset), 0.5, 0.5)
}
*/

/*
func mergeNewsLine(idx int, newsLine *image.Image, width int, height int, dc *gg.Context) {

	topOffset := 80 + 24*idx

	dc.SetRGBA(0, 0, 0, 1)
	for dy := -stroke; dy <= stroke; dy++ {
		for dx := -stroke; dx <= stroke; dx++ {
			x := float64(width/2) + float64(dx)
			y := float64(topOffset) + float64(dy)
			dc.DrawStringAnchored(text, x, y, 0.5, 0.5)
		}
	}

	dc.SetRGBA(1, 1, 1, 1)
	dc.DrawStringAnchored(text, float64(width/2), float64(topOffset), 0.5, 0.5)
}
*/

func createTextImage(text string, dc *gg.Context) image.Image {

	const stroke = 1
	stringWidth, _ := dc.MeasureString(text)

	ctx := gg.NewContext(int(stringWidth+stroke*2), newsLineHeight)
	ctx.SetRGBA(1, 1, 1, 0)
	ctx.Clear()

	ctx.SetRGBA(0, 0, 0, 1)
	for dy := -stroke; dy <= stroke; dy++ {
		for dx := -stroke; dx <= stroke; dx++ {
			x := float64(stringWidth/2) + float64(dx)
			y := float64(newsLineHeight/2) + float64(dy)
			ctx.DrawStringAnchored(text, x, y, 0.5, 0.5)
		}
	}

	ctx.SetRGBA(1, 1, 1, 1)
	ctx.DrawStringAnchored(text, float64(stringWidth/2), float64(newsLineHeight/2), 0.5, 0.5)

	return ctx.Image()
}

func main() {

	//	newsLines := getHeadlines()
	//	fmt.Println(len(newsLines), cap(newsLines), newsLines[2])

	const width, height = 380, 180
	var images []*image.Paletted
	var delays []int

	/*
		const S = 1024
		dc := gg.NewContext(S, S)
		dc.SetRGB(1, 1, 1)
		dc.Clear()

		if err := dc.LoadFontFace("/usr/share/fonts/truetype/freefont/FreeSans.ttf", 18); err != nil {
			panic(err)
		}
		dc.SetRGB(0, 0, 0)
		s := "ONE DOES NOT SIMPLY"
		n := 2 // "stroke" size
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
	*/

	//var line1 = "Yoweri Museveni (pictured) is re-elected as President of Uganda."
	//var line2 = "Dutch prime minister Mark Rutte and his cabinet resign as a result of a child welfare fraud scandal."
	//var line3 = "An earthquake on the Indonesian island of Sulawesi kills at least 92 people and injures more than 900 others."
	var line4 = "Donald Trump becomes the first U.S. president to be impeached twice after the House of Representatives charges him with incitement of insurrection."

	var palette color.Palette = color.Palette{
		image.Transparent,
		image.Black,
		image.White,
		color.RGBA{0xEE, 0xEE, 0xEE, 255},
		color.RGBA{0xCC, 0xCC, 0xCC, 255},
		color.RGBA{0x99, 0x99, 0x99, 255},
		color.RGBA{0x66, 0x66, 0x66, 255},
		color.RGBA{0x33, 0x33, 0x33, 255},
	}
	dc := gg.NewContext(width, height)
	dc.SetRGBA(1, 1, 1, 0)
	dc.Clear()

	if err := dc.LoadFontFace("/usr/share/fonts/truetype/freefont/FreeSans.ttf", 18); err != nil {
		panic(err)
	}
	dc.SetRGBA(0, 0, 0, 1)
	//	drawNewsLine(0, line4, width, height, dc)
	//	drawNewsLine(1, line4, width, height, dc)
	//	drawNewsLine(2, line4, width, height, dc)
	//	drawNewsLine(3, line4, width, height, dc)

	imageNewsLine4 := createTextImage(line4, dc)

	rgba := imageNewsLine4.(*image.RGBA)
	cropped := rgba.SubImage(image.Rect(250, 0, 270, 20))
	dc.DrawImage(cropped, 0, 100)

	img1 := dc.Image()
	bounds := img1.Bounds()

	dst := image.NewPaletted(bounds, palette)
	draw.Draw(dst, bounds, img1, bounds.Min, draw.Src)
	images = append(images, dst)
	delays = append(delays, 0)

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
