package notifications

import (
	"bytes"
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

// ResourceNotificationsCalloutBinding attaches or detaches a callout template to/from a notification definition.
func ResourceNotificationsCalloutBinding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAttachmentCreate,
		ReadContext:   resourceAttachmentRead,
		DeleteContext: resourceAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"callout_template_id": {Type: schema.TypeString, Required: true, ForceNew: true},
			"notification_id":     {Type: schema.TypeString, Required: true, ForceNew: true},
		},
	}
}

func resourceAttachmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*client.Config)
	ctID := d.Get("callout_template_id").(string)
	nID := d.Get("notification_id").(string)

	if err := attachCalloutToNotification(ctx, cfg, nID, ctID); err != nil {
		return diag.FromErr(err)
	}

	// Use composite ID to track this binding
	d.SetId(fmt.Sprintf("%s:%s", nID, ctID))
	return resourceAttachmentRead(ctx, d, m)
}

func resourceAttachmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// No-op: state is fully maintained in Terraform
	return nil
}

func resourceAttachmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*client.Config)
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format %s, expected notification_id:callout_template_id", d.Id())
	}
	nID, ctID := parts[0], parts[1]

	if err := detachCalloutFromNotification(ctx, cfg, nID, ctID); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// attachCalloutToNotification fetches the notification, appends the callout ID, and updates via PUT
func attachCalloutToNotification(ctx context.Context, cfg *client.Config, notificationID, calloutID string) error {
	// GET existing notification definition
	req, err := cfg.NewRequest(ctx, http.MethodGet,
		fmt.Sprintf("/notifications/notification-definitions/%s", notificationID), nil)
	if err != nil {
		return err
	}
	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error fetching notification %s: status %d, response: %s", notificationID, resp.StatusCode, string(b))
	}

	// Decode current payload
	var notif map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&notif); err != nil {
		return err
	}
	// Extract and append calloutTemplateIds
	idsIface, _ := notif["calloutTemplateIds"].([]interface{})
	ids := make([]string, 0, len(idsIface)+1)
	for _, v := range idsIface {
		if s, ok := v.(string); ok {
			ids = append(ids, s)
		}
	}
	ids = append(ids, calloutID)

	// Prepare update body
	payload, _ := json.Marshal(map[string]interface{}{"calloutTemplateIds": ids})

	// PUT update notification definition
	req, err = cfg.NewRequest(ctx, http.MethodPut,
		fmt.Sprintf("/notifications/notification-definitions/%s", notificationID),
		bytes.NewReader(payload))
	if err != nil {
		return err
	}
	resp, err = cfg.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error attaching callout to notification %s: status %d, response: %s", notificationID, resp.StatusCode, string(b))
	}

	return nil
}

// detachCalloutFromNotification fetches the notification, removes the callout ID, and updates via PUT
func detachCalloutFromNotification(ctx context.Context, cfg *client.Config, notificationID, calloutID string) error {
	// GET existing notification definition
	req, err := cfg.NewRequest(ctx, http.MethodGet,
		fmt.Sprintf("/notifications/notification-definitions/%s", notificationID), nil)
	if err != nil {
		return err
	}
	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error fetching notification %s: status %d, response: %s", notificationID, resp.StatusCode, string(b))
	}

	// Decode current payload
	var notif map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&notif); err != nil {
		return err
	}
	// Filter out the calloutID
	idsIface, _ := notif["calloutTemplateIds"].([]interface{})
	ids := make([]string, 0, len(idsIface))
	for _, v := range idsIface {
		if s, ok := v.(string); ok && s != calloutID {
			ids = append(ids, s)
		}
	}

	// Prepare update body
	payload, _ := json.Marshal(map[string]interface{}{"calloutTemplateIds": ids})

	// PUT update notification definition
	req, err = cfg.NewRequest(ctx, http.MethodPut,
		fmt.Sprintf("/notifications/notification-definitions/%s", notificationID),
		bytes.NewReader(payload))
	if err != nil {
		return err
	}
	resp, err = cfg.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error detaching callout from notification %s: status %d, response: %s", notificationID, resp.StatusCode, string(b))
	}

	return nil
}
