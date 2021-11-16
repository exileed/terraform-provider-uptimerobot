package provider

import (
	"context"

	"github.com/exileed/uptimerobotapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func Provider(version string) *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UPTIMEROBOT_API_KEY", nil),
				Description: "API token for UptimeRobot API",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"uptimerobot_account": dataSourceAccount(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"uptimerobot_alert_contact": resourceAlertContact(),
			"uptimerobot_monitor":       resourceMonitor(),
		},
	}
	p.ConfigureContextFunc = configure(version, p)

	return p
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(cnt context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		userAgent := p.UserAgent("terraform-provider-uptimerobot", version)

		c := uptimerobotapi.ClientConfig{
			APIToken:  d.Get("api_key").(string),
			UserAgent: &userAgent,
		}
		api := uptimerobotapi.NewClientWithConfig(&c)

		return *api, nil
	}
}
