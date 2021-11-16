package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUptimeRobotDataSourceAccount(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testUptimeRobotDataSourceAccount,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.uptimerobot_account.test", "email"),
					resource.TestCheckResourceAttrSet("data.uptimerobot_account.test", "monitor_limit"),
					resource.TestCheckResourceAttrSet("data.uptimerobot_account.test", "monitor_interval"),
					resource.TestCheckResourceAttrSet("data.uptimerobot_account.test", "up_monitors"),
					resource.TestCheckResourceAttrSet("data.uptimerobot_account.test", "down_monitors"),
					resource.TestCheckResourceAttrSet("data.uptimerobot_account.test", "paused_monitors"),
				),
			},
		},
	})
}

const testUptimeRobotDataSourceAccount = `
data "uptimerobot_account" "test" {}
`
