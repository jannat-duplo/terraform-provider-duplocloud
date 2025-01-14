---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "duplocloud_native_host_image Data Source - terraform-provider-duplocloud"
subcategory: ""
description: |-
  duplocloud_native_host_image retrieves details of a specific image for a given tenant.
---

# duplocloud_native_host_image (Data Source)

`duplocloud_native_host_image` retrieves details of a specific image for a given tenant.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **tenant_id** (String) The tenant ID

### Optional

- **is_kubernetes** (Boolean)
- **name** (String)

### Read-Only

- **id** (String) The ID of this resource.
- **image_id** (String)
- **k8s_version** (String)
- **os** (String)
- **region** (String)
- **tags** (List of Object) (see [below for nested schema](#nestedatt--tags))
- **username** (String)

<a id="nestedatt--tags"></a>
### Nested Schema for `tags`

Read-Only:

- **key** (String)
- **value** (String)


