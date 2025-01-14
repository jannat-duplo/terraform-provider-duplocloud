---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "duplocloud_aws_ssm_parameter Resource - terraform-provider-duplocloud"
subcategory: ""
description: |-
  duplocloud_aws_ssm_parameter manages an AWS SSM parameter in Duplo.
---

# duplocloud_aws_ssm_parameter (Resource)

`duplocloud_aws_ssm_parameter` manages an AWS SSM parameter in Duplo.

## Example Usage

```terraform
resource "duplocloud_tenant" "myapp" {
  account_name = "myapp"
  plan_id      = "default"
}

resource "duplocloud_aws_ssm_parameter" "ssm_param" {
  tenant_id = duplocloud_tenant.myapp.tenant_id
  name      = "ssm_param"
  type      = "String"
  value     = "ssm_param_value"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) The name of the SSM parameter.
- **tenant_id** (String) The GUID of the tenant that the SSM parameter will be created in.
- **type** (String) The type of the SSM parameter.

### Optional

- **allowed_pattern** (String)
- **description** (String) The description of the SSM parameter.
- **key_id** (String)
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- **value** (String) The value of the SSM parameter.

### Read-Only

- **id** (String) The ID of this resource.
- **last_modified_date** (String)
- **last_modified_user** (String)

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)
- **update** (String)

## Import

Import is supported using the following syntax:

```shell
# Example: Importing an existing AWS SSM Parameter
#  - *TENANT_ID* is the tenant GUID
#  - *NAME* The name for the created Amazon SSM Parameter.
#
terraform import duplocloud_aws_ssm_parameter.ssm_param *TENANT_ID*/*NAME*
```
