# hdns_zone

Provides a Hetzner DNS Zone resource. This can be used to create, modify, and delete zones. 

## Example Usage

```hcl-terraform
resource "hdns_zone" "example" {
    name = "example.org"
}
```
 
## Argument Reference

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
