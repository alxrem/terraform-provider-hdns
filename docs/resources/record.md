# hdns_record

Provides a Hetzner DNS Record resource. This can be used to create, modify, and delete records. 

## Example usage

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

## Argument reference

The following arguments are supported:

* `name` &mdash; (Required, string) Name of record.
* `ttl` &mdash; (Optional, int) TTL of record.
* `type` &mdash; (Required, string) Type of the record. `A`, `AAAA`, `NS`, `MX`, `CNAME`, `RP`, `TXT`, `SOA`, `HINFO`,
  `SRV`, `DANE`, `TLSA`, `DS`, `CAA`.
* `value` &mdash; (Required, string) Value of record (e.g. `127.0.0.1`, `1.1.1.1`).
* `zone_id` &mdash; (Required, string) ID of zone this record is associated with.

## Attributes Reference

The following attributes are exported:

* `name` &mdash; (string) Name of record.
* `ttl` &mdash; (int) TTL of record.
* `type` &mdash; (string) Type of the record.
* `value` &mdash; (string) Value of record.
* `zone_id` &mdash; (string) ID of zone this record is associated with.
