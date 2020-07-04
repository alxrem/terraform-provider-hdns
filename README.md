# Terraform Provider for the Hetzner DNS

Terraform provider for managing zones and records in the [Hetzner DNS service](https://dns.hetzner.com/).

## Installation

### terraform 0.13+

Add into your Terraform configuration this code:

```hcl-terraform
terraform {
  required_providers {
    hdns = {
      source = "alxrem/hdns"
    }
  }
}
```

and run `terraform init`

### terraform 0.12 and earlier

1. Download archive with the latest version of provider for your operating system from
   [Github releases page](https://github.com/alxrem/terraform-provider-hdns/releases).
2. Unpack provider to `$HOME/.terraform.d/plugins`, i.e.
   ```
   unzip terraform-provider-hdns_X.Y.Z_linux_amd64.zip terraform-provider-hdns_* -d $HOME/.terraform.d/plugins/
   ```
3. Init your terraform project
   ```
   terraform init
   ```

## Usage

Read the [documentation on Terraform Registry site](https://registry.terraform.io/providers/alxrem/hdns/latest/docs).
