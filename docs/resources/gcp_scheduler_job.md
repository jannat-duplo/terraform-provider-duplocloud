---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "duplocloud_gcp_scheduler_job Resource - terraform-provider-duplocloud"
subcategory: ""
description: |-
  duplocloud_gcp_scheduler_job manages a GCP scheduler job in Duplo.
---

# duplocloud_gcp_scheduler_job (Resource)

`duplocloud_gcp_scheduler_job` manages a GCP scheduler job in Duplo.

## Example Usage

```terraform
resource "duploscheduler_tenant" "myapp" {
  account_name = "myapp"
  plan_id      = "default"
}

// A simple scheduler job with an HTTPS target, running at 9 am daily.
resource "duploscheduler_gcp_scheduler_job" "myjob" {
  tenant_id = local.tenant_id

  name = "myjob"

  schedule = "* 9 * * *"
  timezone = "America/New_York"

  http_target {
    method = "GET"
    uri    = "https://www.google.com"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) The short name of the scheduler job.  Duplo will add a prefix to the name.  You can retrieve the full name from the `fullname` attribute.
- **schedule** (String) The desired schedule, in cron format.
- **tenant_id** (String) The GUID of the tenant that the scheduler job will be created in.
- **timezone** (String) The timezone used to determine the schedule, in UNIX format

### Optional

- **app_engine_target** (Block List, Max: 1) Specifies an App Engine target for the scheduler job. (see [below for nested schema](#nestedblock--app_engine_target))
- **attempt_deadline** (String) The attempt deadline for the scheduler job.
- **description** (String) The description of the scheduler job.
- **http_target** (Block List, Max: 1) Specifies an HTTP target for the scheduler job. (see [below for nested schema](#nestedblock--http_target))
- **pubsub_target** (Block List, Max: 1) Specifies a pubsub target for the scheduler job. (see [below for nested schema](#nestedblock--pubsub_target))
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- **fullname** (String) The full name of the scheduler job.
- **id** (String) The ID of this resource.
- **self_link** (String) The SelfLink of the scheduler job.

<a id="nestedblock--app_engine_target"></a>
### Nested Schema for `app_engine_target`

Required:

- **method** (String) The HTTP method to use.
- **relative_uri** (String) The relative URI.

Optional:

- **body** (String) The HTTP request body to send.
- **headers** (Map of String) The HTTP headers to send.
- **routing** (Block List, Max: 1) Specifies App Engine routing. (see [below for nested schema](#nestedblock--app_engine_target--routing))

<a id="nestedblock--app_engine_target--routing"></a>
### Nested Schema for `app_engine_target.routing`

Optional:

- **host** (String) The App Engine host.
- **instance** (String) The App Engine instance.
- **service** (String) The App Engine service.
- **version** (String) The App Engine service version.



<a id="nestedblock--http_target"></a>
### Nested Schema for `http_target`

Required:

- **method** (String) The HTTP method to use.
- **uri** (String) The request URI.

Optional:

- **body** (String) The HTTP request body to send.
- **headers** (Map of String) The HTTP headers to send.
- **oauth_token** (Block List, Max: 1) Specifies OAuth authentication. (see [below for nested schema](#nestedblock--http_target--oauth_token))
- **oidc_token** (Block List, Max: 1) Specifies OIDC authentication. (see [below for nested schema](#nestedblock--http_target--oidc_token))

<a id="nestedblock--http_target--oauth_token"></a>
### Nested Schema for `http_target.oauth_token`

Optional:

- **enabled** (Boolean) Must be set to `true`. Defaults to `true`.
- **scope** (String) The OAuth token scope.

Read-Only:

- **service_account_email** (String) The OAuth token service account email.


<a id="nestedblock--http_target--oidc_token"></a>
### Nested Schema for `http_target.oidc_token`

Optional:

- **audience** (String) The OIDC token audience.
- **enabled** (Boolean) Must be set to `true`. Defaults to `true`.

Read-Only:

- **service_account_email** (String) The OIDC token service account email.



<a id="nestedblock--pubsub_target"></a>
### Nested Schema for `pubsub_target`

Required:

- **topic_name** (String) The name of the topic to target

Optional:

- **attributes** (Map of String) The attributes to send to the pubsub target.
- **data** (String) The data to send to the pubsub topic.


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)
- **update** (String)

## Import

Import is supported using the following syntax:

```shell
# Example: Importing an existing GCP scheduler job
#  - *TENANT_ID* is the tenant GUID
#  - *SHORT_NAME* is the short name of the GCP scheduler job
#
terraform import duploscheduler_gcp_scheduler_job.myjob *TENANT_ID*/*SHORT_NAME*
```