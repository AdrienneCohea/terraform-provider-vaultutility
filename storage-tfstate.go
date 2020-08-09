package main

import (
	"context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/vault/api"
	"google.golang.org/api/option"
)

func StoreSecrets(init *api.InitResponse, data *schema.ResourceData) error {
	googleStorageConfig, exists := data.GetOk("google_secrets_manager")
	if exists {
		ctx := context.Background()
		configMap := googleStorageConfig.(map[string]interface{})
		client, err := secretmanager.NewClient(ctx, option.WithCredentialsFile(configMap["credentials_file"].(string)))
		if err != nil {
			return err
		}

		return StoreInGoogleSecretsManager(ctx, client, configMap["project_id"].(string), configMap["secret_id"].(string), init)
	}

	data.Set("keys", init.Keys)
	data.Set("keys_b64", init.KeysB64)
	data.Set("recovery_keys", init.RecoveryKeys)
	data.Set("recovery_keys_b64", init.RecoveryKeysB64)
	data.Set("root_token", init.RootToken)

	return nil
}
