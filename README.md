# Terraform Provider: Zuora

A custom Terraform provider to manage Zuora callout templates and notification definitions via Zuora’s REST API.

## Prerequisites

* **Go 1.24+** (to build the provider binary)
* **Terraform 1.0+** (to use the provider)
* **Zuora API credentials** (Client ID, Client Secret, and API Endpoint)

## Quick Start

### 1. Build locally

```bash
# From provider root
go build -o terraform-provider-zuora_v1.0.0
```

### 2. Install for Terraform CLI

```bash
# Adjust OS/ARCH if needed (e.g. darwin_arm64)
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/local/zuora/1.0.0/darwin_arm64
mv terraform-provider-zuora_v1.0.0 \
  ~/.terraform.d/plugins/registry.terraform.io/local/zuora/1.0.0/darwin_arm64/
```

### 3. Use in your Terraform config

```hcl
terraform {
  required_providers {
    zuora = { source = "local/zuora", version = "1.0.0" }
  }
}
provider "zuora" {
  client_id     = var.client_id
  client_secret = var.client_secret
  endpoint      = var.endpoint
}
```

### 4. Resources

* **Callout Template**: `zuora_notifications_callout_template`

    * Create a callout template in Zuora.
    * Docs: [https://developer.zuora.com/v1-api-reference/api/operation/CreateCalloutTemplate/](https://developer.zuora.com/v1-api-reference/api/operation/CreateCalloutTemplate/)

* **Notification Definition**: `zuora_notifications_notification_definition`

    * Create or update a notification in Zuora.
    * Docs: [https://developer.zuora.com/v1-api-reference/api/operation/POST\_Create\_Notification\_Definition/](https://developer.zuora.com/v1-api-reference/api/operation/POST_Create_Notification_Definition/)

Pass your full JSON payload in the `body` attribute for each resource.

## Release Workflow

1. **Tag a new version** in Git:

   ```bash
   git tag v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
2. **GitHub Actions** automatically builds and zips binaries for all platforms, and publishes a GitHub Release under that tag.
3. **Terraform Cloud Private Registry** is connected to the GitHub repo; it will auto-discover new releases on each tag—no manual refresh needed.
4. **Consume** in Terraform by setting:

   ```hcl
   terraform {
     required_providers {
       zuora = { source = "acsbe/zuora", version = "1.0.0" }
     }
   }
   ```

---
