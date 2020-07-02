package hdns

import (
	"fmt"
	"github.com/alxrem/hdns-go/hdns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func init() {
	resource.AddTestSweepers("data_source_zone", &resource.Sweeper{
		Name: "hdns_zone_data_source",
		F:    testSweepZones,
	})
}

func TestAccHdnsDataSourceZone(t *testing.T) {
	var zone hdns.Zone
	rName := acctest.RandomWithPrefix("hdns-testacc")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccHdnsPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHdnsDataSourceZoneConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccHdnsCheckZoneExists("hdns_zone.foobar", &zone),
					resource.TestCheckResourceAttr("data.hdns_zone.foobar", "name", testAccZoneName(rName)),
					resource.TestCheckResourceAttr("data.hdns_zone.foobar", "ttl", "900"),
					resource.TestCheckResourceAttrSet("data.hdns_zone.foobar", "is_secondary_dns"),
					resource.TestCheckResourceAttrSet("data.hdns_zone.foobar", "paused"),
					resource.TestCheckResourceAttrSet("data.hdns_zone.foobar", "records_count"),
					resource.TestCheckResourceAttrSet("data.hdns_zone.foobar", "status"),
				),
			},
		},
	})
}

//noinspection GoSnakeCaseUsage
func testAccHdnsDataSourceZoneConfig(rName string) string {
	return fmt.Sprintf(`
resource "hdns_zone" "foobar" {
  name = "%s"
  ttl  = 900
}

data "hdns_zone" "foobar" {
  id = hdns_zone.foobar.id
}
`, testAccZoneName(rName))
}
