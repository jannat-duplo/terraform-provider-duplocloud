---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "duplocloud_plan_certificates Data Source - terraform-provider-duplocloud"
subcategory: ""
description: |-
  duplocloud_plan_certificates retrieves a list of cerificates for a given plan.
---

# duplocloud_plan_certificates (Data Source)

`duplocloud_plan_certificates` retrieves a list of cerificates for a given plan.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **plan_id** (String) The plan ID

### Read-Only

- **certificates** (List of Object) The list of certificates for this plan. (see [below for nested schema](#nestedatt--certificates))
- **id** (String) The ID of this resource.

<a id="nestedatt--certificates"></a>
### Nested Schema for `certificates`

Read-Only:

- **arn** (String)
- **name** (String)


