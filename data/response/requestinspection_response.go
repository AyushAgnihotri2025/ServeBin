// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package response

type IPResponse struct {
	IP []interface{} `json:"ip,omitempty"`
}

type UserAgentResponse struct {
	UserAgent string `json:"user-agent"`
}

type HeaderResponse struct {
	Header interface{} `json:"headers,omitempty"`
}
