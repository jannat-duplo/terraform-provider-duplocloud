---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "duplocloud_admin_system_setting Resource - terraform-provider-duplocloud"
subcategory: ""
description: |-
  duplocloud_admin_system_setting manages an admin system setting in Duplo.
---

# duplocloud_admin_system_setting (Resource)

`duplocloud_admin_system_setting` manages an admin system setting in Duplo.

## Example Usage

```terraform
resource "duplocloud_admin_system_setting" "test-setting" {
  key   = "EnableVPN"
  value = "true"
  type  = "Flags"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **key** (String) Key name for the system setting.
- **type** (String) Type of the system setting.
- **value** (String) Value for the system setting.

### Optional

- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- **id** (String) The ID of this resource.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)

## Import

Import is supported using the following syntax:

```shell
# Example: Importing an existing duplo admin setting
#  - *KEY_TYPE* is the type of setting key
#  - *KEY* is the key name
#
terraform import duplocloud_admin_system_setting.mySetting *KEY_TYPE*/*KEY*
```
