package main

import (
	"fmt"
	"net/http"

	"github.com/datamillcloud/mastro/commons/abstract"
	"github.com/datamillcloud/mastro/commons/utils/conf"
	"github.com/datamillcloud/mastro/commons/utils/errors"
	"github.com/datamillcloud/mastro/commons/utils/queries"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	assetsRestEndpoint string = "assets"
	assetRestEndpoint  string = "asset"
	// placeholders for the values actually passed to the endpoint
	assetIDParam   string = "asset_id"
	assetNameParam string = "asset_name"
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
		result, saveErr := assetService.UpsertAssets(&[]abstract.Asset{asset})
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
		result, saveErr := assetService.UpsertAssets(&assets)
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
	asset, getErr := assetService.GetAssetByID(nameID)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
	} else {
		c.JSON(http.StatusOK, asset)
	}
}

// GetAssetByName ... retrieves an asset description by its Unique Name
func GetAssetByName(c *gin.Context) {
	nameID := c.Param(assetNameParam)
	asset, getErr := assetService.GetAssetByName(nameID)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
	} else {
		c.JSON(http.StatusOK, asset)
	}
}

// SearchAssetsByTags ... retrieves any asset matching all specified tags or error if empty
func SearchAssetsByTags(c *gin.Context) {
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
			assets, getErr := assetService.SearchAssetsByTags(query.Tags)
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
	assets, err := assetService.ListAllAssets()
	if err != nil {
		c.JSON(err.Status, err)
	} else {
		c.JSON(http.StatusOK, assets)
	}
}

var router = gin.Default()

// StartEndpoint ... starts the service endpoint
func StartEndpoint(cfg *conf.Config) {
	// https://github.com/gin-contrib/cors
	// allow all origins
	router.Use(cors.Default())

	// init service
	assetService.Init(cfg)

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

	// list all assets
	router.GET(fmt.Sprintf("%s/", assetsRestEndpoint), ListAllAssets)

	// run router as standalone service
	// todo: do we need to run multiple endpoints from the main?
	router.Run(fmt.Sprintf(":%s", cfg.Details["port"]))
}
