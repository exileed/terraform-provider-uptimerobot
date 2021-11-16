package provider

import (
	"context"
	"strconv"

	"github.com/exileed/uptimerobotapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccount() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information about the current UptimeRobot account.",

		ReadContext: dataSourceAccountRead,

		Schema: map[string]*schema.Schema{
			"email":            {Computed: true, Type: schema.TypeString},
			"monitor_limit":    {Computed: true, Type: schema.TypeInt},
			"monitor_interval": {Computed: true, Type: schema.TypeInt},
			"up_monitors":      {Computed: true, Type: schema.TypeInt},
			"down_monitors":    {Computed: true, Type: schema.TypeInt},
			"paused_monitors":  {Computed: true, Type: schema.TypeInt},
		},
	}
}

func dataSourceAccountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(uptimerobotapi.Client)

	var resp *uptimerobotapi.AccountResp
	var err error

	err = retryTime(func() error {
		resp, err = client.Account.GetAccountDetails()
		return err
	}, timeoutMinutes)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	account := resp.Account

	d.SetId(strconv.Itoa(account.UserId))
	d.Set("email", account.Email)
	d.Set("monitor_limit", account.MonitorLimit)
	d.Set("monitor_interval", account.MonitorInterval)
	d.Set("up_monitors", account.UpMonitors)
	d.Set("down_monitors", account.DownMonitors)
	d.Set("paused_monitors", account.PausedMonitors)

	return nil
}
