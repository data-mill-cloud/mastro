package abstract

import (
	"fmt"
	"log"
	"strings"

	"github.com/data-mill-cloud/mastro/commons/utils/conf"
)

// ConnectorProvider ... The interface each connector must implement
type ConnectorProvider interface {
	ValidateDataSourceDefinition(*conf.DataSourceDefinition) error
	InitConnection(*conf.DataSourceDefinition)
	CloseConnection()
}

type ConfigurableConnector struct {
	RequiredFields map[string]string
	OptionalFields map[string]string
}

// ValidateDataSourceDefinition ... validates the provided data source definition
func (c *ConfigurableConnector) ValidateDataSourceDefinition(def *conf.DataSourceDefinition) error {
	// check all required fields are available
	var missingFields []string
	for _, reqvalue := range c.RequiredFields {
		if _, exist := def.Settings[reqvalue]; !exist {
			missingFields = append(missingFields, reqvalue)
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("the following %d fields are missing from the data source configuration: %s", len(missingFields), strings.Join(missingFields[:], ","))
	}

	log.Println("Successfully validated data source definition")
	return nil
}
