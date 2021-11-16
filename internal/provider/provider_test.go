package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	testAccProviders         map[string]*schema.Provider
	testAccProvider          *schema.Provider
	testAccProviderFactories map[string]func() (*schema.Provider, error)
)

func init() {
	testAccProvider = Provider("dev")
	testAccProviders = map[string]*schema.Provider{
		"uptimerobot": testAccProvider,
	}
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"uptimerobot": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider("dev").InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("UPTIMEROBOT_API_KEY"); v == "" {
		t.Fatal("UPTIMEROBOT_API_KEY must be set for acceptance tests")
	}
}
