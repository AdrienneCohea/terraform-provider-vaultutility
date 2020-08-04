# vaultutility_initialization Resource

Initializes a Vault and stores its unseal keys, recovery keys, and root token in the state.

## Example Usage

```hcl
resource "vaultutility_initialization" "create" {
  address = "https://vault.example.com:8200"
}
```

## Argument Reference

* `address` - (Required) The Vault address

## Attribute Reference

* `keys` - Unseal keys in hexadecimal representation
* `keys_b64` - Unseal keys in base64 representation
* `recovery_keys` - Recovery keys in hexadecimal representation
* `recovery_keys_b64` - Recovery keys in base64 representation
* `root_token` - Root token
