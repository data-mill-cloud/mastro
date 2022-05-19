package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/data-mill-cloud/mastro/commons/utils/errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	embeddingRestEndpoint string = "embedding"
	embeddingIDParam      string = "embedding_id"
	partIDParam
	embeddingNameParam string = "embedding_name"

	limitParam string = "limit"
	pageParam  string = "page"
)

func getLimitAndPageNumber(req *http.Request) (limit int, page int, err error) {
	if limit, err = strconv.Atoi(req.URL.Query().Get(limitParam)); err != nil {
		err = fmt.Errorf(fmt.Sprintf("%s parameter is not a valid integer number", limitParam))
		return
	}
	if page, err = strconv.Atoi(req.URL.Query().Get(pageParam)); err != nil {
		err = fmt.Errorf(fmt.Sprintf("%s parameter is not a valid integer number", pageParam))
	}
	return
}

// Ping ... replies to a ping message for healthcheck purposes
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// GetEmbeddingByID ... retrieves an embedding by the provided ID
func GetEmbeddingByID(c *gin.Context) {
	id := c.Param(embeddingIDParam)
	//partitions := c.Param(partIDParam)
	partitions := []string{}
	fs, getErr := embeddingService.GetEmbeddingByID(id, partitions)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
	} else {
		c.JSON(http.StatusOK, fs)
	}
}

// GetEmbeddingByName ... retrieves embeddings by the provided Name
func GetEmbeddingByName(c *gin.Context) {
	name := c.Param(embeddingNameParam)

	limit, page, err := getLimitAndPageNumber(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	fs, getErr := embeddingService.GetEmbeddingByName(name, limit, page)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
	} else {
		c.JSON(http.StatusOK, fs)
	}
}

// CreateEmbedding ... creates an embedding
func CreateEmbedding(c *gin.Context) {
	e := abstract.Embedding{}
	if err := c.ShouldBindJSON(&e); err != nil {
		restErr := errors.GetBadRequestError("Invalid JSON Body")
		c.JSON(restErr.Status, restErr)
	} else {
		// call service to add the embedding
		result, saveErr := embeddingService.CreateEmbedding(e)
		if saveErr != nil {
			c.JSON(saveErr.Status, saveErr)
		} else {
			c.JSON(http.StatusCreated, result)
		}
	}
}

// SimilarToThisId ... retrieves similar embeddings to the provided id
func SimilarToThisId(c *gin.Context) {
	id := c.Param(embeddingIDParam)

	limit, page, err := getLimitAndPageNumber(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	em, getErr := embeddingService.SimilarToThisId(id, limit, page)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
	} else {
		c.JSON(http.StatusOK, em)
	}
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

	/// --------------------------------

	// get feature set as embedding/id/:fs_id with :fs_id being a placeholder for the value passed
	router.GET(fmt.Sprintf("%s/id/:%s", embeddingRestEndpoint, embeddingIDParam), GetEmbeddingByID)
	// get feature set as embedding/name/:fs_name with :fs_name being a placeholder for the value passed
	router.GET(fmt.Sprintf("%s/name/:%s", embeddingRestEndpoint, embeddingNameParam), GetEmbeddingByName)

	// search by query string
	router.POST(fmt.Sprintf("%s/similarid", embeddingRestEndpoint), SimilarToThisId)

	// put embedding
	router.PUT(fmt.Sprintf("%s/", embeddingRestEndpoint), CreateEmbedding)

	////////////////////////////////

	// run router as standalone service
	// todo: do we need to run multiple endpoints from the main?
	router.Run(fmt.Sprintf(":%s", cfg.Details["port"]))
}
