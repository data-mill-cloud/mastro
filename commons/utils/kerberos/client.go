package kerberos

import (
	"github.com/datamillcloud/mastro/commons/utils/conf"

	krbClient "github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/keytab"
)

// GetKerberosClient ... returns a gokrb5 kerberos client
func GetKerberosClient(details *conf.KerberosDetails) *krbClient.Client {
	// https://github.com/jcmturner/gokrb5/blob/master/v8/USAGE.md
	// Replace with a valid credentialed client.
	cfg, err := config.Load(details.KrbConfigPath)
	if err != nil {
		panic(err)
	}

	var krb5Client *krbClient.Client

	if len(details.KeytabPath) > 0 {
		kt, err := keytab.Load(details.KeytabPath)
		if err != nil {
			panic(err)
		}
		krb5Client = krbClient.NewWithKeytab(
			details.Username,
			details.Realm, kt, cfg,
			krbClient.DisablePAFXFAST(details.DisablePAFXFAST), krbClient.AssumePreAuthentication(false),
		)
	} else {
		krb5Client = krbClient.NewWithPassword(
			details.Username,
			details.Realm, details.Password, cfg,
			krbClient.DisablePAFXFAST(true),
		)
	}

	return krb5Client
}
