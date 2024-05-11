// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"ServeBin/helper"
	"bytes"
	"fmt"
	"github.com/biessek/golang-ico"
	"github.com/kettek/apng"
	"github.com/nickalie/go-webpbin"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"time"
)

type MyError struct{}

func (m *MyError) Error() string {
	return "boom"
}

// GeneratePNG implements APIService
func (t *APIServiceImpl) GeneratePNG() ([]byte, error) {
	// Generate Image
	img := helper.GenerateImage()

	// Encode the image as PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GenerateJPEG implements APIService
func (t *APIServiceImpl) GenerateJPEG() ([]byte, error) {
	// Generate Image
	img := helper.GenerateImage(color.White)

	// Encode the image as JPEG
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GenerateSVG implements APIService
func (t *APIServiceImpl) GenerateSVG() ([]byte, error) {
	// Get current time and date in UTC
	currentTime := time.Now().UTC()
	formattedTime := currentTime.Format("03:04:05 PM MST")
	formattedDate := currentTime.Format("02 January 2006")

	// Define SVG content
	svgContent := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="300" height="70">
		<text x="10" y="30" font-family="Arial" font-size="20" fill="black">%s</text>
		<text x="10" y="60" font-family="Arial" font-size="20" fill="black">%s</text>
	</svg>`, formattedTime, formattedDate)

	return []byte(svgContent), nil
}

// GenerateGIF implements APIService
func (t *APIServiceImpl) GenerateGIF() ([]byte, error) {
	// Create a GIF
	anim := gif.GIF{}

	for i := 0; i < 10; i++ { // Generating 10 frames for a 5-second GIF (assuming 100ms delay per frame)
		// Get current time and date in UTC
		currentTime := time.Now().UTC().Add(time.Duration(i) * time.Second) // Increment time for each frame
		img := helper.GenerateImage(color.White, currentTime)

		// Convert RGBA image to Paletted image
		paletted := image.NewPaletted(img.Bounds(), color.Palette([]color.Color{color.White, color.Black}))

		// Draw the RGBA image onto the Paletted image
		draw.Draw(paletted, paletted.Rect, img, img.Bounds().Min, draw.Src)

		anim.Image = append(anim.Image, paletted)
		anim.Delay = append(anim.Delay, 100) // 100ms delay for each frame
	}

	// Encode the GIF
	var gifBytes bytes.Buffer
	if err := gif.EncodeAll(&gifBytes, &anim); err != nil {
		return nil, err
	}

	return gifBytes.Bytes(), nil
}

// GenerateWEBP implements APIService
func (t *APIServiceImpl) GenerateWEBP() ([]byte, error) {
	// Generate Image
	img := helper.GenerateImage()

	// Encode the image as WEBP
	var buf bytes.Buffer
	// Requires the WebP encoder cwebp
	if err := webpbin.Encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GenerateTIFF implements APIService
func (t *APIServiceImpl) GenerateTIFF() ([]byte, error) {
	// Generate Image
	img := helper.GenerateImage()

	// Encode the image as TIFF
	var buf bytes.Buffer
	if err := tiff.Encode(&buf, img, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GenerateBMP implements APIService
func (t *APIServiceImpl) GenerateBMP() ([]byte, error) {
	// Generate Image
	img := helper.GenerateImage(color.White)

	// Encode the image as BMP
	var buf bytes.Buffer
	if err := bmp.Encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GenerateAPNG implements APIService
func (t *APIServiceImpl) GenerateAPNG() ([]byte, error) {
	// Create a new APNG
	anim := apng.APNG{}

	for i := 0; i < 10; i++ { // Generating 10 frames for a 5-second GIF (assuming 100ms delay per frame)
		// Get current time and date in UTC
		currentTime := time.Now().UTC().Add(time.Duration(i) * time.Second) // Increment time for each frame
		img := helper.GenerateImage(color.White, currentTime)

		// Encode image to PNG
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, img, nil); err != nil {
			return nil, err
		}

		// Create a new APNG frame
		frame := apng.Frame{
			Image:            img,
			DelayNumerator:   1, // Delay in number of ticks
			DelayDenominator: 1, // Delay in seconds = DelayNum / DelayDen
		}

		// Append the frame to APNG
		anim.Frames = append(anim.Frames, frame)
	}

	// Encode the APNG to a buffer
	var apngBuffer bytes.Buffer
	if err := apng.Encode(&apngBuffer, anim); err != nil {
		return nil, err
	}

	return apngBuffer.Bytes(), nil
}

// GenerateAVIF implements APIService
func (t *APIServiceImpl) GenerateAVIF() ([]byte, error) {
	// Check if the .cache directory exists, if not, create it
	cacheDir := ".cache"
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			return nil, err
		}
	}

	img := helper.GenerateImage()

	// Create a temporary PNG file
	pngFile, err := os.CreateTemp(".cache", "image_*.png")
	if err != nil {
		return nil, err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			return
		}
	}(pngFile.Name()) // Clean up the temporary PNG file
	defer func(pngFile *os.File) {
		err := pngFile.Close()
		if err != nil {
			return
		}
	}(pngFile)

	// Save the PNG image to the temporary file
	if err := png.Encode(pngFile, img); err != nil {
		return nil, err
	}

	// Create a temporary AVIF file
	avifFile, err := os.CreateTemp(".cache", "image_*.avif")
	if err != nil {
		return nil, err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			return
		}
	}(avifFile.Name()) // Clean up the temporary AVIF file

	defer func(avifFile *os.File) {
		err := avifFile.Close()
		if err != nil {
			return
		}
	}(avifFile)

	// Convert PNG to AVIF using avifenc
	cmd := exec.Command("avifenc", pngFile.Name(), "-o", avifFile.Name())
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// Read the converted AVIF file
	avifBytes, err := os.ReadFile(avifFile.Name())
	if err != nil {
		return nil, err
	}

	return avifBytes, nil
}

// GenerateICO implements APIService
func (t *APIServiceImpl) GenerateICO() ([]byte, error) {
	// Read Image
	f, err := os.Open("static/sample_image/gopher.png")
	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}

	// Encode the image as ICO
	var buf bytes.Buffer
	if err := ico.Encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
