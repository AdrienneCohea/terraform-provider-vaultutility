package main

import (
	"context"
	"encoding/json"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func StoreInGoogleSecretsManager(ctx context.Context, client *secretmanager.Client, projectId string, secretId string, plainText interface{}) error {
	bytes, err := json.Marshal(plainText)
	if err != nil {
		return err
	}

	secret, err := client.CreateSecret(ctx, &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", projectId),
		SecretId: secretId,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	})

	// Todo: handle error when secret already exists
	if err != nil {
		return err
	}

	_, err = client.AddSecretVersion(ctx, &secretmanagerpb.AddSecretVersionRequest{
		Parent: secret.Name,
		Payload: &secretmanagerpb.SecretPayload{
			Data: bytes,
		},
	})

	return err
}

func RetrieveFromGoogleSecretsManager(ctx context.Context, client *secretmanager.Client, projectId string, secretId string) ([]byte, error) {
	secret, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectId, secretId),
	})
	if err != nil {
		return nil, err
	}

	return secret.Payload.Data, nil
}
