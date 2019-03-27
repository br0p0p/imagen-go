package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	// "image/png"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

var errInvalidFormat = errors.New("invalid format")

func Handler(w http.ResponseWriter, r *http.Request) {
	config := ImageConfig{}

	query_params := r.URL.Query()

	fmt.Println("*** NEW REQUEST ***")
	fmt.Println("query params:")
	for k, v := range query_params {
		fmt.Printf("key[%s] value[%s]\n", k, v)
	}

	if width, err := strconv.Atoi(query_params.Get("width")); err == nil {
		config.Width = width
	}

	if height, err := strconv.Atoi(query_params.Get("height")); err == nil {
		if height == 0 {
			config.Height = config.Width
		} else {
			config.Height = height

			if config.Width == 0 {
				config.Width = config.Height
			}
		}
	} else {
		config.Height = config.Width
	}

	/* Fail on 0x0 dimensions */
	if config.Height == 0 && config.Width == 0 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		c := query_params.Get("color")
		fmt.Printf("%T, %v\n", c, c)

		if len(c) == 0 {
			fmt.Println("color query param not found")
		} else {
			config.Color = c
		}

		fmt.Printf("ImageConfig: %+v\n", config)

		img := NewImage(config)

		img_buffer := EncodeImageToBuffer(img, "jpeg")

		// Set headers
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(img_buffer.Bytes())))
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(img_buffer.Bytes()); err != nil {
			fmt.Println("unable to write image...")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func ParseHexColor(s string) color.RGBA {
	normalize := func(s string) (string, error) {
		var b strings.Builder

		switch len(s) {
		case 8:
			b.WriteString(s)
		case 6:
			b.WriteString(s)
			b.WriteString("FF")
		case 3:
			for _, char := range s {
				b.WriteRune(char)
				b.WriteRune(char)
			}
			b.WriteString("FF")
		default:
			return s, errors.New("couldn't normalize the color string")
		}

		s = b.String()

		fmt.Println("normalized color:", s)

		return s, nil
	}

	colorStr, err := normalize(s)
	if err != nil {
		return RandomColor()
	}

	b, err := hex.DecodeString(colorStr)
	if err != nil {
		return RandomColor()
	}

	return color.RGBA{b[0], b[1], b[2], b[3]}
}

func RandomChannelValue() uint8 {
	return uint8(rand.Intn(255))
}

func RandomColor() color.RGBA {
	r := RandomChannelValue()
	g := RandomChannelValue()
	b := RandomChannelValue()
	a := RandomChannelValue()
	return color.RGBA{r, g, b, a}
}

type ImageConfig struct {
	Width  int
	Height int
	Color  string
}

func NewImage(image_config ImageConfig) *image.RGBA {
	upLeft := image.Point{0, 0}
	downRight := image.Point{image_config.Width, image_config.Height}

	img_color := ParseHexColor(image_config.Color)

	fmt.Printf("Image color: %+v\n", img_color)

	img := image.NewRGBA(image.Rectangle{upLeft, downRight})
	draw.Draw(img, img.Bounds(), &image.Uniform{img_color}, image.ZP, draw.Src)

	return img
}

func EncodeImageToBuffer(img *image.RGBA, format string) *bytes.Buffer {
	buffer := new(bytes.Buffer)

	if format == "jpeg" {
		jpeg.Encode(buffer, img, nil)
	} else {
		fmt.Printf("unkown image format %s\n", format)
	}

	return buffer
}

/* Uncomment to develop locally */

// func main() {
// 	http.HandleFunc("/", Handler)

// 	port := 8080

// 	// Handle static files
// 	fs := http.FileServer(http.Dir("static/"))
// 	http.Handle("/static/", http.StripPrefix("/static/", fs))

// 	fmt.Println(fmt.Sprintf("Server listening at http://localhost:%d", port))
// 	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
// }
