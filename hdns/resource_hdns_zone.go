package hdns

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"gitlab.com/alxrem/hdns-go/hdns"
	"log"
)

func resourceZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceZoneCreate,
		Read:   resourceZoneRead,
		Update: resourceZoneUpdate,
		Delete: resourceZoneDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Default:  86400,
				Optional: true,
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
func resourceZoneCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*providerConfig)
	ctx := context.Background()
	opts := hdns.ZoneCreateOpts{
		Name: d.Get("name").(string),
		TTL:  d.Get("ttl").(int),
	}

	z, _, err := config.client.Zone.Create(ctx, opts)
	if err != nil {
		return err
	}

	d.SetId(z.ID)

	return resourceZoneRead(d, m)
}

//noinspection GoUnhandledErrorResult
func resourceZoneRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*providerConfig)
	ctx := context.Background()

	id := d.Id()

	z, _, err := config.client.Zone.GetByID(ctx, id)
	if err != nil {
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

func resourceZoneUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*providerConfig)
	ctx := context.Background()

	opts := hdns.ZoneUpdateOpts{
		Name: d.Get("name").(string),
		TTL:  d.Get("ttl").(int),
	}

	_, _, err := config.client.Zone.Update(ctx, d.Id(), opts)
	if err != nil {
		if resourceZoneIsNotFound(err, d) {
			return nil
		}
		return err
	}

	return resourceZoneRead(d, m)
}

func resourceZoneDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*providerConfig)
	ctx := context.Background()

	id := d.Id()

	r, err := config.client.Zone.Delete(ctx, id)
	if err != nil && r != nil && r.StatusCode != 404 {
		return err
	}

	d.SetId("")

	return nil
}

func resourceZoneIsNotFound(err error, d *schema.ResourceData) bool {
	if hderr, ok := err.(hdns.Error); ok && hderr.Code == hdns.ErrorCodeNotFound {
		log.Printf("[WARN] Zone (%s) not found, removing from state", d.Id())
		d.SetId("")
		return true
	}
	return false
}
