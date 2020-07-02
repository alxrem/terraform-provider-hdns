package hdns

import (
	"context"
	"fmt"
	"github.com/alxrem/hdns-go/hdns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

	zoneID := d.Get("zone_id").(string)
	opts := hdns.RecordCreateOpts{
		Name:   d.Get("name").(string),
		TTL:    d.Get("ttl").(int),
		Type:   d.Get("type").(string),
		Value:  d.Get("value").(string),
		ZoneID: zoneID,
	}

	mu := config.Mutex(zoneID)
	defer func() {
		mu.Unlock()
		log.Printf("[DEBUG] Released lock for zone %s", zoneID)
	}()

	mu.Lock()
	log.Printf("[DEBUG] Acquired lock for zone %s", zoneID)

	record, _, err := config.client.Record.Create(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to create record \"%s %s\" in the zone %s: %s", opts.Type, opts.Name, opts.ZoneID, err.Error())
	}

	if _, _, err := config.client.Record.GetByID(ctx, record.ID); err != nil {
		return fmt.Errorf("failed to create record \"%s %s\" in the zone %s: record is not exists just after creation", opts.Type, opts.Name, opts.ZoneID)
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
		if resourceRecordIsNotFound(err, d) {
			return nil
		}
		return fmt.Errorf("failed to read record %s: %s", id, err.Error())
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

	zoneID := d.Get("zone_id").(string)
	opts := hdns.RecordUpdateOpts{
		Name:   d.Get("name").(string),
		TTL:    d.Get("ttl").(int),
		Type:   d.Get("type").(string),
		Value:  d.Get("value").(string),
		ZoneID: zoneID,
	}

	mu := config.Mutex(zoneID)
	defer func() {
		mu.Unlock()
		log.Printf("[DEBUG] Released lock for zone %s", zoneID)
	}()

	mu.Lock()
	log.Printf("[DEBUG] Acquired lock for zone %s", zoneID)

	_, _, err := config.client.Record.Update(ctx, d.Id(), opts)
	if err != nil {
		if resourceRecordIsNotFound(err, d) {
			return nil
		}
		return fmt.Errorf("failed to create record \"%s %s\" in the zone %s: %s", opts.Type, opts.Name, opts.ZoneID, err.Error())
	}

	return resourceRecordRead(d, m)
}

func resourceRecordDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*providerConfig)
	ctx := context.Background()

	id := d.Id()

	zoneID := d.Get("zone_id").(string)
	mu := config.Mutex(zoneID)
	defer func() {
		mu.Unlock()
		log.Printf("[DEBUG] Released lock for zone %s", zoneID)
	}()

	mu.Lock()
	log.Printf("[DEBUG] Acquired lock for zone %s", zoneID)

	r, err := config.client.Record.Delete(ctx, id)
	if err != nil && r != nil && r.StatusCode != 404 {
		return fmt.Errorf("failed to delete record %s: %s", id, err.Error())
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
