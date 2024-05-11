// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package response

type GzipResponse struct {
	HeaderResponse
	IPResponse
	Gzipped bool `json:"gzipped"  example:"true"`
}

type BrotliResponse struct {
	HeaderResponse
	IPResponse
	Compressed bool `json:"compressed"  example:"true"`
}

type DeflateResponse struct {
	HeaderResponse
	IPResponse
	Deflated bool `json:"deflated"  example:"true"`
}

type ZstdResponse struct {
	HeaderResponse
	IPResponse
	Compressed bool `json:"compressed"  example:"true"`
}
