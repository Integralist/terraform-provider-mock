# terraform-provider-mock

This is an empty boilerplate repository for creating a terraform provider. It doesn't even include the code for generating an HTTP client which you would want to do in order to make the relevant API requests to a service for which your terraform provider would be designed to manage (although I have plenty of code comments to help demonstrate that).

## Requirements for creating a terraform provider

1. Code (inc. Vendors)
2. Docs
3. Tagged Release

When creating a terraform provider you'll need to write the provider/plugin code and run `go mod vendor` to ensure all dependant code (e.g. API clients etc are bundled into the compiled `terraform-provider-mock` binary). If you intend on publishing the provider on `registry.terraform.io` you'll need to follow [these steps](https://www.terraform.io/docs/registry/providers/publishing.html) which includes generating documentation (i.e. `make generate-docs`). you'll also want to tag a commit to be used as the release version (which you'd then reference in the `version` field in your terraform code which imports this provider, for example:

```tf
terraform {
  required_providers {
    mock = {
      source = "integralist/mock"
      version = "<tag_version_here>"
    }
  }
}
```

## How to use this provider

To consume this provider without it being published to the terraform registry, follow these steps:

- Clone this repo and build the `terraform-provider-mock` binary:  
  ```bash
  make build
  ```
- Create a separate directory for your own terraform project (e.g. `/example-tf`).
- Create a `dev.tfrc` file in your terraform directory:
  ```tf
  provider_installation {
    dev_overrides {
      "integralist/mock" = "/path/to/terraform-provider-mock" // wherever the binary is.
    }
    direct {}
  }
  ```
- Set `TF_CLI_CONFIG_FILE` environment variable (e.g. `export TF_CLI_CONFIG_FILE=/example-tf/dev.tfrc`).
- 

```bash
$ TF_LOG=trace terraform init

Initializing the backend...

Initializing provider plugins...
- Finding latest version of integralist/mock...

Warning: Provider development overrides are in effect

The following provider development overrides are set in the CLI configuration:
 - integralist/mock in /Users/integralist/Code/terraform/terraform-provider-mock

The behavior may therefore not match any released version of the provider and
applying changes may cause the state to become incompatible with published
releases.


Error: Failed to query available provider packages

Could not retrieve the list of available versions for provider
integralist/mock: provider registry registry.terraform.io does not have a
provider named registry.terraform.io/integralist/mock

If you have just upgraded directly from Terraform v0.12 to Terraform v0.14
then please upgrade to Terraform v0.13 first and follow the upgrade guide for
that release, which might help you address this problem.
```

```bash
$ terraform plan

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # mock_example.testing will be created
  + resource "mock_example" "testing" {
      + id                    = (known after apply)
      + last_updated          = (known after apply)
      + not_computed_required = "some value"
      + some_list             = [
          + "a",
          + "b",
          + "c",
        ]

      + baz {
          + qux = "x"
        }
      + baz {
          + qux = "y"
        }
      + baz {
          + qux = "z"
        }

      + foo {
          + bar {
              + number  = 1
              + version = (known after apply)
            }
        }
      + foo {
          + bar {
              + number  = 2
              + version = (known after apply)
            }
        }
      + foo {
          + bar {
              + number  = 3
              + version = (known after apply)
            }
        }
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + last_updated = (known after apply)

------------------------------------------------------------------------

Note: You didn't specify an "-out" parameter to save this plan, so Terraform
can't guarantee that exactly these actions will be performed if
"terraform apply" is subsequently run.
```
