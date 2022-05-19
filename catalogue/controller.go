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
	assetsRestEndpoint string = "assets"
	assetRestEndpoint  string = "asset"
	// placeholders for the values actually passed to the endpoint
	assetIDParam   string = "asset_id"
	assetNameParam string = "asset_name"

	limitParam string = "limit"
	pageParam  string = "page"
)

// Ping ... replies to a ping message for healthcheck purposes
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// UpsertAsset ... creates an asset description entry
func UpsertAsset(c *gin.Context) {
	asset := abstract.Asset{}
	if err := c.ShouldBindJSON(&asset); err != nil {
		restErr := errors.GetBadRequestError("Invalid JSON Body")
		c.JSON(restErr.Status, restErr)
	} else {
		result, saveErr := catalogueService.UpsertAssets(&[]abstract.Asset{asset})
		if saveErr != nil {
			c.JSON(saveErr.Status, saveErr)
		} else {
			c.JSON(http.StatusCreated, result)
		}
	}
}

// BulkUpsert ... bulk upsert
func BulkUpsert(c *gin.Context) {
	assets := []abstract.Asset{}
	if err := c.ShouldBindJSON(&assets); err != nil {
		restErr := errors.GetBadRequestError("Invalid JSON Body")
		c.JSON(restErr.Status, restErr)
	} else {
		result, saveErr := catalogueService.UpsertAssets(&assets)
		if saveErr != nil {
			c.JSON(saveErr.Status, saveErr)
		} else {
			c.JSON(http.StatusCreated, result)
		}
	}
}

// GetAssetByID ... retrieves an asset description by its Unique Name ID
func GetAssetByID(c *gin.Context) {
	nameID := c.Param(assetIDParam)
	asset, getErr := catalogueService.GetAssetByID(nameID)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
	} else {
		c.JSON(http.StatusOK, asset)
	}
}

// GetAssetByName ... retrieves an asset description by its Unique Name
func GetAssetByName(c *gin.Context) {
	nameID := c.Param(assetNameParam)
	asset, getErr := catalogueService.GetAssetByName(nameID)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
	} else {
		c.JSON(http.StatusOK, asset)
	}
}

// SearchAssetsByTags ... retrieves any asset matching all specified tags or error if empty
func SearchAssetsByTags(c *gin.Context) {
	/*
		limit, err := strconv.ParseInt(c.Request.URL.Query().Get("limit"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		page, err := strconv.ParseInt(c.Request.URL.Query().Get("page"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
	*/
	query := queries.ByTags{}
	err := c.BindJSON(&query)
	if err != nil {
		restErr := errors.GetBadRequestError("Invalid query by tag :: invalid input json format")
		c.JSON(restErr.Status, restErr)
	} else {
		if query.Tags == nil || len(query.Tags) == 0 {
			restErr := errors.GetBadRequestError("Invalid query by tag :: empty tag list")
			c.JSON(restErr.Status, restErr)
		} else {
			assets, getErr := catalogueService.SearchAssetsByTags(query.Tags,
				//limit,
				query.Limit,
				//page,
				query.Page,
			)
			if getErr != nil {
				c.JSON(getErr.Status, getErr)
			} else {
				c.JSON(http.StatusOK, assets)
			}
		}
	}
}

// ListAllAssets ... returns all assets
func ListAllAssets(c *gin.Context) {

	limit, page, err := getLimitAndPageNumber(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	assets, getErr := catalogueService.ListAllAssets(limit, page)
	if err != nil {
		c.JSON(getErr.Status, getErr)
	} else {
		c.JSON(http.StatusOK, assets)
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
			assets, getErr := catalogueService.Search(query.Query, query.Limit, query.Page)
			if getErr != nil {
				c.JSON(getErr.Status, getErr)
			} else {
				c.JSON(http.StatusOK, assets)
			}
		}
	}

}

func getLimitAndPageNumber(req *http.Request) (limit int, page int, err error) {
	if limit, err = strconv.Atoi(req.URL.Query().Get(limitParam)); err != nil {
		return
	}
	page, err = strconv.Atoi(req.URL.Query().Get(pageParam))
	return
}

var router = gin.Default()

// StartEndpoint ... starts the service endpoint
func StartEndpoint(cfg *conf.Config) {
	// https://github.com/gin-contrib/cors
	// allow all origins
	router.Use(cors.Default())

	// init service
	catalogueService.Init(cfg)

	// add an healthcheck for the endpoint
	router.GET(fmt.Sprintf("healthcheck/%s", assetRestEndpoint), Ping)

	// get specific asset as asset/:id or asset/:name
	router.GET(fmt.Sprintf("%s/id/:%s", assetRestEndpoint, assetIDParam), GetAssetByID)
	router.GET(fmt.Sprintf("%s/name/:%s", assetRestEndpoint, assetNameParam), GetAssetByName)

	// put 1 asset as asset/
	router.PUT(fmt.Sprintf("%s/", assetRestEndpoint), UpsertAsset)
	// put n assets as asset/
	router.PUT(fmt.Sprintf("%s/", assetsRestEndpoint), BulkUpsert)

	// get any asset matching tags
	router.POST(fmt.Sprintf("%s/tags", assetsRestEndpoint), SearchAssetsByTags)
	router.POST(fmt.Sprintf("%s/search", assetsRestEndpoint), Search)

	// list all assets
	router.GET(fmt.Sprintf("%s/", assetsRestEndpoint), ListAllAssets)

	// run router as standalone service
	// todo: do we need to run multiple endpoints from the main?
	router.Run(fmt.Sprintf(":%s", cfg.Details["port"]))
}
