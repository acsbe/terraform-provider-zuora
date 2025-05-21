package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/acsbe/terraform-provider-zuora/zuora/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceNotificationsNotificationDefinition manages Zuora notification definitions using raw JSON.
// The body field must contain valid JSON and at minimum include the 'name' attribute.
func ResourceNotificationsNotificationDefinition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotificationCreate,
		ReadContext:   resourceNotificationRead,
		UpdateContext: resourceNotificationUpdate,
		DeleteContext: resourceNotificationDelete,

		Schema: map[string]*schema.Schema{
			"body": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Raw JSON payload for the notification definition. Must include the 'name' field.",
			},
			"notification_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the notification definition.",
			},
		},
	}
}

func resourceNotificationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*client.Config)
	body := d.Get("body").(string)

	// ensure valid JSON at create
	var tmp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &tmp); err != nil {
		return diag.Errorf("invalid JSON in body: %s", err)
	}

	req, err := cfg.NewRequest(ctx, "POST", "/notifications/notification-definitions", strings.NewReader(body))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return diag.Errorf("error creating notification definition: status %d, response: %s", resp.StatusCode, string(respBody))
	}

	var res struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return diag.Errorf("error decoding create response: %s", err)
	}
	d.SetId(res.ID)
	d.Set("notification_id", res.ID)

	return resourceNotificationRead(ctx, d, m)
}

func resourceNotificationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*client.Config)
	id := d.Id()

	req, err := cfg.NewRequest(ctx, "GET", fmt.Sprintf("/notifications/notification-definitions/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return diag.Errorf("error reading notification definition: status %d, response: %s", resp.StatusCode, string(respBody))
	}

	// only refresh ID; body remains as user provided
	d.Set("notification_id", id)
	return nil
}

func resourceNotificationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*client.Config)
	id := d.Id()
	body := d.Get("body").(string)

	req, err := cfg.NewRequest(ctx, "PUT", fmt.Sprintf("/notifications/notification-definitions/%s", id), strings.NewReader(body))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return diag.Errorf("error updating notification definition: status %d, response: %s", resp.StatusCode, string(respBody))
	}

	return resourceNotificationRead(ctx, d, m)
}

func resourceNotificationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*client.Config)
	id := d.Id()

	req, err := cfg.NewRequest(ctx, "DELETE", fmt.Sprintf("/notifications/notification-definitions/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 && resp.StatusCode != http.StatusNotFound {
		respBody, _ := io.ReadAll(resp.Body)
		return diag.Errorf("error deleting notification definition: status %d, response: %s", resp.StatusCode, string(respBody))
	}

	d.SetId("")
	return nil
}
