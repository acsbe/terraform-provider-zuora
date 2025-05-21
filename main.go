package main

import (
	"github.com/acsbe/terraform-provider-zuora/zuora"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: zuora.Provider,
	})
}
