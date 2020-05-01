# Terraform Provider for the Hetzner DNS

Terraform provider for managing zones and records in the [Hetzner DNS service](https://dns.hetzner.com/).

## Installation

1. Download archive with the latest version of provider for your operating system from
   [Gitlab releases page](https://gitlab.com/alxrem/terraform-provider-hdns/-/releases).
2. Unpack provider to `$HOME/.terraform.d/plugins`, i.e.
   ```
   unzip terraform-provider-hdns_vX.Y.Z-linux-amd64.zip -d $HOME/.terraform.d/plugins/
   ```
3. Init your terraform project
   ```
   terraform init
   ```
4. Download from the [releases page](https://gitlab.com/alxrem/terraform-provider-hdns/-/releases) the file
   containing checksum of downloaded version of plugin and compare it with checksum contained in the file
   `.terraform/plugins/<os_arch>/lock.json` in the root of your terraform project.

## Hetzner DNS Provider

You should [obtain API token](https://dns.hetzner.com/settings/api-token) to access Hetzner DNS API. 

### Argument Reference

The following arguments are supported in the provider block:

* `token` &mdash; (Required) This is the Hetzner Cloud API Token, can also be specified with the `HDNS_TOKEN`
   environment variable.

## hdns_zone

Provides a Hetzner DNS Zone resource. This can be used to create, modify, and delete zones. 

## Example Usage

```hcl-terraform
resource "hdns_zone" "example" {
    name = "example.org"
}
```
 
### Argument Reference

The following arguments are supported:

* `name` &mdash; (Required, string) Name of zone.
* `ttl` &mdash; (Optional, int) TTL of zone. Default value is `86400`.

### Attributes Reference

The following attributes are exported:

* `id` &mdash; (string) ID of zone.
* `is_secondary_dns` &mdash; (boolean) Indicates if a zone is a secondary DNS zone.
* `legacy_dns_host` &mdash; (string)
* `legacy_ns` &mdash; (Array of strings)
* `name` &mdash; (string) Name of zone.
* `ns` &mdash; (Array of strings)
* `owner` &mdash; (string) Owner of zone.
* `paused` &mdash; (boolean)
* `permission` &mdash; (string) Zone's permissions.
* `project` &mdash; (string)
* `records_count` &mdash; (int) Amount of records associated to this zone.
* `registrar` &mdash; (string)
* `status` &mdash; (string) Status of zone. `verified`, `failed` or `pending`.
* `ttl` &mdash; (int) TTL of zone.

## hdns_record

Provides a Hetzner DNS Record resource. This can be used to create, modify, and delete records. 

### Example usage

```hcl-terraform
resource "hdns_zone" "example" {
    name = "example.org"
}

resource "hdns_record" "www" {
    zone_id = hdns_zone.example.id

    name  = "www"
    type  = "A"
    value = "127.0.0.1"
}
```

### Argument reference

The following arguments are supported:

* `name` &mdash; (Required, string) Name of record.
* `ttl` &mdash; (Optional, int) TTL of record.
* `type` &mdash; (Required, string) Type of the record. `A`, `AAAA`, `NS`, `MX`, `CNAME`, `RP`, `TXT`, `SOA`, `HINFO`,
  `SRV`, `DANE`, `TLSA`, `DS`, `CAA`.
* `value` &mdash; (Required, string) Value of record (e.g. `127.0.0.1`, `1.1.1.1`).
* `zone_id` &mdash; (Required, string) ID of zone this record is associated with.

### Attributes Reference

The following attributes are exported:

* `name` &mdash; (string) Name of record.
* `ttl` &mdash; (int) TTL of record.
* `type` &mdash; (string) Type of the record.
* `value` &mdash; (string) Value of record.
* `zone_id` &mdash; (string) ID of zone this record is associated with.
