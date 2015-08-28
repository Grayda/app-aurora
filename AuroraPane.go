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
	getWeather  *time.Timer
	dataTimeout *time.Timer

	bz      bool
	speed   bool
	density bool
	kp      bool

	image util.Image
}

var current int = 0
var results, resultsKp map[int]map[string]float64

func NewAuroraPane(conn *ninja.Connection) *AuroraPane {

	pane := &AuroraPane{
		siteModel: conn.GetServiceClient("$home/services/SiteModel"),
		image:     util.LoadImage("images/loading.gif"),
	}

	pane.dataTimeout = time.AfterFunc(0, func() {
		pane.bz = false
		pane.kp = false
		pane.density = false
		pane.speed = false
		current = 0
	})

	go func() {
		for {
			results = aurora.Get()
			resultsKp = aurora.GetKp()
			switch {
			case resultsKp[0]["Kp"] < 5:
				pane.image = util.LoadImage("images/green.gif")
			case resultsKp[0]["Kp"] >= 5 && resultsKp[0]["Kp"] <= 7:
				pane.image = util.LoadImage("images/orange.gif")
			case resultsKp[0]["Kp"] > 7:
				pane.image = util.LoadImage("images/red.gif")

			}

			time.Sleep(time.Minute)
		}
	}()

	return pane
}

func (p *AuroraPane) IsEnabled() bool {
	return true
}

func (p *AuroraPane) KeepAwake() bool {
	return false
}

func (p *AuroraPane) Gesture(gesture *gestic.GestureMessage) {
	if gesture.Tap.Active() {

		switch current {
		case 0:
			p.kp = true
			p.bz = false
			p.speed = false
			p.density = false
			current++
		case 1:
			p.kp = false
			p.bz = true
			p.speed = false
			p.density = false
			current++
		case 2:
			p.kp = false
			p.bz = false
			p.speed = true
			p.density = false
			current++
		case 3:
			p.kp = false
			p.bz = false
			p.speed = false
			p.density = true
			current = 0
		}

		p.dataTimeout.Reset(time.Second * 5)

	}
}

func (p *AuroraPane) Render() (*image.RGBA, error) {
	if p.bz || p.kp || p.density || p.speed {
		img := image.NewRGBA(image.Rect(0, 0, 16, 16))

		switch {
		case p.kp == true:
			drawText("K:", color.RGBA{253, 151, 32, 255}, 1, img)
			drawText(fmt.Sprintf("%1.1f", resultsKp[0]["Kp"]), color.RGBA{253, 151, 32, 255}, 8, img)
		case p.bz == true:
			drawText("B:", color.RGBA{253, 151, 32, 255}, 1, img)
			drawText(fmt.Sprintf("%2.1f", results[0]["Bz"]), color.RGBA{253, 151, 32, 255}, 8, img)
		case p.speed == true:
			drawText("S:", color.RGBA{253, 151, 32, 255}, 1, img)
			drawText(fmt.Sprintf("%4.0f", results[0]["Speed"]), color.RGBA{253, 151, 32, 255}, 8, img)
		case p.density == true:
			drawText("D:", color.RGBA{253, 151, 32, 255}, 1, img)
			drawText(fmt.Sprintf("%3.1f", results[0]["Density"]), color.RGBA{253, 151, 32, 255}, 8, img)

		}
		return img, nil

	} else {

		drawText(fmt.Sprintf("%1.f", resultsKp[0]["Kp"]), color.RGBA{253, 151, 32, 255}, 1, image.NewRGBA(image.Rect(0, 0, 16, 16)))
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
