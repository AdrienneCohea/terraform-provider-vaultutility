terraform {
  required_version = "~> 1.0"

  required_providers {
    vaultutility = {
      source  = "localhost/providers/vaultutility"
      version = "~> 0.1"
    }
  }
}
