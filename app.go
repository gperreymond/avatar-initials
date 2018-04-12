package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"github.com/julienschmidt/httprouter"
	"github.com/manyminds/api2go"
)

var (
	blue = map[int]*color.RGBA{
		0: &color.RGBA{0x1C, 0x90, 0xF3, 255},
		1: &color.RGBA{0x11, 0x47, 0xCC, 255},
	}
	orange = map[int]*color.RGBA{
		0: &color.RGBA{0xFA, 0xD9, 0x61, 255},
		1: &color.RGBA{0xF7, 0x6B, 0x1C, 255},
	}
	green = map[int]*color.RGBA{
		0: &color.RGBA{0x11, 0xBB, 0xB0, 255},
		1: &color.RGBA{0x18, 0x94, 0x8C, 255},
	}
)

func main() {
	port := 21974

	api := api2go.NewAPI("api")
	handler := api.Handler().(*httprouter.Router)

	handler.GET("/hc", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		type Data struct {
			Alive bool `json:"alive"`
		}
		data := Data{
			Alive: true,
		}
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(data)
	})

	handler.GET("/square", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		q := r.URL.Query()
		// Initialize vars
		text := q.Get("text")
		if text == "" {
			text = "??"
		}
		text = string(text[0:2])
		text = strings.ToUpper(text)
		size, err := strconv.ParseFloat(q.Get("size"), 64)
		if err != nil {
			size = 90
		}
		scale := float64(size / 90)
		width := float64(size)
		height := float64(size)
		// Initialize colors
		i := strings.Index("ABCDEFGHIJKLMNOPQRSTUVWXYZ", string(text[0]))
		j := strings.Index("ABCDEFGHIJKLMNOPQRSTUVWXYZ", string(text[1]))
		k := 0
		if i != -1 && j != -1 {
			k = (i + j*10) % 3
		}
		color := blue
		if k == 0 {
			color = orange
		}
		if k == 1 {
			color = green
		}
		// Initialize image
		dest := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
		dc := gg.NewContextForRGBA(dest)
		// Grandient background
		grad := gg.NewLinearGradient(width, height, 0, 0)
		grad.AddColorStop(0, color[0])
		grad.AddColorStop(1, color[1])
		dc.DrawRoundedRectangle(0, 0, width, height, 5*scale)
		dc.SetFillStyle(grad)
		dc.Fill()
		// Fonts and text
		dc.LoadFontFace("./fonts/RobotoMono-Bold.ttf", 50*scale)
		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(text, width/2, height/2, 0.5, 0.5)
		w.Header().Set("Content-type", "image/png")
		png.Encode(w, dest)
	})

	fmt.Println("Listening on", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}
