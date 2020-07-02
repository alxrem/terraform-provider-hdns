package hdns

import (
	"context"
	"fmt"
	"github.com/alxrem/hdns-go/hdns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceZoneRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_secondary_dns": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"legacy_dns_host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"legacy_ns": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ns": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"paused": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"permission": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"records_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"registrar": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

//noinspection GoUnhandledErrorResult
func dataSourceZoneRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*providerConfig)
	ctx := context.Background()

	id := d.Get("id").(string)

	z, _, err := config.client.Zone.GetByID(ctx, id)
	if err != nil {
		if hderr, ok := err.(hdns.Error); ok && hderr.Code == hdns.ErrorCodeNotFound {
			return fmt.Errorf("no zone found with id %v", id)
		}
		return err
	}

	d.SetId(z.ID)
	d.Set("name", z.Name)
	d.Set("ttl", z.TTL)
	d.Set("is_secondary_dns", z.IsSecondaryDNS)
	d.Set("legacy_dns_host", z.LegacyDNSHost)
	d.Set("legacy_ns", z.LegacyNS)
	d.Set("ns", z.NS)
	d.Set("owner", z.Owner)
	d.Set("paused", z.Paused)
	d.Set("permission", z.Permission)
	d.Set("project", z.Project)
	d.Set("records_count", z.RecordsCount)
	d.Set("registrar", z.Registrar)
	d.Set("status", z.Status)

	return nil
}
