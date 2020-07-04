# Data Source: hdns_zone

Provides details about Hetzner DNS Zone.

## Example usage

```hcl-terraform
data "hdns_zone" "example" {
  id = "tv7gT6IoAdvvoCBPkNwyAc"
}
```

## Argument Reference

The following arguments are supported:

* `id` &mdash; (Required, string) ID of zone.

## Attributes Reference

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
