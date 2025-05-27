# Terraform Provider: Zuora

A Terraform provider for managing [Zuora](https://www.zuora.com/) callout templates and their bindings to notification definitions.

👉 [View on Terraform Registry](https://registry.terraform.io/providers/acsbe/zuora/latest)

---

## Requirements

- [Terraform 1.0+](https://www.terraform.io/downloads)
- [Go 1.24+](https://go.dev/dl/) (for development only)

---

## Usage

```hcl
terraform {
  required_providers {
    zuora = {
      source  = "acsbe/zuora"
      version = "~> 1.0"
    }
  }
}

provider "zuora" {
  client_id     = var.client_id
  client_secret = var.client_secret
  endpoint      = var.endpoint
}
```

---

## Resources

- [`zuora_notifications_callout_template`]()  
  Create or manage callout templates using raw JSON payloads.

- [`zuora_notifications_callout_binding`]()  
  Attach or detach a callout template to/from a notification definition.

---

## Development

To build the provider locally:

```bash
go build -o terraform-provider-zuora
```
