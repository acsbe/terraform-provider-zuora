package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/acsbe/terraform-provider-zuora/zuora/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"io"
)

// ResourceNotificationsCalloutTemplate manages Zuora callout templates using raw JSON.
// The body field must contain valid JSON and at minimum include "name", "calloutBaseurl", and "httpMethod".
func ResourceNotificationsCalloutTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCalloutTemplateCreate,
		ReadContext:   resourceCalloutTemplateRead,
		UpdateContext: resourceCalloutTemplateUpdate,
		DeleteContext: resourceCalloutTemplateDelete,

		Schema: map[string]*schema.Schema{
			"body": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				Description:  "Raw JSON payload for the callout template. Must include 'name', 'calloutBaseurl', and 'httpMethod'.",
			},
			"callout_template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the callout template.",
			},
		},
	}
}

func resourceCalloutTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*client.Config)
	body := d.Get("body").(string)

	// validate JSON and required fields
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return diag.Errorf("invalid JSON in body: %s", err)
	}
	for _, f := range []string{"name", "calloutBaseurl", "httpMethod"} {
		if _, ok := payload[f]; !ok {
			return diag.Errorf("missing required field %q in body", f)
		}
	}

	// POST to create
	req, err := cfg.NewRequest(ctx, "POST", "/notifications/callout-templates", strings.NewReader(body))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return diag.Errorf("error creating callout template: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	// parse response
	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.Errorf("error decoding create response: %s", err)
	}
	d.SetId(result.ID)
	d.Set("callout_template_id", result.ID)

	return resourceCalloutTemplateRead(ctx, d, m)
}

func resourceCalloutTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*client.Config)

	// GET existing
	req, err := cfg.NewRequest(ctx, "GET", fmt.Sprintf("/notifications/callout-templates/%s", d.Id()), nil)
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
		return diag.Errorf("error reading callout template: status %d", resp.StatusCode)
	}

	// refresh computed ID
	d.Set("callout_template_id", d.Id())
	return nil
}

func resourceCalloutTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*client.Config)
	body := d.Get("body").(string)

	// validate JSON
	if !json.Valid([]byte(body)) {
		return diag.Errorf("invalid JSON in body")
	}

	// PUT to update
	req, err := cfg.NewRequest(ctx, "PUT", fmt.Sprintf("/notifications/callout-templates/%s", d.Id()), strings.NewReader(body))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 {
		return diag.Errorf("error updating callout template: status %d", resp.StatusCode)
	}

	return resourceCalloutTemplateRead(ctx, d, m)
}

func resourceCalloutTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*client.Config)

	// DELETE resource
	req, err := cfg.NewRequest(ctx, "DELETE", fmt.Sprintf("/notifications/callout-templates/%s", d.Id()), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 && resp.StatusCode != http.StatusNotFound {
		return diag.Errorf("error deleting callout template: status %d", resp.StatusCode)
	}

	// clear state
	d.SetId("")
	return nil
}
