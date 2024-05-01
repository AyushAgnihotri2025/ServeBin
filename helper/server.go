// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package helper

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Opens the specified URL in the default browser of the user.
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Can't open the server url in the browser! Kindly open it manually " + url)
	}
	if err != nil {
		log.Fatal(err)
	}
}

// Get the server Host
func GetHost() string {
	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	var protocol string
	ssl := os.Getenv("IS_SSL")
	if strings.ToLower(ssl) == "true" {
		protocol = "https"
	} else {
		protocol = "http"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	addr := host + ":" + port

	UrlPrint(addr, protocol)

	if os.Getenv("ENV") != "production" {
		openBrowser(addr)
	}

	return addr
}

// Logs the server URL
func UrlPrint(address string, protocol string) {
	// Parse the URL
	parsedURL, err1 := url.Parse(protocol + "://" + address)
	if err1 != nil {
		fmt.Println("Error parsing URL:", err1)
		return
	}

	// Print the URL
	fmt.Println()
	fmt.Println("ServeBin server started successfully.")
	fmt.Println("\t"+"Server Address:"+"\t", parsedURL.String())
	fmt.Println()
}
