// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package response

type ParamResponse struct {
	Parma map[string]interface{} `json:"arg,omitempty"`
}

type DataResponse struct {
	Data interface{} `json:"data,omitempty"`
}

type FileResponse struct {
	File interface{} `json:"files,omitempty"`
}

type FormResponse struct {
	Form interface{} `json:"form,omitempty"`
}

type JsonResponse struct {
	Json interface{} `json:"json,omitempty"`
}

type EmptyResponse struct {
	ParamResponse
	HeaderResponse
	IPResponse
	Url    string `json:"url,omitempty"`
	Method string `json:"method,omitempty"`
}

type BodyDataResponse struct {
	ParamResponse
	DataResponse
	FileResponse
	FormResponse
	HeaderResponse
	JsonResponse
	IPResponse
	Url    string `json:"url,omitempty"`
	Method string `json:"method,omitempty"`
}

type PutResponse struct {
	ParamResponse
	DataResponse
	FileResponse
	FormResponse
	HeaderResponse
	JsonResponse
	IPResponse
	Url    string `json:"url,omitempty"`
	Method string `json:"method,omitempty"`
}

type DeleteResponse struct {
	ParamResponse
	DataResponse
	FileResponse
	FormResponse
	HeaderResponse
	JsonResponse
	IPResponse
	Url    string `json:"url,omitempty"`
	Method string `json:"method,omitempty"`
}
