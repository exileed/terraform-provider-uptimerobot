package provider

import (
	"context"
	"strconv"

	"github.com/exileed/uptimerobotapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var alertContactType = map[string]int{
	"email":      2,
	"twitter":    3,
	"boxcar":     4,
	"webhook":    5,
	"pushbullet": 6,
	"zapier":     7,
	"sms":        8,
	"pushover":   9,
	"hipchat":    10,
	"slack":      11,
	"phone":      13,
	"pagerduty":  16,
	"splunk":     15,
	"telegram":   18,
	"teams":      20,
	"hangouts":   21,
	"discord":    23,
}

func resourceAlertContact() *schema.Resource {
	return &schema.Resource{
		Description: "Uptimerobot alert contact resource",

		CreateContext: resourceAlertContactCreate,
		ReadContext:   resourceAlertContactRead,
		UpdateContext: resourceAlertContactUpdate,
		DeleteContext: resourceAlertContactDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"friendly_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(mapKeys(alertContactType), false),
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlertContactCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(uptimerobotapi.Client)

	acName := d.Get("friendly_name").(string)
	acValue := d.Get("value").(string)
	acType := d.Get("type").(string)

	acTypeStr := strconv.Itoa(alertContactType[acType])

	params := uptimerobotapi.NewAlertContactParams{TypeContact: acTypeStr, Value: acValue, FriendlyName: acName}

	ac, err := client.AlertContact.NewAlertContact(params)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	idStr := strconv.Itoa(ac.AlertContact.Id)

	d.SetId(idStr)

	getParams := uptimerobotapi.GetAlertContactsParams{
		AlertContacts: &idStr,
	}

	acs, err := client.AlertContact.GetAlertContacts(getParams)

	if acs.Total == 0 {
		return diag.Errorf("AlertContact %s not found", acName)
	}

	fillAlertContact(d, acs.AlertContacts[0])

	return nil
}

func resourceAlertContactRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(uptimerobotapi.Client)
	id := d.Id()

	getParams := uptimerobotapi.GetAlertContactsParams{
		AlertContacts: &id,
	}

	ac, err := client.AlertContact.GetAlertContacts(getParams)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	if ac.Total == 0 {
		return diag.Errorf("AlertContact %s not found", id)
	}

	alertContact := ac.AlertContacts[0]
	fillAlertContact(d, alertContact)

	return nil
}

func resourceAlertContactUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(uptimerobotapi.Client)

	id := d.Id()

	acName := d.Get("friendly_name").(string)
	acValue := d.Get("value").(string)

	idStr, err := strconv.Atoi(id)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	params := uptimerobotapi.EditAlertContactParams{Id: idStr, Value: &acValue, FriendlyName: &acName}

	_, err = client.AlertContact.EditAlertContact(params)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	return nil
}

func resourceAlertContactDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(uptimerobotapi.Client)

	id := d.Id()

	idStr, err := strconv.Atoi(id)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	_, err = client.AlertContact.DeleteAlertContact(idStr)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	return nil
}

func fillAlertContact(d *schema.ResourceData, ac uptimerobotapi.AlertContact) {
	d.Set("friendly_name", ac.FriendlyName)
	d.Set("value", ac.Value)
	d.Set("type", intToString(alertContactType, ac.Type))
	d.Set("status", ac.Status)
}
