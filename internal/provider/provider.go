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

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"api_key": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("UPTIMEROBOT_API_KEY", nil),
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
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(cnt context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		// Setup a User-Agent for your API client (replace the provider name for yours):
		// userAgent := p.UserAgent("terraform-provider-scaffolding", version)
		// TODO: myClient.UserAgent = userAgent

		api := uptimerobotapi.NewClient(d.Get("api_key").(string))

		return *api, nil
	}
}
