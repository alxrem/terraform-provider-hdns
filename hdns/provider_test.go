package hdns

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"hdns": testAccProvider,
	}
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccHdnsPreCheck(t *testing.T) {
	if v := os.Getenv("HDNS_TOKEN"); v == "" {
		t.Fatal("HDNS_TOKEN must be set for acceptance tests")
	}
}
