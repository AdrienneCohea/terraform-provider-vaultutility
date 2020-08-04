# vaultutility Provider

Provides utility resources for Vault clusters, including initialization.

## Example Usage

```hcl
resource "vaultutility_initialization" "create" {
  address = "https://vault.example.com:8200"
}
```
