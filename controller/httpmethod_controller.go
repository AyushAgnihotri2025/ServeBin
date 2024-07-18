// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	"ServeBin/data/response"
	"ServeBin/helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ResponseData 		ServeBin
// @Tags				HTTP Methods
// @Summary				Returns the request parameters.
// @Description			Returns different type of request parameters like form data, json data, raw data, headers etc.
// @Param        		customheader  	header  string  false  "Header"
// @Param        		queryparam  	query  	string  false  "Query Paramater"
// @Default				200			{object}	response.EmptyResponse
// @Success				200			{object}	response.EmptyResponse
// @Failure      		400
// @Failure      		404
// @Failure      		500  		{object}  	response.HTTPError
// @Router				/get		[get]
func (controller *APIController) ResponseData(ctx *gin.Context) {
	// Get URL parameters (args)
	args, _ := controller.apiService.ReturnArguments(ctx)

	// Get request headers
	header := helper.GetHeaders(ctx)

	// Get origin
	origin := controller.apiService.FindIP(ctx)

	// Get URL
	url := ctx.Request.URL.String()

	// Add method
	method := ctx.Request.Method

	webResponse := response.EmptyResponse{
		ParamResponse:  response.ParamResponse{Parma: args},
		HeaderResponse: response.HeaderResponse{Header: header},
		IPResponse:     response.IPResponse{IP: []interface{}{origin}},
		Url:            url,
		Method:         method,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

// ResponseHeaderData	ServeBin
// @Tags				HTTP Methods
// @Summary				Returns the request parameters.
// @Description			Returns different type of request parameters like form data, json data, raw data, headers etc.
// @Param        		customheader  	header  	string  false  "Header"
// @Default				200
// @Success				200
// @Failure      		400
// @Failure      		404
// @Failure      		500  		{object}  	response.HTTPError
// @Router				/head		[head]
// @Router				/options	[options]
func (controller *APIController) ResponseHeaderData(ctx *gin.Context) {

	// Get request headers
	for key, values := range ctx.Request.Header {
		for _, value := range values {
			ctx.Writer.Header().Add(key, value)
		}
	}

	// Get origin
	origin := controller.apiService.FindIP(ctx)

	// Get URL
	url := ctx.Request.URL.String()

	// Add method
	method := ctx.Request.Method

	ctx.Header("origin", origin)
	ctx.Header("url", url)
	ctx.Header("method", method)
	ctx.JSON(http.StatusOK, nil)
}

// ResponseBodyData 	ServeBin
// @Tags				HTTP Methods
// @Summary				Returns the request parameters.
// @Description			Returns different type of request parameters like form data, json data, raw data, headers etc.
// @Param        		body 			formData	string  false  "Body"
// @Param        		formdata  		formData  	file  	false  "Form Data"
// @Param        		customheader  	header  	string  false  "Header"
// @Param        		queryparam  	query  		string  false  "Query Paramater"
// @Default				200			{object}	response.BodyDataResponse
// @Success				200			{object}	response.BodyDataResponse
// @Failure      		400
// @Failure      		404
// @Failure      		500  		{object}  	response.HTTPError
// @Router				/post		[post]
// @Router				/put		[put]
// @Router				/patch 		[patch]
func (controller *APIController) ResponseBodyData(ctx *gin.Context) {
	// Get URL parameters (args)
	args, _ := controller.apiService.ReturnArguments(ctx)

	// Initialized MultiPart Form Data
	ctx.Request.ParseMultipartForm(10 << 20) // maxMemory 10 MB

	// Get form data
	form, _ := controller.apiService.ReturnFormData(ctx)

	// Get form files
	files, _ := controller.apiService.ReturnFormFile(ctx)

	// Get raw data (JSON payload, text, etc.)
	rawData, _ := controller.apiService.ReturnJson_RawData(ctx)

	// Get request headers
	header := helper.GetHeaders(ctx)

	// Get origin
	origin := controller.apiService.FindIP(ctx)

	// Get URL
	url := ctx.Request.URL.String()

	// Add method
	method := ctx.Request.Method

	webResponse := response.BodyDataResponse{
		ParamResponse:  response.ParamResponse{Parma: args},
		DataResponse:   response.DataResponse{Data: rawData["rawData"]},
		FileResponse:   response.FileResponse{File: files},
		FormResponse:   response.FormResponse{Form: form},
		HeaderResponse: response.HeaderResponse{Header: header},
		JsonResponse:   response.JsonResponse{Json: rawData["json"]},
		IPResponse:     response.IPResponse{IP: []interface{}{origin}},
		Url:            url,
		Method:         method,
	}

	ctx.JSON(http.StatusOK, webResponse)
}
