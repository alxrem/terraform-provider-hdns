package hdns

import (
	"context"
	"fmt"
	"github.com/alxrem/hdns-go/hdns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccHdnsRecord_Create(t *testing.T) {
	var record hdns.Record
	rName := acctest.RandomWithPrefix("hdns-testacc")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccHdnsPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccHdnsCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHdnsCheckRecordConfig_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccHdnsCheckRecordExists("hdns_record.bar", &record),
					resource.TestCheckResourceAttr("hdns_record.bar", "name", "www"),
					resource.TestCheckResourceAttr("hdns_record.bar", "ttl", "86400"),
					resource.TestCheckResourceAttr("hdns_record.bar", "type", "A"),
					resource.TestCheckResourceAttr("hdns_record.bar", "value", "1.1.1.1"),
					resource.TestCheckResourceAttrSet("hdns_record.bar", "zone_id"),
				),
			},
		},
	})
}

func TestAccHdnsRecord_Update(t *testing.T) {
	var record hdns.Record
	rName := acctest.RandomWithPrefix("hdns-testacc")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccHdnsPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccHdnsCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHdnsCheckRecordConfig_create(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccHdnsCheckRecordExists("hdns_record.bar", &record),
					resource.TestCheckResourceAttr("hdns_record.bar", "name", "www"),
					resource.TestCheckResourceAttr("hdns_record.bar", "ttl", "86400"),
					resource.TestCheckResourceAttr("hdns_record.bar", "type", "A"),
					resource.TestCheckResourceAttr("hdns_record.bar", "value", "1.1.1.1"),
					resource.TestCheckResourceAttrSet("hdns_record.bar", "zone_id"),
				),
			},
			{
				Config: testAccHdnsCheckRecordConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccHdnsCheckRecordExists("hdns_record.bar", &record),
					resource.TestCheckResourceAttr("hdns_record.bar", "name", "ddd"),
					resource.TestCheckResourceAttr("hdns_record.bar", "ttl", "1800"),
					resource.TestCheckResourceAttr("hdns_record.bar", "type", "A"),
					resource.TestCheckResourceAttr("hdns_record.bar", "value", "2.2.2.2"),
					resource.TestCheckResourceAttrSet("hdns_record.bar", "zone_id"),
				),
			},
		},
	})
}

func TestAccHdnsRecord_MultiRecords(t *testing.T) {
	var record hdns.Record
	rName := acctest.RandomWithPrefix("hdns-testacc")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccHdnsPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccHdnsCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHdnsCheckRecordConfig_multi_records(rName),
				Check: resource.ComposeTestCheckFunc(func() []resource.TestCheckFunc {
					checkFunks := make([]resource.TestCheckFunc, 0)
					for i := 0; i < 10; i++ {
						resourceName := fmt.Sprintf("hdns_record.bar.%d", i)
						checkFunks = append(checkFunks, []resource.TestCheckFunc{
							testAccHdnsCheckRecordExists(resourceName, &record),
							resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("www-%d", i)),
							resource.TestCheckResourceAttr(resourceName, "type", "A"),
							resource.TestCheckResourceAttr(resourceName, "value", fmt.Sprintf("1.1.1.%d", i)),
							resource.TestCheckResourceAttrSet(resourceName, "zone_id"),
						}...)
					}
					return checkFunks
				}()...,
				),
			},
		},
	})
}

func TestAccHdnsRecord_MultiZonesMultiRecords(t *testing.T) {
	//var record hdns.Record
	rName := acctest.RandomWithPrefix("hdns-testacc")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccHdnsPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccHdnsCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHdnsCheckRecordConfig_multi_zones_multi_records(rName),
				Check:  resource.ComposeTestCheckFunc(
				//testAccHdnsCheckRecordExists(`hdns_record.bar["0-0"]`, &record),
				),
			},
		},
	})
}

//noinspection GoSnakeCaseUsage
func testAccHdnsCheckRecordConfig_create(rName string) string {
	return fmt.Sprintf(`
resource "hdns_zone" "foo" {
  name = "%s.dev"
}

resource "hdns_record" "bar" {
  name    = "www"
  type    = "A"
  value   = "1.1.1.1"
  zone_id = hdns_zone.foo.id 
}
`, rName)
}

//noinspection GoSnakeCaseUsage
func testAccHdnsCheckRecordConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "hdns_zone" "foo" {
  name = "%s.dev"
}

resource "hdns_record" "bar" {
  name    = "ddd"
  ttl     = 1800
  type    = "A"
  value   = "2.2.2.2"
  zone_id = hdns_zone.foo.id 
}
`, rName)
}

//noinspection GoSnakeCaseUsage
func testAccHdnsCheckRecordConfig_multi_records(rName string) string {
	return fmt.Sprintf(`
resource "hdns_zone" "foo" {
  name = "%s.dev"
}

resource "hdns_record" "bar" {
  count = 10

  name    = "www-${count.index}"
  type    = "A"
  value   = "1.1.1.${count.index}"
  zone_id = hdns_zone.foo.id
}
`, rName)
}

//noinspection GoSnakeCaseUsage
func testAccHdnsCheckRecordConfig_multi_zones_multi_records(rName string) string {
	return fmt.Sprintf(`
resource "hdns_zone" "foo1" {
  name = "%s-1.dev"
}

resource "hdns_record" "bar1" {
  count = 4

  name    = "www-${count.index}"
  type    = "A"
  value   = "1.1.1.${count.index}"
  zone_id = hdns_zone.foo1.id
}

resource "hdns_zone" "foo2" {
  name = "%s-2.dev"
}

resource "hdns_record" "bar2" {
  count = 4

  name    = "www-${count.index}"
  type    = "A"
  value   = "2.1.1.${count.index}"
  zone_id = hdns_zone.foo2.id
}

resource "hdns_zone" "foo3" {
  name = "%s-3.dev"
}

resource "hdns_record" "bar3" {
  count = 4

  name    = "www-${count.index}"
  type    = "A"
  value   = "3.1.1.${count.index}"
  zone_id = hdns_zone.foo3.id
}
`, rName, rName, rName)
}

func testAccHdnsCheckRecordDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*providerConfig)
	client := config.client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hdns_record" {
			continue
		}

		record, r, err := client.Record.GetByID(context.Background(), rs.Primary.ID)
		if err != nil && r != nil && r.StatusCode != 404 {
			return fmt.Errorf(
				"error checking if Record (%s) is deleted: %v",
				rs.Primary.ID, err)
		}
		if record != nil {
			return fmt.Errorf("record (%s) has not been deleted", rs.Primary.ID)
		}
	}

	return nil
}

func testAccHdnsCheckRecordExists(r string, record *hdns.Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]

		if !ok {
			return fmt.Errorf("not found: %s", r)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		config := testAccProvider.Meta().(*providerConfig)
		client := config.client

		// Try to find the key
		foundRecord, _, err := client.Record.GetByID(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundRecord == nil {
			return fmt.Errorf("record not found")
		}

		*record = *foundRecord
		return nil
	}
}
