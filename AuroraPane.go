package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/Grayda/go-aurora"
	"github.com/ninjasphere/gestic-tools/go-gestic-sdk"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/model"
	"github.com/ninjasphere/sphere-go-led-controller/fonts/O4b03b"
	"github.com/ninjasphere/sphere-go-led-controller/util"
)

type AuroraPane struct {
	siteModel   *ninja.ServiceClient
	site        *model.Site
	dataTimeout *time.Timer
	loaded      bool
	image       util.Image
}

var current int = -1

var results, resultsKp map[int]map[string]float64
var score map[string]int

func NewAuroraPane(conn *ninja.Connection) *AuroraPane {

	pane := &AuroraPane{
		siteModel: conn.GetServiceClient("$home/services/SiteModel"),
		image:     util.LoadImage("images/loading.gif"),
	}

	go func() {
		for {
			fmt.Println("Loading results..")
			results = aurora.Get()
			fmt.Println("Loading Kp results..")
			resultsKp = aurora.GetKp()
			fmt.Println("Calculating score..")
			score := aurora.Check(results, resultsKp, 0)
			fmt.Println(score["Score"])
			pane.loaded = true
			switch {
			case score["Score"] < 25:
				pane.image = util.LoadImage("images/green.gif")
			case score["Score"] >= 25 && score["Score"] <= 50:
				pane.image = util.LoadImage("images/yellow.gif")
			case score["Score"] > 50 && score["Score"] <= 75:
				pane.image = util.LoadImage("images/orange.gif")
			case score["Score"] > 75:
				pane.image = util.LoadImage("images/red.gif")
			}

			time.Sleep(time.Minute)
		}
	}()

	pane.dataTimeout = time.AfterFunc(0, func() {
		current = -1
	})

	return pane
}

func (p *AuroraPane) IsEnabled() bool {
	return true
}

func (p *AuroraPane) KeepAwake() bool {
	return false
}

func (p *AuroraPane) Gesture(gesture *gestic.GestureMessage) {
	if gesture.Tap.Active() && p.loaded == true {
		current++
		if current == 5 {
			current = 0
		}

		p.dataTimeout.Reset(time.Second * 5)

	}
}

func (p *AuroraPane) Render() (*image.RGBA, error) {
	if current > -1 {
		img := image.NewRGBA(image.Rect(0, 0, 16, 16))
		var colour color.RGBA
		// score, bz, speed, density, kp
		switch {
		case score["Score"] == 0 || score["Bz"] == 0 || score["Speed"] == 0 || score["Density"] == 0 || score["Kp"] == 0:
			colour = color.RGBA{0, 255, 0, 255}
		case score["Score"] == 1 || score["Bz"] == 1 || score["Speed"] == 1 || score["Density"] == 1 || score["Kp"] == 1:
			colour = color.RGBA{255, 255, 0, 255}
		case score["Score"] == 2 || score["Bz"] == 2 || score["Speed"] == 2 || score["Density"] == 2 || score["Kp"] == 2:
			colour = color.RGBA{255, 128, 0, 255}
		case score["Score"] == 3 || score["Bz"] == 3 || score["Speed"] == 3 || score["Density"] == 3 || score["Kp"] == 3:
			colour = color.RGBA{255, 0, 0, 255}

		}

		switch current {
		case 0:
			drawText("SC:", colour, 1, img)
			drawText(fmt.Sprintf("%d", score["Score"]), colour, 8, img)
		case 1:
			drawText("KP:", colour, 1, img)
			drawText(fmt.Sprintf("%1.1f", resultsKp[0]["Kp"]), colour, 8, img)
		case 2:
			drawText("BZ:", colour, 1, img)
			drawText(fmt.Sprintf("%2.1f", results[0]["Bz"]), colour, 8, img)
		case 3:
			drawText("SP:", colour, 1, img)
			drawText(fmt.Sprintf("%f", results[0]["Speed"]), colour, 8, img)
		case 4:
			drawText("DN:", colour, 1, img)
			drawText(fmt.Sprintf("%3.1f", results[0]["Density"]), colour, 8, img)

		}
		return img, nil

	} else {
		return p.image.GetNextFrame(), nil
	}
}

func (p *AuroraPane) IsDirty() bool {
	return true
}

// drawText is a helper function to draw a string of text into an image
func drawText(text string, col color.RGBA, top int, img *image.RGBA) {
	width := O4b03b.Font.DrawString(img, 0, 8, text, color.Black)
	start := int(16 - width - 1)

	O4b03b.Font.DrawString(img, start, top, text, col)
}
