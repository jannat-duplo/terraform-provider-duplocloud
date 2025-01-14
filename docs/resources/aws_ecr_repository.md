---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "duplocloud_aws_ecr_repository Resource - terraform-provider-duplocloud"
subcategory: ""
description: |-
  duplocloud_aws_ecr_repository manages an aws ecr repository in Duplo.
---

# duplocloud_aws_ecr_repository (Resource)

`duplocloud_aws_ecr_repository` manages an aws ecr repository in Duplo.

## Example Usage

```terraform
resource "duplocloud_tenant" "myapp" {
  account_name = "myapp"
  plan_id      = "default"
}

resource "duplocloud_aws_ecr_repository" "test-ecr" {
  tenant_id                 = duplocloud_tenant.myapp.tenant_id
  name                      = "test-ecr"
  enable_scan_image_on_push = true
  enable_tag_immutability   = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) The name of the ECR Repository.
- **tenant_id** (String) The GUID of the tenant that the aws ecr repository will be created in.

### Optional

- **enable_scan_image_on_push** (Boolean) Indicates whether images are scanned after being pushed to the repository (true) or not scanned (false).
- **enable_tag_immutability** (Boolean) The tag mutability setting for the repository.
- **kms_encryption_key** (String) The ARN of the KMS key to use.
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- **arn** (String) Full ARN of the repository.
- **id** (String) The ID of this resource.
- **registry_id** (String) The registry ID where the repository was created.
- **repository_url** (String) The URL of the repository.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)

## Import

Import is supported using the following syntax:

```shell
# Example: Importing an existing AWS ECR repository
#  - *TENANT_ID* is the tenant GUID
#  - *SHORT_NAME* is the short name of the AWS ECR repository
#
terraform import duplocloud_aws_ecr_repository.myecr *TENANT_ID*/*SHORT_NAME*
```
