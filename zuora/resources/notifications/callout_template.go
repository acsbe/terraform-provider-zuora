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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceNotificationsCalloutTemplate manages Zuora callout templates using raw JSON.
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
		},
	}
}

func resourceCalloutTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*client.Config)
	body := d.Get("body").(string)

	// validate JSON
	if !json.Valid([]byte(body)) {
		return diag.Errorf("invalid JSON in body")
	}

	req, err := cfg.NewRequest(ctx, "POST", "/notifications/callout-templates", strings.NewReader(body))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return diag.Errorf("error creating callout template: status %d, response: %s", resp.StatusCode, string(b))
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.Errorf("error decoding create response: %s", err)
	}
	d.SetId(result.ID)
	return resourceCalloutTemplateRead(ctx, d, m)
}

func resourceCalloutTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*client.Config)

	req, err := cfg.NewRequest(ctx, "GET", fmt.Sprintf("/notifications/callout-templates/%s", d.Id()), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode >= 300 {
		return diag.Errorf("error reading callout template: status %d", resp.StatusCode)
	}

	return nil
}

func resourceCalloutTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*client.Config)
	body := d.Get("body").(string)

	if !json.Valid([]byte(body)) {
		return diag.Errorf("invalid JSON in body")
	}

	req, err := cfg.NewRequest(ctx, "PUT", fmt.Sprintf("/notifications/callout-templates/%s", d.Id()), strings.NewReader(body))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return diag.Errorf("error updating callout template: status %d, response: %s", resp.StatusCode, string(b))
	}

	return resourceCalloutTemplateRead(ctx, d, m)
}

func resourceCalloutTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*client.Config)

	req, err := cfg.NewRequest(ctx, "DELETE", fmt.Sprintf("/notifications/callout-templates/%s", d.Id()), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 && resp.StatusCode != http.StatusNotFound {
		b, _ := io.ReadAll(resp.Body)
		return diag.Errorf("error deleting callout template: status %d, response: %s", resp.StatusCode, string(b))
	}

	d.SetId("")
	return nil
}
