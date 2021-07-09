package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/vault/api"
)

func noop(d *schema.ResourceData, m interface{}) error {
	return nil
}

func forget(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}

func expandStringList(list []interface{}) []string {
	vs := make([]string, 0, len(list))
	for _, v := range list {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}

func getVaultClient(address string, caCert string, clientCert string, clientKey string, tlsServerName string) (*api.Client, error) {
	config := api.DefaultConfig()
	config.Address = address

	if caCert != "" || clientCert != "" || clientKey != "" || tlsServerName != "" {
		err := config.ConfigureTLS(&api.TLSConfig{
			CACert:        caCert,
			ClientCert:    clientCert,
			ClientKey:     clientKey,
			TLSServerName: tlsServerName,
		})
		if err != nil {
			return nil, err
		}
	}

	return api.NewClient(config)
}
