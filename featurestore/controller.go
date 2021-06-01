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
	featureSetRestEndpoint string = "featureset"
	featureSetIDParam      string = "featureset_id"
	featureSetNameParam    string = "featureset_name"
)

// Ping ... replies to a ping message for healthcheck purposes
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// CreateFeatureSet ... creates a featureSet
func CreateFeatureSet(c *gin.Context) {
	fs := abstract.FeatureSet{}
	if err := c.ShouldBindJSON(&fs); err != nil {
		restErr := errors.GetBadRequestError("Invalid JSON Body")
		c.JSON(restErr.Status, restErr)
	} else {
		// call service to add the featureset
		result, saveErr := featureSetService.CreateFeatureSet(fs)
		if saveErr != nil {
			c.JSON(saveErr.Status, saveErr)
		} else {
			c.JSON(http.StatusCreated, result)
		}
	}
}

// parseFeatureSetID ... attempts parsing the fs id from a string param
func parseFeatureSetID(param string) (int64, *errors.RestErr) {
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, errors.GetBadRequestError("invalid feature set id, it should be an integer number")
	}
	return id, nil
}

// GetFeatureSetByID ... retrieves a featureSet by the provided ID
func GetFeatureSetByID(c *gin.Context) {
	//id, err := parseFeatureSetID(c.Param(featureSetIDParam))
	id := c.Param(featureSetIDParam)
	/*
		if err != nil {
			c.JSON(err.Status, err)
		} else {
	*/
	fs, getErr := featureSetService.GetFeatureSetByID(id)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
	} else {
		c.JSON(http.StatusOK, fs)
	}
	//}
}

// GetFeatureSetByName ... retrieves a featureSet by the provided Name
func GetFeatureSetByName(c *gin.Context) {
	//id, err := parseFeatureSetName(c.Param(featureSetNameParam))
	name := c.Param(featureSetNameParam)
	/*
		if err != nil {
			c.JSON(err.Status, err)
		} else {
	*/
	fs, getErr := featureSetService.GetFeatureSetByName(name)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
	} else {
		c.JSON(http.StatusOK, fs)
	}
	//}
}

// ListAllFeatureSets ... lists all featuresets in the DB
func ListAllFeatureSets(c *gin.Context) {
	fsets, err := featureSetService.ListAllFeatureSets()
	if err != nil {
		c.JSON(err.Status, err)
	} else {
		c.JSON(http.StatusOK, fsets)
	}
}

var router = gin.Default()

// StartEndpoint ... handles requests for the endpoint on the specified port
func StartEndpoint(cfg *conf.Config) {
	// https://github.com/gin-contrib/cors
	// allow all origins
	router.Use(cors.Default())

	// init service
	featureSetService.Init(cfg)

	// add an healthcheck for the endpoint
	router.GET(fmt.Sprintf("healthcheck/%s", featureSetRestEndpoint), Ping)

	// get feature set as featureset/id/:fs_id with :fs_id being a placeholder for the value passed
	router.GET(fmt.Sprintf("%s/id/:%s", featureSetRestEndpoint, featureSetIDParam), GetFeatureSetByID)
	// get feature set as featureset/name/:fs_name with :fs_name being a placeholder for the value passed
	router.GET(fmt.Sprintf("%s/name/:%s", featureSetRestEndpoint, featureSetNameParam), GetFeatureSetByName)

	// put feature set as featureset/
	router.PUT(fmt.Sprintf("%s/", featureSetRestEndpoint), CreateFeatureSet)

	// list all feature sets
	router.GET(fmt.Sprintf("%s/", featureSetRestEndpoint), ListAllFeatureSets)

	// run router as standalone service
	// todo: do we need to run multiple endpoints from the main?
	router.Run(fmt.Sprintf(":%s", cfg.Details["port"]))
}
