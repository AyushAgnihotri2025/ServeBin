// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package helper

import (
	"image"
	"image/color"
	"image/draw"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func GenerateImage(options ...interface{}) image.Image {
	// Generated Timestamp if not provided
	var backgroundColor color.Color = color.Transparent // Default background color is transparent
	var currentTime time.Time = time.Now().UTC()        // Default timestamp is current time

	// Process optional parameters
	for _, option := range options {
		switch opt := option.(type) {
		case color.Color:
			backgroundColor = opt
		case time.Time:
			currentTime = opt
		}
	}

	// Format time in HH:MM:SS AM/PM TZ
	formattedTime := currentTime.Format("03:04:05 PM MST")
	formattedDate := currentTime.Format("02 January 2006")

	// Create a new RGBA image with transparent background
	img := image.NewRGBA(image.Rect(0, 0, 300, 70))
	//transparent := color.RGBA{0, 0, 0, 0}
	draw.Draw(img, img.Bounds(), &image.Uniform{backgroundColor}, image.Point{}, draw.Src)

	// Add the current time text to the image
	fontColor := color.RGBA{0, 0, 0, 255} // Black color

	// Write the text to the image
	drawText(img, 100, 30, formattedTime, fontColor)
	drawText(img, 100, 50, formattedDate, fontColor)

	return img
}

// Function to draw text on an image
func drawText(img *image.RGBA, x, y int, text string, fontColor color.RGBA) {
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(fontColor),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)
}
