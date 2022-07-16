package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/data-mill-cloud/mastro/commons/utils/errors"
	"github.com/data-mill-cloud/mastro/commons/utils/queries"
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
	fs, getErr := embeddingService.GetEmbeddingByID(id)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
	} else {
		c.JSON(http.StatusOK, fs)
	}
}

// GetEmbeddingByName ... retrieves embeddings by the provided Name
func GetEmbeddingByName(c *gin.Context) {
	name := c.Param(embeddingNameParam)

	fs, getErr := embeddingService.GetEmbeddingByName(name)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
	} else {
		c.JSON(http.StatusOK, fs)
	}
}

// Upsert ... upsert embeddings
func UpsertEmbeddings(c *gin.Context) {
	embeddings := []abstract.Embedding{}
	if err := c.ShouldBindJSON(&embeddings); err != nil {
		restErr := errors.GetBadRequestError("Invalid JSON Body")
		c.JSON(restErr.Status, restErr)
	} else {
		// call service to add the embedding
		saveErr := embeddingService.UpsertEmbeddings(embeddings)
		if saveErr != nil {
			c.JSON(saveErr.Status, saveErr)
		} else {
			c.Writer.WriteHeader(http.StatusCreated)
		}
	}
}

// SimilarToThis ... retrieves similar embeddings to the provided one
func SimilarToThis(c *gin.Context) {
	query := queries.ByVector{}
	err := c.BindJSON(&query)

	if err != nil {
		restErr := errors.GetBadRequestError("Invalid query by vector :: invalid input json format")
		c.JSON(restErr.Status, restErr)
	} else {
		if query.Vector == nil || len(query.Vector) == 0 {
			restErr := errors.GetBadRequestError("Invalid query by vector :: missing or empty vector embedding")
			c.JSON(restErr.Status, restErr)
		} else {
			em, getErr := embeddingService.SimilarToThis(query.Vector, query.K)
			if getErr != nil {
				c.JSON(getErr.Status, getErr)
			} else {
				c.JSON(http.StatusOK, em)
			}
		}
	}
}

func DeleteEmbeddingByID(c *gin.Context) {
	id := c.Param(embeddingIDParam)
	getErr := embeddingService.DeleteEmbeddingByIds(id)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
	} else {
		c.Writer.WriteHeader(http.StatusOK)
	}
}

func DeleteEmbeddingByName(c *gin.Context) {
	name := c.Param(embeddingNameParam)
	getErr := embeddingService.DeleteEmbeddingByName(name)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
	} else {
		c.Writer.WriteHeader(http.StatusOK)
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
	router.POST(fmt.Sprintf("%s/similar", embeddingRestEndpoint), SimilarToThis)

	// put embedding
	router.PUT(fmt.Sprintf("%s/", embeddingRestEndpoint), UpsertEmbeddings)

	router.DELETE(fmt.Sprintf("%s/id/:%s", embeddingRestEndpoint, embeddingIDParam), DeleteEmbeddingByID)
	router.DELETE(fmt.Sprintf("%s/name/:%s", embeddingRestEndpoint, embeddingNameParam), DeleteEmbeddingByName)

	////////////////////////////////

	// run router as standalone service
	// todo: do we need to run multiple endpoints from the main?
	router.Run(fmt.Sprintf(":%s", cfg.Details["port"]))
}
