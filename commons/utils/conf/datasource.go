package conf

// DataSourceDefinition ... connection details for a data source connector
type DataSourceDefinition struct {
	Name              string            `yaml:"name"`
	Type              string            `yaml:"type"`
	CrawlerDefinition CrawlerDefinition `yaml:"crawler,omitempty"`
	Settings          map[string]string `yaml:"settings,omitempty"`
	// optional kerberos section
	KerberosDetails *KerberosDetails `yaml:"kerberos"`
	// optional tls section
	TLSDetails *TLSDetails `yaml:"tls"`
}

// KerberosDetails ... Connection details for Kerberos
type KerberosDetails struct {
	KrbConfigPath   string `yaml:"krb-config-path,omitempty"`
	SASLMechanism   string `yaml:"sasl-mech,omitempty"`
	EnableSASL      bool   `yaml:"enable-sasl,omitempty"`
	ServiceName     string `yaml:"service-name,omitempty"`
	Realm           string `yaml:"realm,omitempty"`
	Username        string `yaml:"username,omitempty"`
	AuthType        int    `yaml:"auth-type,omitempty"`
	Password        string `yaml:"password,omitempty"`
	KeytabPath      string `yaml:"keytab-path,omitempty"`
	DisablePAFXFAST bool   `yaml:"disable-pafx-fast,omitempty"`
}

// TLSDetails ... TLS Connection details
type TLSDetails struct {
	Enable             bool   `yaml:"enable,omitempty"`
	InsecureSkipVerify bool   `yaml:"insecure-skip-verify,omitempty"`
	ClientCertFile     string `yaml:"client-cert-file"`
	ClientKeyFile      string `yaml:"client-key-file"`
	CaCertFile         string `yaml:"ca-cert-file"`
}
