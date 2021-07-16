---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "duplocloud_gcp_storage_bucket Resource - terraform-provider-duplocloud"
subcategory: ""
description: |-
  duplocloud_gcp_storage_bucket manages a GCP storage bucket in Duplo.
---

# duplocloud_gcp_storage_bucket (Resource)

`duplocloud_gcp_storage_bucket` manages a GCP storage bucket in Duplo.

## Example Usage

```terraform
resource "duplocloud_tenant" "myapp" {
  account_name = "myapp"
  plan_id      = "default"
}

resource "duplocloud_gcp_storage_bucket" "mybucket" {

  tenant_id         = duplocloud_tenant.this.tenant_id
  name              = "mybucket"
  enable_versioning = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) The short name of the storage bucket.  Duplo will add a prefix to the name.  You can retrieve the full name from the `fullname` attribute.
- **tenant_id** (String) The GUID of the tenant that the storage bucket will be created in.

### Optional

- **enable_versioning** (Boolean) Whether or not versioning is enabled for the storage bucket. Defaults to `false`.
- **labels** (Map of String) The labels assigned to this storage bucket.
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- **fullname** (String) The full name of the storage bucket.
- **id** (String) The ID of this resource.
- **self_link** (String) The SelfLink of the storage bucket.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)
- **update** (String)

## Import

Import is supported using the following syntax:

```shell
# Example: Importing an existing GCP storage bucket
#  - *TENANT_ID* is the tenant GUID
#  - *SHORT_NAME* is the short name of the GCP storage bucket
#
terraform import duplocloud_gcp_storage_bucket.mybucket *TENANT_ID*/*SHORT_NAME*
```