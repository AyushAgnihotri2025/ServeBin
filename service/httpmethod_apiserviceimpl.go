// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"ServeBin/data/response"
	"ServeBin/helper"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

// ReturnArguments implements APIService
func (t *APIServiceImpl) ReturnArguments(ctx *gin.Context) (map[string]interface{}, error) {

	// Get URL parameters (args)
	args := make(map[string]interface{})
	for key, values := range ctx.Request.URL.Query() {
		if len(values) == 1 {
			args[key] = values[0]
		} else {
			args[key] = values
		}
	}

	return args, nil
}

// ReturnFormData implements APIService
func (t *APIServiceImpl) ReturnFormData(ctx *gin.Context) (map[string]interface{}, error) {

	form := make(map[string]interface{})
	var err error

	// Generate FormData Response
	if ctx.Request.MultipartForm != nil {
		// Get Form Data
		for key, values := range ctx.Request.MultipartForm.Value {
			if len(values) == 1 {
				form[key] = values[0]
			} else {
				form[key] = values
			}
		}
	} else {
		err = errors.New("No Form Data Found")
	}

	return form, err
}

// ReturnFormFile implements APIService
func (t *APIServiceImpl) ReturnFormFile(ctx *gin.Context) (map[string]interface{}, error) {

	files := make(map[string]interface{})
	var err error

	// Generate FormFile Response
	if ctx.Request.MultipartForm != nil {
		if ctx.Request.MultipartForm.File != nil {
			// Get files
			formFiles := ctx.Request.MultipartForm.File
			for key, values := range formFiles {
				var fileInfo gin.H
				for _, fileHeader := range values {
					file, err := fileHeader.Open()
					if err != nil {
						fmt.Println("Error opening file:", err)
						continue
					}
					defer file.Close()
					fileContents, err := ioutil.ReadAll(file)
					if err != nil {
						fmt.Println("Error reading file:", err)
						continue
					}

					header := helper.GetFileHeaders(fileHeader)
					fileHeaderResponse := response.HeaderResponse{
						Header: header,
					}

					fileInfo = gin.H{
						"Filename":            fileHeader.Filename,
						"Header":              fileHeaderResponse.Header,
						"Size":                fileHeader.Size,
						"Human_Readable_Size": helper.HumanBytes(fileHeader.Size),
						"Content":             fmt.Sprintf("data:%s;base64,%s", fileHeader.Header["Content-Type"][0], base64.StdEncoding.EncodeToString(fileContents)),
					}
				}
				files[key] = fileInfo
			}
		} else {
			err = errors.New("No Form Data Found")
		}
	}

	return files, err
}

// ReturnJson_RawData implements APIService
func (t *APIServiceImpl) ReturnJson_RawData(ctx *gin.Context) (map[string]interface{}, error) {

	data := make(map[string]interface{})
	var err error

	// Get raw data (JSON payload, text, etc.)
	rawData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		data["rawData"] = ""
		data["json"] = ""
	} else {
		var jsonData interface{}
		err = json.Unmarshal(rawData, &jsonData)
		if err == nil {
			data["rawData"] = string(rawData) // assuming raw data is text
			data["json"] = jsonData
		} else {
			data["rawData"] = string(rawData) // assuming raw data is text
			data["json"] = ""
		}
	}

	return data, err
}
