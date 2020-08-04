package main

import (
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// DefineUnsealResource returns the definition of the unseal resource
func DefineUnsealResource() *schema.Resource {
	return &schema.Resource{
		// Creation means to initialize Vault
		Create: unsealVault,

		// These are not meaningful operations
		Read:   noop,
		Update: nil,

		// We do not implement seal
		Delete: forget,
		Schema: defineUnsealSchema(),
	}
}

func defineUnsealSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"ca_cert_file": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"client_key_file": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"client_cert_file": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"unseal_keys": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Required: true,
			ForceNew: true,
		},
	}
}

func unsealVault(data *schema.ResourceData, meta interface{}) error {
	client, err := getVaultClient(
		data.Get("address").(string),
		data.Get("ca_cert_file").(string),
		"",
		"",
		"",
	)
	if err != nil {
		return err
	}

	return backoff.Retry(func() error {
		for _, key := range expandStringList(data.Get("unseal_keys").([]interface{})) {
			unsealResponse, unsealError := client.Sys().Unseal(key)
			if unsealError != nil {
				return unsealError
			}

			if !unsealResponse.Sealed {
				data.SetId(data.Get("address").(string))
				return nil
			}
		}

		return nil
	}, backoff.NewConstantBackOff(time.Second*300))
}
