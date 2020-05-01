package hdns

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"gitlab.com/alxrem/hdns-go/hdns"
	"log"
)

func resourceRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceRecordCreate,
		Read:   resourceRecordRead,
		Update: resourceRecordUpdate,
		Delete: resourceRecordDelete,

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
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

//noinspection GoUnhandledErrorResult
func resourceRecordCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*providerConfig)
	ctx := context.Background()
	opts := hdns.RecordCreateOpts{
		Name:   d.Get("name").(string),
		TTL:    d.Get("ttl").(int),
		Type:   d.Get("type").(string),
		Value:  d.Get("value").(string),
		ZoneID: d.Get("zone_id").(string),
	}

	record, _, err := config.client.Record.Create(ctx, opts)
	if err != nil {
		return err
	}

	d.SetId(record.ID)
	d.Set("name", record.Name)
	d.Set("ttl", record.TTL)
	d.Set("type", record.Type)
	d.Set("value", record.Value)
	d.Set("zone_id", record.ZoneID)

	return nil
}

//noinspection GoUnhandledErrorResult
func resourceRecordRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*providerConfig)
	ctx := context.Background()

	id := d.Id()

	record, _, err := config.client.Record.GetByID(ctx, id)
	if err != nil {
		return err
	}

	d.SetId(record.ID)
	d.Set("name", record.Name)
	d.Set("ttl", record.TTL)
	d.Set("type", record.Type)
	d.Set("value", record.Value)
	d.Set("zone_id", record.ZoneID)

	return nil
}

func resourceRecordUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*providerConfig)
	ctx := context.Background()

	opts := hdns.RecordUpdateOpts{
		Name:   d.Get("name").(string),
		TTL:    d.Get("ttl").(int),
		Type:   d.Get("ttl").(string),
		Value:  d.Get("value").(string),
		ZoneID: d.Get("zone_id").(string),
	}

	_, _, err := config.client.Record.Update(ctx, d.Id(), opts)
	if err != nil {
		if resourceRecordIsNotFound(err, d) {
			return nil
		}
		return err
	}

	return resourceRecordRead(d, m)
}

func resourceRecordDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*providerConfig)
	ctx := context.Background()

	id := d.Id()

	r, err := config.client.Record.Delete(ctx, id)
	if err != nil && r != nil && r.StatusCode != 404 {
		return err
	}

	d.SetId("")

	return nil
}

func resourceRecordIsNotFound(err error, d *schema.ResourceData) bool {
	if hderr, ok := err.(hdns.Error); ok && hderr.Code == hdns.ErrorCodeNotFound {
		log.Printf("[WARN] Record (%s) not found, removing from state", d.Id())
		d.SetId("")
		return true
	}
	return false
}
