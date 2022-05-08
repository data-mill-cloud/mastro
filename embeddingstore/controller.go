package main

import (
	"fmt"
	"net/http"

	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	embeddingRestEndpoint string = "embedding"
	embeddingIDParam      string = "embedding_id"
	embeddingNameParam    string = "embedding_name"

	limitParam string = "limit"
	pageParam  string = "page"
)

// Ping ... replies to a ping message for healthcheck purposes
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

var router = gin.Default()

// StartEndpoint ... handles requests for the endpoint on the specified port
func StartEndpoint(cfg *conf.Config) {
	// https://github.com/gin-contrib/cors
	// allow all origins
	router.Use(cors.Default())

	// init service
	embeddingService.Init(cfg)

	// add an healthcheck for the endpoint
	router.GET(fmt.Sprintf("healthcheck/%s", embeddingRestEndpoint), Ping)

	// run router as standalone service
	// todo: do we need to run multiple endpoints from the main?
	router.Run(fmt.Sprintf(":%s", cfg.Details["port"]))
}
