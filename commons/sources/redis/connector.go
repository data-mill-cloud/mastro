package redis

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/go-redis/redis/v8"
)

var requiredFields = map[string]string{
	"redisUser": "username",
	"redisPwd":  "password",
	"redisHost": "host",
	"redisDb":   "db",
}

var optionalFields = map[string]string{
	"cert": "cert",
}

// NewRedisConnector factory
func NewRedisConnector() *Connector {
	return &Connector{}
}

// Connector ... Redis connector
type Connector struct {
	Client *redis.Client
}

// ValidateDataSourceDefinition ... Validates the input data source definition
func (c *Connector) ValidateDataSourceDefinition(def *conf.DataSourceDefinition) error {
	// check all required fields are available
	var missingFields []string
	for _, reqvalue := range requiredFields {
		if _, exist := def.Settings[reqvalue]; !exist {
			missingFields = append(missingFields, reqvalue)
		}
	}

	if len(missingFields) > 0 {
		// https://stackoverflow.com/questions/28799110/how-to-join-a-slice-of-strings-into-a-single-string
		return fmt.Errorf("The following fields are missing from the data source configuration: %s", strings.Join(missingFields, ","))
	}

	_, err := strconv.Atoi(def.Settings[requiredFields["redisDB"]])
	if err != nil {
		return fmt.Errorf("Impossible to convert redisDB to integer")
	}

	log.Println("Successfully validated data source definition")
	return nil
}

// InitConnection ... inits connection
func (c *Connector) InitConnection(def *conf.DataSourceDefinition) {
	redisConf := &redis.Options{}

	redisHost := def.Settings[requiredFields["redisHost"]]
	//stringutils.SplitAndTrim(def.Settings.Values[requiredFields["redisHost"]], ",")

	redisConf.Addr = redisHost
	redisConf.Username = def.Settings[requiredFields["redisUser"]]
	redisConf.Password = def.Settings[requiredFields["redisPwd"]]
	// assuming we did already validate the conversion to int
	redisConf.DB, _ = strconv.Atoi(def.Settings[requiredFields["redisDb"]])

	c.Client = redis.NewClient(redisConf)
}

// CloseConnection ... terminates the connection
func (c *Connector) CloseConnection() {
	c.CloseConnection()
}
