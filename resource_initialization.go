package main

import (
	backoff "github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/vault/api"
)

// DefineInitializationResource returns the definition of the initialization resource
func DefineInitializationResource() *schema.Resource {
	return &schema.Resource{
		// Creation means to initialize Vault
		Create: initializeVault,

		// These are not meaningful operations
		Read:   noop,
		Update: nil,

		// Once initialized, Vault cannot transition to uninitialized
		Delete: forget,
		Schema: defineSchema(),
	}
}

func defineSchema() map[string]*schema.Schema {
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
		"secret_shares": {
			Type:     schema.TypeInt,
			Default:  5,
			Optional: true,
			ForceNew: true,
		},
		"secret_threshold": {
			Type:     schema.TypeInt,
			Default:  3,
			Optional: true,
			ForceNew: true,
		},
		"recovery_shares": {
			Type:     schema.TypeInt,
			Default:  5,
			Optional: true,
			ForceNew: true,
		},
		"recovery_threshold": {
			Type:     schema.TypeInt,
			Default:  3,
			Optional: true,
			ForceNew: true,
		},
		"root_token_pgp_key": {
			Type:     schema.TypeString,
			Default:  "",
			Optional: true,
			ForceNew: true,
		},
		"pgp_keys": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			ForceNew: true,
		},
		"recovery_pgp_keys": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			ForceNew: true,
		},
		"keys": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Computed: true,
		},
		"keys_b64": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Computed: true,
		},
		"recovery_keys": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Computed: true,
		},
		"recovery_keys_b64": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Computed: true,
		},
		"root_token": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"google_secretmanager": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"credentials_file": {
						Type: schema.TypeString,
					},
					"project_id": {
						Type:     schema.TypeString,
						Required: true,
					},
					"secret_id": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
			Optional:   true,
			ConfigMode: schema.SchemaConfigModeBlock,
		},
	}
}

func initializeVault(data *schema.ResourceData, meta interface{}) error {
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

	var initialized bool

	err = backoff.Retry(func() error {
		init, initStatusError := client.Sys().InitStatus()
		if initStatusError != nil {
			return initStatusError
		}
		initialized = init
		return nil
	}, backoff.NewExponentialBackOff())

	if err != nil {
		return err
	}

	if initialized {
		data.SetId("unknown cluster id")
		return nil
	}

	return backoff.Retry(func() error {
		initialization, initError := client.Sys().Init(&api.InitRequest{
			SecretShares:      data.Get("secret_shares").(int),
			SecretThreshold:   data.Get("secret_threshold").(int),
			PGPKeys:           expandStringList(data.Get("pgp_keys").([]interface{})),
			RecoveryShares:    data.Get("recovery_shares").(int),
			RecoveryThreshold: data.Get("recovery_threshold").(int),
			RecoveryPGPKeys:   expandStringList(data.Get("recovery_pgp_keys").([]interface{})),
			RootTokenPGPKey:   data.Get("root_token_pgp_key").(string),
		})

		if initError != nil {
			return initError
		}

		health, healthError := client.Sys().Health()
		if healthError != nil {
			return healthError
		}

		if health.ClusterID != "" {
			data.SetId(health.ClusterID)
		} else {
			data.SetId("unknown cluster id")
		}

		return StoreSecrets(initialization, data)
	}, backoff.NewExponentialBackOff())
}
