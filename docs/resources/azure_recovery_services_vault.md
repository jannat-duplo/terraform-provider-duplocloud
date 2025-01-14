---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "duplocloud_azure_recovery_services_vault Resource - terraform-provider-duplocloud"
subcategory: ""
description: |-
  duplocloud_azure_recovery_services_vault manages an Azure Recovery Services Vault in Duplo.
---

# duplocloud_azure_recovery_services_vault (Resource)

`duplocloud_azure_recovery_services_vault` manages an Azure Recovery Services Vault in Duplo.

## Example Usage

```terraform
resource "duplocloud_azure_recovery_services_vault" "recovery_services_vault" {
  infra_name          = "demo"
  resource_group_name = "duploinfra-demo"
  name                = "test"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **infra_name** (String) The name of the infrastructure. Infrastructure names are globally unique and less than 13 characters.
- **name** (String) Specifies the name of the Recovery Services Vault. Recovery Service Vault name must be 2 - 50 characters long, start with a letter, contain only letters, numbers and hyphens. Changing this forces a new resource to be created.

### Optional

- **resource_group_name** (String) The name of the resource group in which to create the Recovery Services Vault. Changing this forces a new resource to be created.
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- **wait_until_ready** (Boolean) Whether or not to wait until Recovery Services Vault to be ready, after creation. Defaults to `true`.

### Read-Only

- **azure_id** (String)
- **id** (String) The ID of this resource.
- **location** (String)
- **sku** (String)

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)

## Import

Import is supported using the following syntax:

```shell
# Example: Importing an existing Azure Recovery Services Vault
#  - *INFRA_NAME* is the name of duplo infra.
#  - *NAME* is the name of the Recovery Services Vault
#
terraform import duplocloud_azure_recovery_services_vault.recovery_services_vault *INFRA_NAME*/*NAME*
```
