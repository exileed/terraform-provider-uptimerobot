package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUptimeRobotResourceAlertContact(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertContact,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("uptimerobot_alert_contact.test", "friendly_name", "me+test@exileed.com"),
					resource.TestCheckResourceAttr("uptimerobot_alert_contact.test", "type", "email"),
					resource.TestCheckResourceAttr("uptimerobot_alert_contact.test", "value", "me+test@exileed.com"),
				),
			},
		},
	})
}

const testAccResourceAlertContact = `
resource "uptimerobot_alert_contact" "test" {
  friendly_name = "me+test@exileed.com"
  type = "email"
  value = "me+test@exileed.com"
}
`
