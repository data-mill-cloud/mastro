package redis

import (
	"log"
	"strconv"

	"github.com/data-mill-cloud/mastro/commons/abstract"
	"github.com/data-mill-cloud/mastro/commons/utils/conf"
	"github.com/go-redis/redis/v8"
)

// NewRedisConnector factory
func NewRedisConnector() *Connector {
	return &Connector{
		ConfigurableConnector: abstract.ConfigurableConnector{
			RequiredFields: map[string]string{
				"redisUser": "username",
				"redisPwd":  "password",
				"redisHost": "host",
				"redisDb":   "db",
			},
			OptionalFields: map[string]string{
				"cert": "cert",
			},
		},
	}
}

// Connector ... Redis connector
type Connector struct {
	abstract.ConfigurableConnector
	Client *redis.Client
}

// InitConnection ... inits connection
func (c *Connector) InitConnection(def *conf.DataSourceDefinition) {
	redisConf := &redis.Options{}

	redisHost := def.Settings[c.RequiredFields["redisHost"]]
	//stringutils.SplitAndTrim(def.Settings.Values[requiredFields["redisHost"]], ",")

	redisConf.Addr = redisHost
	redisConf.Username = def.Settings[c.RequiredFields["redisUser"]]
	redisConf.Password = def.Settings[c.RequiredFields["redisPwd"]]

	var err error
	if redisConf.DB, err = strconv.Atoi(def.Settings[c.RequiredFields["redisDb"]]); err != nil {
		log.Fatalln("Impossible to convert redisDB to integer")
	}

	c.Client = redis.NewClient(redisConf)
}

// CloseConnection ... terminates the connection
func (c *Connector) CloseConnection() {
	c.Client.Close()
}
