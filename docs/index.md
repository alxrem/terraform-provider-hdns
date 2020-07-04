# Hetzner DNS Provider

Terraform provider for managing zones and records in the [Hetzner DNS service](https://dns.hetzner.com/).

## Example Usage

```hcl-terraform
provider "hdns" {
  token = var.hdns_token
}

resource "hdns_zone" "example" {
    name = "example.org"
}
```

## Argument Reference

The following arguments are supported in the provider block:

* `token` &mdash; (Required) This is the Hetzner Cloud API Token, can also be specified with the `HDNS_TOKEN`
   environment variable.
