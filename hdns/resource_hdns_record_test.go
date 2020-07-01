package hdns

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"gitlab.com/alxrem/hdns-go/hdns"
	"testing"
)

func TestAccHdnsRecord_Basic(t *testing.T) {
	var record hdns.Record
	rName := acctest.RandomWithPrefix("hdns-testacc")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccHdnsPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccHdnsCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHdnsCheckRecordConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccHdnsCheckRecordExists("hdns_record.basic_record", &record),
					resource.TestCheckResourceAttr(
						"hdns_record.basic_record", "name", "www"),
					resource.TestCheckResourceAttr(
						"hdns_record.basic_record", "ttl", "86400"),
				),
			},
			{
				Config: testAccHdnsCheckRecordConfig_multiple_records(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccHdnsCheckRecordExists("hdns_record.multi-record[0]", &record),
					//resource.TestCheckResourceAttr(
					//	"hdns_record.bar[0]", "name", "www-0"),
					//resource.TestCheckResourceAttr(
					//	"hdns_record.bar[0]", "ttl", "86400"),
				),
			},
		},
	})
}

//noinspection GoSnakeCaseUsage
func testAccHdnsCheckRecordConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "hdns_zone" "foo" {
  name = "%s.dev"
}

resource "hdns_record" "basic_record" {
  name    = "www"
  type    = "A"
  value   = "1.1.1.1"
  zone_id = hdns_zone.foo.id 
}
`, rName)
}

//noinspection GoSnakeCaseUsage
func testAccHdnsCheckRecordConfig_multiple_records(rName string) string {
	return fmt.Sprintf(`
resource "hdns_zone" "foo" {
  name = "%s.dev"
}

resource "hdns_record" "multiple_record" {
  count = 10

  name    = "www-${count.index}"
  type    = "A"
  value   = "1.1.1.${count.index}"
  zone_id = hdns_zone.foo.id
}
`, rName)
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
