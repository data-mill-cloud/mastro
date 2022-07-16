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
	featureSetRestEndpoint string = "featureset"
	featureSetIDParam      string = "featureset_id"
	featureSetNameParam    string = "featureset_name"

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

// CreateFeatureSet ... creates a featureSet
func CreateFeatureSet(c *gin.Context) {
	fs := abstract.FeatureSet{}
	if err := c.ShouldBindJSON(&fs); err != nil {
		restErr := errors.GetBadRequestError("Invalid JSON Body")
		c.JSON(restErr.Status, restErr)
	} else {
		// call service to add the featureset
		result, saveErr := featureStoreService.CreateFeatureSet(fs)
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
	fs, getErr := featureStoreService.GetFeatureSetByID(id)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
	} else {
		c.JSON(http.StatusOK, fs)
	}
}

// GetFeatureSetByName ... retrieves a featureSet by the provided Name
func GetFeatureSetByName(c *gin.Context) {
	//id, err := parseFeatureSetName(c.Param(featureSetNameParam))
	name := c.Param(featureSetNameParam)

	limit, page, err := getLimitAndPageNumber(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.GetBadRequestError(err.Error()))
		return
	}

	fs, getErr := featureStoreService.GetFeatureSetByName(name, limit, page)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
	} else {
		c.JSON(http.StatusOK, fs)
	}
}

// SearchFeatureSetsByLabels ... retrieves any featureset matching all specified labels or error if empty
func SearchFeatureSetsByLabels(c *gin.Context) {
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
			fsets, getErr := featureStoreService.SearchFeatureSetsByLabels(query.Labels, query.Limit, query.Page)
			if getErr != nil {
				c.JSON(getErr.Status, getErr)
			} else {
				c.JSON(http.StatusOK, fsets)
			}
		}
	}
}

func SearchFeatureSetsByQueryLabels(c *gin.Context) {
	limit, page, err := getLimitAndPageNumber(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.GetBadRequestError(err.Error()))
		return
	}

	query := c.Request.URL.Query()

	if len(query) == 0 {
		restErr := errors.GetBadRequestError("Invalid query by labels :: empty label dict")
		c.JSON(restErr.Status, restErr)
	} else {
		q := make(map[string]string)
		for k, l := range query {
			if k != limitParam && k != pageParam {
				q[k] = l[0]
			}
		}
		fsets, getErr := featureStoreService.SearchFeatureSetsByLabels(q, limit, page)
		if getErr != nil {
			c.JSON(getErr.Status, getErr)
		} else {
			c.JSON(http.StatusOK, fsets)
		}

	}
}

// Search ... search by a full text query param
func Search(c *gin.Context) {
	query := queries.ByText{}
	err := c.BindJSON(&query)
	if err != nil {
		restErr := errors.GetBadRequestError("Invalid text query :: invalid input json format")
		c.JSON(restErr.Status, restErr)
	} else {
		if len(query.Query) == 0 {
			restErr := errors.GetBadRequestError("Invalid text query :: empty text")
			c.JSON(restErr.Status, restErr)
		} else {
			fsets, getErr := featureStoreService.Search(query.Query, query.Limit, query.Page)
			if getErr != nil {
				c.JSON(getErr.Status, getErr)
			} else {
				c.JSON(http.StatusOK, fsets)
			}
		}
	}

}

// ListAllFeatureSets ... lists all featuresets in the DB
func ListAllFeatureSets(c *gin.Context) {

	limit, page, err := getLimitAndPageNumber(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.GetBadRequestError(err.Error()))
		return
	}

	fsets, svcErr := featureStoreService.ListAllFeatureSets(limit, page)
	if svcErr != nil {
		c.JSON(svcErr.Status, svcErr)
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
	featureStoreService.Init(cfg)

	// add an healthcheck for the endpoint
	router.GET(fmt.Sprintf("healthcheck/%s", featureSetRestEndpoint), Ping)

	// get feature set as featureset/id/:fs_id with :fs_id being a placeholder for the value passed
	router.GET(fmt.Sprintf("%s/id/:%s", featureSetRestEndpoint, featureSetIDParam), GetFeatureSetByID)
	// get feature set as featureset/name/:fs_name with :fs_name being a placeholder for the value passed
	router.GET(fmt.Sprintf("%s/name/:%s", featureSetRestEndpoint, featureSetNameParam), GetFeatureSetByName)

	// search by query string
	router.POST(fmt.Sprintf("%s/search", featureSetRestEndpoint), Search)

	// put feature set as featureset/
	router.PUT(fmt.Sprintf("%s/", featureSetRestEndpoint), CreateFeatureSet)

	router.POST(fmt.Sprintf("%s/labels", featureSetRestEndpoint), SearchFeatureSetsByLabels)
	router.GET(fmt.Sprintf("%s/labels", featureSetRestEndpoint), SearchFeatureSetsByQueryLabels)

	// list all feature sets
	router.GET(fmt.Sprintf("%s/", featureSetRestEndpoint), ListAllFeatureSets)

	// run router as standalone service
	// todo: do we need to run multiple endpoints from the main?
	router.Run(fmt.Sprintf(":%s", cfg.Details["port"]))
}
