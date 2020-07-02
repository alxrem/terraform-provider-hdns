package hdns

import (
	"context"
	"fmt"
	"github.com/alxrem/hdns-go/hdns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"log"
	"os"
	"strings"
	"testing"
)

const testAccPrefix = "hdns-testacc"

func testAccZoneName(rName string) string {
	return fmt.Sprintf("%s.dev", rName)
}

func init() {
	resource.AddTestSweepers("hdns_zone", &resource.Sweeper{
		Name: "hdns_zone",
		F:    testSweepZones,
	})
}

func TestAccHdnsZone_Basic(t *testing.T) {
	var zone hdns.Zone
	rName := acctest.RandomWithPrefix("hdns-testacc")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccHdnsPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccHdnsCheckZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHdnsCheckZoneConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccHdnsCheckZoneExists("hdns_zone.foobar", &zone),
					resource.TestCheckResourceAttr(
						"hdns_zone.foobar", "name", testAccZoneName(rName)),
					resource.TestCheckResourceAttr(
						"hdns_zone.foobar", "ttl", "900"),
				),
			},
			{
				Config: testAccHdnsCheckZoneConfig_changedTtl(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccHdnsCheckZoneExists("hdns_zone.foobar", &zone),
					resource.TestCheckResourceAttr(
						"hdns_zone.foobar", "name", testAccZoneName(rName)),
					resource.TestCheckResourceAttr(
						"hdns_zone.foobar", "ttl", "1800"),
				),
			},
			{
				ResourceName:      "hdns_zone.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

//noinspection GoSnakeCaseUsage
func testAccHdnsCheckZoneConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "hdns_zone" "foobar" {
  name = "%s"
  ttl  = 900
}
`, testAccZoneName(rName))
}

//noinspection GoSnakeCaseUsage
func testAccHdnsCheckZoneConfig_changedTtl(rName string) string {
	return fmt.Sprintf(`
resource "hdns_zone" "foobar" {
  name = "%s"
  ttl  = 1800
}
`, testAccZoneName(rName))
}

func testAccHdnsCheckZoneDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*providerConfig)
	client := config.client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hdns_zone" {
			continue
		}

		zone, r, err := client.Zone.GetByID(context.Background(), rs.Primary.ID)
		if err != nil && r != nil && r.StatusCode != 404 {
			return fmt.Errorf(
				"error checking if Zone (%s) is deleted: %v",
				rs.Primary.ID, err)
		}
		if zone != nil {
			return fmt.Errorf("zone (%s) has not been deleted", rs.Primary.ID)
		}
	}

	return nil
}

func testAccHdnsCheckZoneExists(z string, zone *hdns.Zone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[z]

		if !ok {
			return fmt.Errorf("not found: %s", z)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		config := testAccProvider.Meta().(*providerConfig)
		client := config.client

		// Try to find the key
		foundZone, _, err := client.Zone.GetByID(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundZone == nil {
			return fmt.Errorf("record not found")
		}

		*zone = *foundZone
		return nil
	}
}

func testSweepZones(_ string) error {
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	zones, err := client.Zone.All(ctx)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Found %d zones to sweep", len(zones))

	for _, z := range zones {
		if strings.HasPrefix(z.Name, testAccPrefix) {
			log.Printf("Deleting zone %s", z.Name)
			if _, err := client.Zone.Delete(ctx, z.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func createClient() (*hdns.Client, error) {
	if os.Getenv("HDNS_TOKEN") == "" {
		return nil, fmt.Errorf("empty HDNS_TOKEN")
	}
	opts := []hdns.ClientOption{
		hdns.WithToken(os.Getenv("HDNS_TOKEN")),
	}
	return hdns.NewClient(opts...), nil
}
