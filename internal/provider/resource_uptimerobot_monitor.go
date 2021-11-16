package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/exileed/uptimerobotapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var monitorType = map[string]int{
	"http":    1,
	"keyword": 2,
	"ping":    3,
	"port":    4,
}
var monitorSubType = map[string]int{
	"http":   1,
	"https":  2,
	"ftp":    3,
	"smtp":   4,
	"pop3":   5,
	"imap":   6,
	"custom": 99,
}

var monitorHTTPAuthType = map[string]int{
	"basic":  1,
	"digest": 2,
}

var monitorHTTPMethodType = map[string]int{
	"head":    1,
	"get":     2,
	"post":    3,
	"put":     4,
	"patch":   5,
	"delete":  6,
	"options": 7,
}

var monitorStatusType = map[string]int{
	"paused":          0,
	"not_checked_yet": 1,
	"up":              2,
	"seems_down":      8,
	"down":            9,
}

func resourceMonitor() *schema.Resource {
	return &schema.Resource{
		Description: "Uptimerobot monitor resource",

		CreateContext: resourceMonitorCreate,
		ReadContext:   resourceMonitorRead,
		UpdateContext: resourceMonitorUpdate,
		DeleteContext: resourceMonitorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"friendly_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The friendly name of the monitor.",
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(mapKeys(monitorType), false),
			},
			"sub_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(mapKeys(monitorSubType), false),
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"http_username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"http_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"http_auth_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(mapKeys(monitorHTTPAuthType), false),
			},
			//"http_method": {
			//	Type:         schema.TypeString,
			//	Optional:     true,
			//	ValidateFunc: validation.StringInSlice(mapKeys(monitorHTTPMethodType), false),
			//},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ignore_ssl_errors": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"alert_contact": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"recurrence": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
					},
				},
			},
		},
	}
}

func resourceMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(uptimerobotapi.Client)

	mType := d.Get("type").(string)
	mInterval := d.Get("interval").(int)
	mTimeout := d.Get("timeout").(int)

	request := uptimerobotapi.NewMonitorsParams{
		FriendlyName: d.Get("friendly_name").(string),
		Url:          d.Get("url").(string),
		Type:         monitorType[mType],
		////KeywordType: &monitorSubType[mSubType], //@todo
		////KeywordCaseType: &monitorSubType[mSubType],
		////KeywordValue: &monitorSubType[mSubType],
		Interval: &mInterval,
		Timeout:  &mTimeout,
	}

	mPort := d.Get("port")
	if mPort != nil {
		mPortInt := mPort.(int)
		request.Port = &mPortInt
	}

	mSubType := d.Get("sub_type")
	if mSubType != nil {
		mSubTypeInt := monitorSubType[mSubType.(string)]
		request.SubType = &mSubTypeInt
	}

	mHttpUsername := d.Get("http_username")
	if mHttpUsername != nil {
		mHttpUsernameString := mHttpUsername.(string)
		request.HttpUsername = &mHttpUsernameString
	}

	mHttpPassword := d.Get("http_password")
	if mHttpPassword != nil {
		mHttpPasswordString := mHttpPassword.(string)
		request.HttpPassword = &mHttpPasswordString
	}

	mHttpAuthType := d.Get("http_auth_type")
	if mHttpAuthType != nil {
		mHttpAuthTypeInt := monitorHTTPAuthType[mHttpAuthType.(string)]
		request.HttpAuthType = &mHttpAuthTypeInt
	}

	alertContactMap := d.Get("alert_contact").([]interface{})
	acStrings := make([]string, len(alertContactMap))

	for k, v := range alertContactMap {
		id := v.(map[string]interface{})["id"].(string)
		threshold := v.(map[string]interface{})["threshold"].(int)
		recurrence := v.(map[string]interface{})["recurrence"].(int)

		acStrings[k] = fmt.Sprintf("%s_%d_%d", id, threshold, recurrence)
	}
	alertContactStr := strings.Join(acStrings, "-")
	request.AlertContacts = &alertContactStr

	var monitor *uptimerobotapi.MonitorsSingResp
	var err error

	err = retryTime(func() error {
		monitor, err = client.Monitor.NewMonitor(request)
		return err
	}, timeoutMinutes)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	d.SetId(strconv.Itoa(monitor.Monitor.Id))

	return nil
}

func resourceMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(uptimerobotapi.Client)
	id := d.Id()

	request := uptimerobotapi.GetMonitorsParams{
		Monitors:      &id,
		AlertContacts: 1,
		SSL:           1,
	}

	var m *uptimerobotapi.MonitorsResp
	var err error

	err = retryTime(func() error {
		m, err = client.Monitor.GetMonitors(request)
		return err
	}, timeoutMinutes)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	if m.Pagination.Total == 0 {
		return diag.Errorf("Monitor %d not found", m.Monitors[0].Id)
	}

	monitor := m.Monitors[0]

	d.Set("id", monitor.Id)
	fillMonitor(d, monitor)

	return nil
}

func resourceMonitorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(uptimerobotapi.Client)

	id := d.Id()

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	mInterval := d.Get("interval").(int)
	mTimeout := d.Get("timeout").(int)

	request := uptimerobotapi.EditMonitorsParams{
		FriendlyName: d.Get("friendly_name").(string),
		Url:          d.Get("url").(string),
		Interval:     &mInterval,
		Timeout:      &mTimeout,
	}

	mPort := d.Get("port")
	if mPort != nil {
		mPortInt := mPort.(int)
		request.Port = &mPortInt
	}

	mSubType := d.Get("sub_type")
	if mSubType != nil {
		mSubTypeInt := monitorSubType[mSubType.(string)]
		request.SubType = &mSubTypeInt
	}

	mHttpUsername := d.Get("http_username")
	if mHttpUsername != nil {
		mHttpUsernameString := mHttpUsername.(string)
		request.HttpUsername = &mHttpUsernameString
	}

	mHttpPassword := d.Get("http_password")
	if mHttpPassword != nil {
		mHttpPasswordString := mHttpPassword.(string)
		request.HttpPassword = &mHttpPasswordString
	}

	mHttpAuthType := d.Get("http_auth_type")
	if mHttpAuthType != nil {
		mHttpAuthTypeInt := monitorHTTPAuthType[mHttpAuthType.(string)]
		request.HttpAuthType = &mHttpAuthTypeInt
	}

	alertContactMap := d.Get("alert_contact").([]interface{})
	acStrings := make([]string, len(alertContactMap))

	for k, v := range alertContactMap {

		id := v.(map[string]interface{})["id"].(string)
		threshold := v.(map[string]interface{})["threshold"].(int)
		recurrence := v.(map[string]interface{})["recurrence"].(int)

		acStrings[k] = fmt.Sprintf("%s_%d_%d", id, threshold, recurrence)
	}
	alertContactStr := strings.Join(acStrings, "-")
	request.AlertContacts = &alertContactStr

	err = retryTime(func() error {
		_, err = client.Monitor.EditMonitor(idInt, request)
		return err
	}, timeoutMinutes)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	return resourceMonitorRead(ctx, d, meta)
}

func resourceMonitorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(uptimerobotapi.Client)

	id := d.Id()

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	err = retryTime(func() error {
		_, err = client.Monitor.DeleteMonitor(idInt)
		return err
	}, timeoutMinutes)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	return nil
}

func fillMonitor(d *schema.ResourceData, m uptimerobotapi.Monitor) {
	d.Set("friendly_name", m.FriendlyName)
	d.Set("url", m.Url)
	d.Set("type", intToString(monitorType, m.Type))
	d.Set("sub_type", m.SubType)
	d.Set("keyword_type", m.KeywordType)
	d.Set("keyword_case_type", m.KeywordCaseType)
	d.Set("keyword_value", m.KeywordValue)
	d.Set("http_username", m.HttpUsername)
	d.Set("http_password", m.HttpPassword)
	d.Set("port", m.Port)
	d.Set("interval", m.Interval)
	d.Set("timeout", m.Timeout)
	d.Set("ignore_ssl_errors", m.SSL.IgnoreErrors == 1)
	d.Set("status", intToString(monitorStatusType, m.Status))

	rawContacts := make([]map[string]interface{}, len(*m.AlertContacts))

	for k, v := range *m.AlertContacts {
		rawContacts[k] = map[string]interface{}{
			"id":         v.Id,
			"recurrence": v.Recurrence,
			"threshold":  v.Threshold,
		}
	}
	if err := d.Set("alert_contact", rawContacts); err != nil {
		fmt.Errorf("error setting alert_contact for resource %s: %s", d.Id(), err.Error())
	}
}
