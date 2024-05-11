// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetImages 	 	ServeBin
// @Tags			Images
// @Summary			Returns image as per the Accept header.
// @Description		Returns a simple image of the type suggest by the Accept header.
// @produce			image/png
// @produce			image/jpeg
// @produce			image/svg
// @produce			image/gif
// @produce			image/webp
// @produce			image/tiff
// @produce			image/bmp
// @produce			image/apng
// @produce			image/avif
// @produce			image/x-icon
// @produce			image/*
// @Success			200
// @Failure      	406
// @Failure      	404
// @Failure      	500  		{object}  	response.HTTPError
// @Router			/image 		[get]
// @Router			/image/png 	[get]
// @Router			/image/jpeg [get]
// @Router			/image/svg 	[get]
// @Router			/image/gif 	[get]
// @Router			/image/webp	[get]
// @Router			/image/tiff	[get]
// @Router			/image/bmp 	[get]
// @Router			/image/apng	[get]
// @Router			/image/avif	[get]
// @Router			/image/ico 	[get]
func (controller *APIController) GetImages(ctx *gin.Context) {

	var imageBytes []byte
	var err error

	switch imagetype := ctx.Param("imagetype"); imagetype {
	case "":
		switch accept := ctx.Request.Header["Accept"][0]; accept {
		case "image/png":
			// Generate PNG
			imageBytes, err = controller.apiService.GeneratePNG()
			if err != nil {
				ctx.String(http.StatusInternalServerError, "Failed to generate PNG: "+err.Error())
				return
			}

			// Set Content-Type header
			ctx.Header("Content-Type", "image/png")

			// Set Content-Disposition header
			ctx.Header("Content-Disposition", "inline; filename=image.png")
			break

		case "image/svg+xml":
			// Generate SVG
			imageBytes, err = controller.apiService.GenerateSVG()
			if err != nil {
				ctx.String(http.StatusInternalServerError, "Failed to generate SVG: "+err.Error())
				return
			}

			// Set Content-Type header
			ctx.Header("Content-Type", "image/svg+xml")

			// Set Content-Disposition header
			ctx.Header("Content-Disposition", "inline; filename=image.svg")
			break

		case "image/gif":
			// Generate GIF
			imageBytes, err = controller.apiService.GenerateGIF()
			if err != nil {
				ctx.String(http.StatusInternalServerError, "Failed to generate GIF: "+err.Error())
				return
			}

			// Set Content-Type header
			ctx.Header("Content-Type", "image/gif")

			// Set Content-Disposition header
			ctx.Header("Content-Disposition", "inline; filename=image.gif")
			break

		case "image/jpeg":
			// Generate JPEG
			imageBytes, err = controller.apiService.GenerateJPEG()
			if err != nil {
				ctx.String(http.StatusInternalServerError, "Failed to generate JPEG: "+err.Error())
				return
			}

			// Set Content-Type header
			ctx.Header("Content-Type", "image/jpeg")

			// Set Content-Disposition header
			ctx.Header("Content-Disposition", "inline; filename=image.jpeg")
			break

		case "image/webp":
			// Generate WEBP
			imageBytes, err = controller.apiService.GenerateWEBP()
			if err != nil {
				ctx.String(http.StatusInternalServerError, "Failed to generate WEBP: "+err.Error())
				return
			}

			// Set Content-Type header
			ctx.Header("Content-Type", "image/webp")

			// Set Content-Disposition header
			ctx.Header("Content-Disposition", "inline; filename=image.webp")
			break

		case "image/tiff":
			// Generate TIFF
			imageBytes, err = controller.apiService.GenerateTIFF()
			if err != nil {
				ctx.String(http.StatusInternalServerError, "Failed to generate TIFF: "+err.Error())
				return
			}

			// Set Content-Type header
			ctx.Header("Content-Type", "image/tiff")

			// Set Content-Disposition header
			ctx.Header("Content-Disposition", "inline; filename=image.tiff")
			break

		case "image/bmp":
			// Generate BMP
			imageBytes, err = controller.apiService.GenerateBMP()
			if err != nil {
				ctx.String(http.StatusInternalServerError, "Failed to generate BMP: "+err.Error())
				return
			}

			// Set Content-Type header
			ctx.Header("Content-Type", "image/bmp")

			// Set Content-Disposition header
			ctx.Header("Content-Disposition", "inline; filename=image.bmp")
			break

		case "image/apng":
			// Generate APNG
			imageBytes, err = controller.apiService.GenerateAPNG()
			if err != nil {
				ctx.String(http.StatusInternalServerError, "Failed to generate APNG: "+err.Error())
				return
			}

			// Set Content-Type header
			ctx.Header("Content-Type", "image/apng")

			// Set Content-Disposition header
			ctx.Header("Content-Disposition", "inline; filename=image.apng")
			break

		case "image/avif":
			// Generate AVIF
			imageBytes, err = controller.apiService.GenerateAVIF()
			if err != nil {
				ctx.String(http.StatusInternalServerError, "Failed to generate AVIF: "+err.Error())
				return
			}

			// Set Content-Type header
			ctx.Header("Content-Type", "image/avif")

			// Set Content-Disposition header
			ctx.Header("Content-Disposition", "inline; filename=image.avif")
			break

		case "image/x-icon":
			// Generate ICO
			imageBytes, err = controller.apiService.GenerateICO()
			if err != nil {
				ctx.String(http.StatusInternalServerError, "Failed to generate ICO: "+err.Error())
				return
			}

			// Set Content-Type header
			ctx.Header("Content-Type", "image/x-icon")

			// Set Content-Disposition header
			ctx.Header("Content-Disposition", "inline; filename=image.ico")
			break

		default:
			// Generate PNG
			imageBytes, err = controller.apiService.GeneratePNG()
			if err != nil {
				ctx.String(http.StatusInternalServerError, "Failed to generate PNG: "+err.Error())
				return
			}

			// Set Content-Type header
			ctx.Header("Content-Type", "image/png")

			// Set Content-Disposition header
			ctx.Header("Content-Disposition", "inline; filename=image.png")
			break
		}
		break

	case "png":
		// Generate PNG
		imageBytes, err = controller.apiService.GeneratePNG()
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to generate PNG: "+err.Error())
			return
		}

		// Set Content-Type header
		ctx.Header("Content-Type", "image/png")

		// Set Content-Disposition header
		ctx.Header("Content-Disposition", "inline; filename=image.png")
		break

	case "svg":
		// Generate SVG
		imageBytes, err = controller.apiService.GenerateSVG()
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to generate SVG: "+err.Error())
			return
		}

		// Set Content-Type header
		ctx.Header("Content-Type", "image/svg+xml")

		// Set Content-Disposition header
		ctx.Header("Content-Disposition", "inline; filename=image.svg")
		break

	case "gif":
		// Generate GIF
		imageBytes, err = controller.apiService.GenerateGIF()
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to generate GIF: "+err.Error())
			return
		}

		// Set Content-Type header
		ctx.Header("Content-Type", "image/gif")

		// Set Content-Disposition header
		ctx.Header("Content-Disposition", "inline; filename=image.gif")
		break

	case "jpeg":
		// Generate JPEG
		imageBytes, err = controller.apiService.GenerateJPEG()
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to generate JPEG: "+err.Error())
			return
		}

		// Set Content-Type header
		ctx.Header("Content-Type", "image/jpeg")

		// Set Content-Disposition header
		ctx.Header("Content-Disposition", "inline; filename=image.jpeg")
		break

	case "webp":
		// Generate WEBP
		imageBytes, err = controller.apiService.GenerateWEBP()
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to generate WEBP: "+err.Error())
			return
		}

		// Set Content-Type header
		ctx.Header("Content-Type", "image/webp")

		// Set Content-Disposition header
		ctx.Header("Content-Disposition", "inline; filename=image.webp")
		break

	case "tiff":
		// Generate TIFF
		imageBytes, err = controller.apiService.GenerateTIFF()
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to generate TIFF: "+err.Error())
			return
		}

		// Set Content-Type header
		ctx.Header("Content-Type", "image/tiff")

		// Set Content-Disposition header
		ctx.Header("Content-Disposition", "inline; filename=image.tiff")
		break

	case "bmp":
		// Generate BMP
		imageBytes, err = controller.apiService.GenerateBMP()
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to generate BMP: "+err.Error())
			return
		}

		// Set Content-Type header
		ctx.Header("Content-Type", "image/bmp")

		// Set Content-Disposition header
		ctx.Header("Content-Disposition", "inline; filename=image.bmp")
		break

	case "apng":
		// Generate APNG
		imageBytes, err = controller.apiService.GenerateAPNG()
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to generate APNG: "+err.Error())
			return
		}

		// Set Content-Type header
		ctx.Header("Content-Type", "image/apng")

		// Set Content-Disposition header
		ctx.Header("Content-Disposition", "inline; filename=image.apng")
		break

	case "avif":
		// Generate AVIF
		imageBytes, err = controller.apiService.GenerateAVIF()
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to generate AVIF: "+err.Error())
			return
		}

		// Set Content-Type header
		ctx.Header("Content-Type", "image/avif")

		// Set Content-Disposition header
		ctx.Header("Content-Disposition", "inline; filename=image.avif")
		break

	case "ico":
		// Generate ICO
		imageBytes, err = controller.apiService.GenerateICO()
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to generate ICO: "+err.Error())
			return
		}

		// Set Content-Type header
		ctx.Header("Content-Type", "image/x-icon")

		// Set Content-Disposition header
		ctx.Header("Content-Disposition", "inline; filename=image.ico")
		break

	default:
		ctx.Status(http.StatusNotAcceptable)
		return
	}

	// Set Content-Length header
	ctx.Header("Content-Length", strconv.Itoa(len(imageBytes)))

	// Write JPEG to response body
	ctx.Writer.Write(imageBytes)
}
