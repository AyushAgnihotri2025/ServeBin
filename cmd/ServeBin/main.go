// Copyright 2024 The ServeBin AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"ServeBin/controller"
	_ "ServeBin/docs"
	"ServeBin/helper"
	"ServeBin/router"
	"ServeBin/service"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

// @title           	ServeBin API
// @version         	v1.0.0
// @description			Welcome to ServeBin API documentation! ServeBin is a cutting-edge HTTP testing and debugging tool, built with the latest technologies in Go. This documentation provides comprehensive details about the endpoints, parameters, and responses offered by ServeBin, empowering developers to streamline their testing workflows and ensure the reliability of their applications. Explore the various features and capabilities of ServeBin API to optimize your development process and elevate your HTTP testing experience.
// @termsOfService  	https://servebin.dev/

// @license.name  		BSD-3-Clause
// @license.url  		https://github.com/AyushAgnihotri2025/ServeBin/blob/master/LICENSE

// @host      			servebin.dev
// @BasePath  			/

// @servers				https://servebin.dev https://s1.servebin.dev https://s2.servebin.dev
// @schemes 			http https

// @tag.name			Status Codes
// @tag.description 	Generates responses with given status code

// @tag.name			Request inspection
// @tag.description 	Inspect the request data
func main() {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	} else {
		log.Printf(".env file doesn't exist, assuming you have already set all env vars")
	}

	// Validator
	validate := validator.New()

	// Service
	tagsService := service.NewAPIServiceImpl(validate)

	// Controller
	tagsController := controller.NewAPIController(tagsService)

	// Router
	routes := router.NewRouter(tagsController)

	server := &http.Server{
		Addr:    helper.GetHost(),
		Handler: routes,
	}

	err := server.ListenAndServe()
	helper.ErrorPanic(err)
}
