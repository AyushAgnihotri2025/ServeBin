// Copyright 2026 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package response

type HTTPBinRequestResponse struct {
	Args    map[string]interface{} `json:"args"`
	Data    string                 `json:"data"`
	Files   map[string]interface{} `json:"files"`
	Form    map[string]interface{} `json:"form"`
	Headers map[string]string      `json:"headers"`
	Json    interface{}            `json:"json"`
	Method  string                 `json:"method,omitempty"`
	Origin  string                 `json:"origin"`
	Url     string                 `json:"url"`
}

type HTTPBinStreamResponse struct {
	HTTPBinRequestResponse
	Id int `json:"id"`
}

type UUIDResponse struct {
	UUID string `json:"uuid"`
}

type CookiesResponse struct {
	Cookies map[string]string `json:"cookies"`
}

type AuthResponse struct {
	Authenticated bool   `json:"authenticated"`
	User          string `json:"user"`
}
