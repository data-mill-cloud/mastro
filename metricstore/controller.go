package main

import (
	"fmt"
	"net/http"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/data-mill-cloud/mastro/commons/utils/errors"
	"github.com/data-mill-cloud/mastro/commons/utils/queries"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	metricStoreRestEndpoint string = "metricstore"
	metricSetIDParam        string = "metricset_id"
	metricSetNameParam      string = "metricset_name"
)

// Ping ... replies to a ping message for healthcheck purposes
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// CreateMetricSet ... creates a metricSet
func CreateMetricSet(c *gin.Context) {
	fs := abstract.MetricSet{}
	if err := c.ShouldBindJSON(&fs); err != nil {
		restErr := errors.GetBadRequestError("Invalid JSON Body")
		c.JSON(restErr.Status, restErr)
	} else {
		// call service to add the metricset
		result, saveErr := metricStoreService.CreateMetricSet(fs)
		if saveErr != nil {
			c.JSON(saveErr.Status, saveErr)
		} else {
			c.JSON(http.StatusCreated, result)
		}
	}
}

// GetMetricSetByID ... retrieves a metricSet by the provided ID
func GetMetricSetByID(c *gin.Context) {
	id := c.Param(metricSetIDParam)

	ms, getErr := metricStoreService.GetMetricSetByID(id)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
	} else {
		c.JSON(http.StatusOK, ms)
	}
}

// GetMetricSetByName ... retrieves a metricSet by the provided Name
func GetMetricSetByName(c *gin.Context) {
	//id, err := parseMetricSetName(c.Param(metricSetNameParam))
	name := c.Param(metricSetNameParam)

	ms, getErr := metricStoreService.GetMetricSetByName(name)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
	} else {
		c.JSON(http.StatusOK, ms)
	}
}

// SearchMetricSetsByLabels ... retrieves any metricset matching all specified labels or error if empty
func SearchMetricSetsByLabels(c *gin.Context) {
	query := queries.ByLabels{}
	err := c.BindJSON(&query)

	if err != nil {
		restErr := errors.GetBadRequestError("Invalid query by labels :: invalid input json format")
		c.JSON(restErr.Status, restErr)
	} else {
		if query.Labels == nil || len(query.Labels) == 0 {
			restErr := errors.GetBadRequestError("Invalid query by labels :: empty label dict")
			c.JSON(restErr.Status, restErr)
		} else {
			assets, getErr := metricStoreService.SearchMetricSetsByLabels(query.Labels)
			if getErr != nil {
				c.JSON(getErr.Status, getErr)
			} else {
				c.JSON(http.StatusOK, assets)
			}
		}
	}
}

// ListAllMetricSets ... lists all metricsets in the DB
func ListAllMetricSets(c *gin.Context) {
	msets, err := metricStoreService.ListAllMetricSets()
	if err != nil {
		c.JSON(err.Status, err)
	} else {
		c.JSON(http.StatusOK, msets)
	}
}

var router = gin.Default()

// StartEndpoint ... handles requests for the endpoint on the specified port
func StartEndpoint(cfg *conf.Config) {
	// https://github.com/gin-contrib/cors
	// allow all origins
	router.Use(cors.Default())

	// init service
	metricStoreService.Init(cfg)

	// add an healthcheck for the endpoint
	router.GET(fmt.Sprintf("healthcheck/%s", metricStoreRestEndpoint), Ping)

	// get metric set as metricset/id/:fs_id with :fs_id being a placeholder for the value passed
	router.GET(fmt.Sprintf("%s/id/:%s", metricStoreRestEndpoint, metricSetIDParam), GetMetricSetByID)
	// get metric set as metricset/name/:fs_name with :fs_name being a placeholder for the value passed
	router.GET(fmt.Sprintf("%s/name/:%s", metricStoreRestEndpoint, metricSetNameParam), GetMetricSetByName)

	// put metricset as metricset/
	router.PUT(fmt.Sprintf("%s/", metricStoreRestEndpoint), CreateMetricSet)

	// get any metricset matching labels
	router.POST(fmt.Sprintf("%s/labels", metricStoreRestEndpoint), SearchMetricSetsByLabels)

	// list all metricsets
	router.GET(fmt.Sprintf("%s/", metricStoreRestEndpoint), ListAllMetricSets)

	// run router as standalone service
	// todo: do we need to run multiple endpoints from the main?
	router.Run(fmt.Sprintf(":%s", cfg.Details["port"]))
}
