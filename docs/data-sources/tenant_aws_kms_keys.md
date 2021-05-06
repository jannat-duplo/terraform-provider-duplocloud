---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "duplocloud_tenant_aws_kms_keys Data Source - terraform-provider-duplocloud"
subcategory: ""
description: |-
  
---

# duplocloud_tenant_aws_kms_keys (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **tenant_id** (String)

### Optional

- **id** (String) The ID of this resource.
- **selectable** (Boolean)

### Read-Only

- **keys** (List of Object) (see [below for nested schema](#nestedatt--keys))

<a id="nestedatt--keys"></a>
### Nested Schema for `keys`

Read-Only:

- **key_arn** (String)
- **key_id** (String)
- **key_name** (String)

