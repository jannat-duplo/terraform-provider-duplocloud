terraform {
  required_providers {
    duplocloud = {
      version = "0.4.10" # RELEASE VERSION
      source = "registry.terraform.io/duplocloud/duplocloud"
    }
    aws = {
      version = "~> 3.29.1"
    }
  }
}

provider "duplocloud" {
  // duplo_host = "https://xxx.duplocloud.net"  # you can also set the duplo_host env var
  // duplo_token = ".."                         # please *ONLY* specify using a duplo_token env var (avoid checking secrets into git)
}

variable "plan_id" {
  type = string
  default = "default"
}

variable "tenant_id" {
  type = string
}

# AWS region retrieval - straight from Duplo!
data "duplocloud_tenant_aws_region" "test" { tenant_id = var.tenant_id }
output "aws_region" { value = data.duplocloud_tenant_aws_region.test.aws_region }

# Use any AWS terraform resource with just-in-time Duplo credentials!
data "duplocloud_tenant_aws_credentials" "test" { tenant_id = var.tenant_id }
provider "aws" {
  # The following credentials are temporary "just in time" credentials created by Duplo.
  access_key = data.duplocloud_tenant_aws_credentials.test.access_key_id
  secret_key = data.duplocloud_tenant_aws_credentials.test.secret_access_key
  token      = data.duplocloud_tenant_aws_credentials.test.session_token
  region     = data.duplocloud_tenant_aws_credentials.test.region
}

# For example, get the account ID
data "aws_caller_identity" "current" {}
output "aws_account_id" { value = data.aws_caller_identity.current.account_id }

# Or, apply additional policy to an S3 bucket created by Duplo.
resource "duplocloud_s3_bucket" "test" {
  tenant_id = var.tenant_id
  name = "joetest"
}
resource "aws_s3_bucket_public_access_block" "test" {
  bucket = duplocloud_s3_bucket.test.fullname
  block_public_acls = true
  ignore_public_acls = true
  block_public_policy = true
  restrict_public_buckets = true
}

# Or, get information on one of your KMS keys
#data "duplocloud_tenant_aws_kms_keys" "test" { tenant_id = var.tenant_id }
#data "aws_kms_key" "test" { key_id = data.duplocloud_tenant_aws_kms_keys.test.keys[1].key_id }
#output "kms_key_state" { value = data.aws_kms_key.test.key_state }

# Or, get information on your Duplo tenant's KMS key
#data "duplocloud_tenant_aws_kms_key" "test" { tenant_id = var.tenant_id }
#data "aws_kms_key" "test" { key_id = data.duplocloud_tenant_aws_kms_key.test.key_id }
#output "kms_key_state" { value = data.aws_kms_key.test.key_state }
