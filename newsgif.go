package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type textImage struct {
	img   *image.RGBA
	width int
}

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

func drawOutlinedText(text string, stroke int, ctx *gg.Context, x float64, y float64) {
	ctx.SetRGBA(0, 0, 0, 1)
	for dy := -stroke; dy <= stroke; dy++ {
		for dx := -stroke; dx <= stroke; dx++ {
			x := x + float64(dx)
			y := y + float64(dy)
			ctx.DrawStringAnchored(text, x, y, 0.5, 0.5)
		}
	}

	ctx.SetRGBA(1, 1, 1, 1)
	ctx.DrawStringAnchored(text, x, y, 0.5, 0.5)
}

func createTextImage(text string, dc *gg.Context) textImage {

	const stroke = 1
	stringWidth, _ := dc.MeasureString(text)
	imgWidth := int(stringWidth + stroke*2 + 300)

	ctx := gg.NewContext(imgWidth, newsLineHeight)
	ctx.SetRGBA(1, 1, 1, 0)
	ctx.Clear()

	if err := ctx.LoadFontFace("/usr/share/fonts/truetype/freefont/FreeSans.ttf", 18); err != nil {
		panic(err)
	}

	drawOutlinedText(text, stroke, ctx, float64(stringWidth/2+300), float64(newsLineHeight/2))

	return textImage{img: ctx.Image().(*image.RGBA), width: imgWidth}
}

func main() {

	newsLines := getHeadlines()

	const width, height = 380, 180
	var images []*image.Paletted
	var delays []int
	var disposals []byte

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

	if err := dc.LoadFontFace("/usr/share/fonts/truetype/freefont/FreeSans.ttf", 24); err != nil {
		panic(err)
	}

	imageNewsLine1 := createTextImage(newsLines[0], dc)
	imageNewsLine2 := createTextImage(newsLines[1], dc)
	imageNewsLine3 := createTextImage(newsLines[2], dc)
	imageNewsLine4 := createTextImage(newsLines[3], dc)

	textMaxWidth := 0
	for i := 0; i < 4; i++ {
		textMaxWidth = imageNewsLine1.width
		if imageNewsLine2.width > textMaxWidth {
			textMaxWidth = imageNewsLine2.width
		}
		if imageNewsLine3.width > textMaxWidth {
			textMaxWidth = imageNewsLine3.width
		}
		if imageNewsLine4.width > textMaxWidth {
			textMaxWidth = imageNewsLine4.width
		}
	}

	for i := 0; i < textMaxWidth; i += 10 {
		dc.SetRGBA(1, 1, 1, 0)
		dc.Clear()

		cropped1 := imageNewsLine1.img.SubImage(image.Rect(i, 0, i+300, 24))
		dc.DrawImage(cropped1, 40-i, 64)
		cropped2 := imageNewsLine2.img.SubImage(image.Rect(i, 0, i+300, 24))
		dc.DrawImage(cropped2, 40-i, 64+1*24)
		cropped3 := imageNewsLine3.img.SubImage(image.Rect(i, 0, i+300, 24))
		dc.DrawImage(cropped3, 40-i, 64+2*24)
		cropped4 := imageNewsLine4.img.SubImage(image.Rect(i, 0, i+300, 24))
		dc.DrawImage(cropped4, 40-i, 64+3*24)

		dc.DrawRoundedRectangle(1, 1, width-2, height-2, 20)
		dc.SetRGBA(0, 0, 0, 1)
		dc.SetLineWidth(3)
		dc.Stroke()
		dc.DrawRoundedRectangle(1, 1, width-2, height-2, 20)
		dc.SetRGBA(1, 1, 1, 1)
		dc.SetLineWidth(1)
		dc.Stroke()

		drawOutlinedText("newsgif", 1, dc, float64(width/2), float64(32))

		img1 := dc.Image()
		bounds := img1.Bounds()

		dst := image.NewPaletted(bounds, palette)
		draw.Draw(dst, bounds, img1, bounds.Min, draw.Src)
		images = append(images, dst)
		delays = append(delays, 20)
		disposals = append(disposals, gif.DisposalBackground)
	}

	f, err := os.OpenFile("rgb.gif", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	gif.EncodeAll(f, &gif.GIF{
		Image:    images,
		Delay:    delays,
		Disposal: disposals,
	})

}
