package zuora

import (
	"context"
	"net/http"

	"github.com/acsbe/terraform-provider-zuora/zuora/client"
	"github.com/acsbe/terraform-provider-zuora/zuora/resources/notifications"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns the Terraform provider for Zuora.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ZUORA_CLIENT_ID", nil),
				Description: "OAuth client ID for Zuora",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ZUORA_CLIENT_SECRET", nil),
				Description: "OAuth client secret for Zuora",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ZUORA_ENDPOINT", nil),
				Description: "Base URL for the Zuora REST API",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"zuora_notifications_callout_template": notifications.ResourceNotificationsCalloutTemplate(),
			"zuora_notification_callout_binding":   notifications.ResourceNotificationsCalloutBinding(),
		},

		ConfigureContextFunc: func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			cfg := &client.Config{
				ClientID:     d.Get("client_id").(string),
				ClientSecret: d.Get("client_secret").(string),
				Endpoint:     d.Get("endpoint").(string),
				HTTPClient:   http.DefaultClient,
			}
			return cfg, nil
		},
	}
}
