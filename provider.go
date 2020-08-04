package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider returns the definition of this provider
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"vaultutility_initialization": DefineInitializationResource(),
			"vaultutility_unseal":         DefineUnsealResource(),
		},
	}
}
