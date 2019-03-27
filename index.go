package main

import (
	"bytes"
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
)

var errInvalidFormat = errors.New("invalid format")

// Stolen from https://stackoverflow.com/a/54200713
func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff

	if s[0] != '#' {
		return c, errInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a'
		case b >= 'A' && b <= 'F':
			return b - 'A'
		}
		err = errInvalidFormat
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		err = errInvalidFormat
	}
	return
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
	Color  string "black"
}

func NewImageConfig() ImageConfig {
	image_config := ImageConfig{}
	image_config.Color = "black"
	return image_config
}

func NewImage(image_config ImageConfig) *image.RGBA {
	upLeft := image.Point{0, 0}
	downRight := image.Point{image_config.Width, image_config.Height}

	img_color, err := ParseHexColor(image_config.Color)

	if err != nil {
		fmt.Printf("unable to parse color %s\n", image_config.Color)
		// img_color := color.RGBA{10, 10, 10, 255}
		img_color = RandomColor()
	}

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

func Handler(w http.ResponseWriter, r *http.Request) {
	config := NewImageConfig()

	query_params := r.URL.Query()

	fmt.Println("query params:")
	for k, v := range query_params {
		fmt.Printf("key[%s] value[%s]\n", k, v)
	}

	if width, err := strconv.Atoi(query_params.Get("width")); err == nil {
		// fmt.Printf("%T, %v\n", width, width)
		config.Width = width
	}

	if height, err := strconv.Atoi(query_params.Get("height")); err == nil {
		// fmt.Printf("%T, %v\n", height, height)
		config.Height = height
	}

	c := query_params.Get("color")
	// fmt.Printf("%T, %v\n", c, c)

	if c == "" {
		fmt.Println("color query param not found")
		// config.Color = "#123321"
		// config.Color = strconv.Parse
	} else {
		if c[0] != '#' {
			c = "#" + c
		}
		config.Color = c
	}

	fmt.Printf("ImageConfig: %+v\n", config)

	// fmt.Fprintf(w, "<h1>Hello from Go!")

	img := NewImage(config)

	img_buffer := EncodeImageToBuffer(img, "jpeg")

	// Set headers
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(img_buffer.Bytes())))

	if _, err := w.Write(img_buffer.Bytes()); err != nil {
		fmt.Println("unable to write image...")
	}
}

func main() {
	http.HandleFunc("/", Handler)

	port := 8080

	// Handle static files
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println(fmt.Sprintf("Server listening at http://localhost:%d", port))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
