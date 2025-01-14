---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "duplocloud_aws_mwaa_environment Resource - terraform-provider-duplocloud"
subcategory: ""
description: |-
  duplocloud_aws_mwaa_environment manages an AWS MWAA Environment resource in Duplo.
---

# duplocloud_aws_mwaa_environment (Resource)

`duplocloud_aws_mwaa_environment` manages an AWS MWAA Environment resource in Duplo.

## Example Usage

```terraform
resource "duplocloud_tenant" "myapp" {
  account_name = "myapp"
  plan_id      = "default"
}

data "duplocloud_tenant_aws_kms_key" "tenant_kms_key" {
  tenant_id = duplocloud_tenant.myapp.tenant_id
}

resource "duplocloud_aws_mwaa_environment" "my-mwaa" {
  tenant_id                       = duplocloud_tenant.myapp.tenant_id
  name                            = "airflow-test"
  source_bucket_arn               = "arn:aws:s3:::duploservices-demo01-dags-140563923322"
  dag_s3_path                     = "AirflowDags/dag"
  kms_key                         = data.duplocloud_tenant_aws_kms_key.tenant_kms_key.key_arn
  schedulers                      = 2
  max_workers                     = 10
  min_workers                     = 1
  airflow_version                 = "2.2.2"
  weekly_maintenance_window_start = "SUN:23:30"
  environment_class               = "mw1.small"

  airflow_configuration_options = {
    "core.log_format" : "[%%(asctime)s] {{%%(filename)s:%%(lineno)d}} %%(levelname)s - %%(message)s"
  }

  logging_configuration {
    dag_processing_logs {
      enabled   = false
      log_level = "INFO"
    }

    scheduler_logs {
      enabled   = false
      log_level = "INFO"
    }

    task_logs {
      enabled   = false
      log_level = "INFO"
    }

    webserver_logs {
      enabled   = false
      log_level = "INFO"
    }

    worker_logs {
      enabled   = false
      log_level = "INFO"
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **dag_s3_path** (String) The relative path to the DAG folder on your Amazon S3 storage bucket.
- **name** (String) The name of the Apache Airflow Environment.
- **source_bucket_arn** (String) The Amazon Resource Name (ARN) of your Amazon S3 storage bucket. For example, arn:aws:s3:::airflow-mybucketname.
- **tenant_id** (String) The GUID of the tenant that the Managed Workflows Apache Airflow will be created in.

### Optional

- **airflow_configuration_options** (Map of String, Sensitive) The `airflow_configuration_options` parameter specifies airflow override options
- **airflow_version** (String) Airflow version of your environment, will be set by default to the latest version that MWAA supports.
- **environment_class** (String) Environment class for the cluster. Possible options are `mw1.small`, `mw1.medium`, `mw1.large`.
- **execution_role_arn** (String) The Execution Role ARN of the Amazon MWAA Environment
- **kms_key** (String) The Amazon Resource Name (ARN) of your KMS key that you want to use for encryption. Will be set to the ARN of the managed KMS key aws/airflow by default.
- **logging_configuration** (Block List, Max: 1) (see [below for nested schema](#nestedblock--logging_configuration))
- **max_workers** (Number) The maximum number of workers that can be automatically scaled up. Value need to be between `1` and `25`.
- **min_workers** (Number) The minimum number of workers that you want to run in your environment.
- **plugins_s3_object_version** (String) The plugins.zip file version you want to use.
- **plugins_s3_path** (String) The relative path to the plugins.zip file on your Amazon S3 storage bucket. For example, plugins.zip. If a relative path is provided in the request, then `plugins_s3_object_version` is required.
- **requirements_s3_object_version** (String) The requirements.txt file version you want to use..
- **requirements_s3_path** (String) The relative path to the requirements.txt file on your Amazon S3 storage bucket. For example, requirements.txt. If a relative path is provided in the request, then requirements_s3_object_version is required.
- **schedulers** (Number) The number of schedulers that you want to run in your environment.
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- **wait_until_ready** (Boolean) Whether or not to wait until Amazon MWAA Environment to be ready, after creation. Defaults to `true`.
- **webserver_access_mode** (String) Specifies whether the webserver should be accessible over the internet or via your specified VPC.  Defaults to `PUBLIC_ONLY`.
- **weekly_maintenance_window_start** (String) Specifies the start date for the weekly maintenance window.

### Read-Only

- **arn** (String) The ARN of the Managed Workflows Apache Airflow.
- **fullname** (String) The full name provided by duplo for Apache Airflow Environment.
- **id** (String) The ID of this resource.
- **last_updated** (List of Object) (see [below for nested schema](#nestedatt--last_updated))
- **status** (String) The status of the Amazon MWAA Environment.
- **tags** (String) Tags.
- **webserver_url** (String) The webserver URL of the MWAA Environment.

<a id="nestedblock--logging_configuration"></a>
### Nested Schema for `logging_configuration`

Optional:

- **dag_processing_logs** (Block List, Max: 1) (see [below for nested schema](#nestedblock--logging_configuration--dag_processing_logs))
- **scheduler_logs** (Block List, Max: 1) (see [below for nested schema](#nestedblock--logging_configuration--scheduler_logs))
- **task_logs** (Block List, Max: 1) (see [below for nested schema](#nestedblock--logging_configuration--task_logs))
- **webserver_logs** (Block List, Max: 1) (see [below for nested schema](#nestedblock--logging_configuration--webserver_logs))
- **worker_logs** (Block List, Max: 1) (see [below for nested schema](#nestedblock--logging_configuration--worker_logs))

<a id="nestedblock--logging_configuration--dag_processing_logs"></a>
### Nested Schema for `logging_configuration.dag_processing_logs`

Optional:

- **enabled** (Boolean)
- **log_level** (String)


<a id="nestedblock--logging_configuration--scheduler_logs"></a>
### Nested Schema for `logging_configuration.scheduler_logs`

Optional:

- **enabled** (Boolean)
- **log_level** (String)


<a id="nestedblock--logging_configuration--task_logs"></a>
### Nested Schema for `logging_configuration.task_logs`

Optional:

- **enabled** (Boolean)
- **log_level** (String)


<a id="nestedblock--logging_configuration--webserver_logs"></a>
### Nested Schema for `logging_configuration.webserver_logs`

Optional:

- **enabled** (Boolean)
- **log_level** (String)


<a id="nestedblock--logging_configuration--worker_logs"></a>
### Nested Schema for `logging_configuration.worker_logs`

Optional:

- **enabled** (Boolean)
- **log_level** (String)



<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)


<a id="nestedatt--last_updated"></a>
### Nested Schema for `last_updated`

Read-Only:

- **created_at** (String)
- **error** (List of Object) (see [below for nested schema](#nestedobjatt--last_updated--error))
- **status** (String)

<a id="nestedobjatt--last_updated--error"></a>
### Nested Schema for `last_updated.error`

Read-Only:

- **error_code** (String)
- **error_message** (String)

## Import

Import is supported using the following syntax:

```shell
# Example: Importing an existing AWS MWAA Environment resource.
#  - *TENANT_ID* is the tenant GUID
#  - *FULL_NAME* is the fullname of the AWS MWAA Environment resource
#
terraform import duplocloud_aws_mwaa_environment.my-mwaa *TENANT_ID*/*FULL_NAME*
```
