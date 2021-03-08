# terraform-provider-mock

This is an empty boilerplate repository for creating a terraform provider. 

The motivation for creating this repo was:

1. To learn how to create a terraform provider.
2. To be a _simple_ resource for others to do the same.
3. Demonstrate how to test a provider locally without needing to publish it.
4. Anything else that might be of interest.

## What is a terraform provider?

A 'provider' is an abstraction over an existing API, that will enable you to manage the creation of resources using terraform. If you don't know what terraform is, then I recommend reading up on that subject first.

In summary: if you have an API (or you are a user of an existing API), then you can manage that API via terraform.

## Things to know about this repository

This is quite literally a _skeleton_ repo. It's intentionally designed that way. Most tutorials online teach the details of a terraform provider by first implementing an API backend, but I personally find this an unnecessary mental hurdle. So I have avoided that in favour of heavily commented code that explains what you need to do, when, and why. This makes it easier for you to strip out what you don't want.

## Requirements for creating a terraform provider

1. Provider code
2. Provider documentation
3. A tagged release

If you intend on publishing a provider on `registry.terraform.io` you'll need to follow [these steps](https://www.terraform.io/docs/registry/providers/publishing.html) which includes generating documentation (for which I have: `make generate-docs` defined in this repo's `Makefile`). 

You'll also want to tag a commit to be used as the release version, which you'd then reference in the `version` field in your terraform code, for example a consumer of this provider should define something like the following:

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

> NOTE: when developing your own provider, remember not just to update the `source` value but also the parent key (in this case `mock`). I've forgotten to do this in the past and had it confuse me for hours because it's such a subtle thing to miss. 

## Linting a Provider

There is no official tool but `tfproviderlint` is written by a HashiCorp software engineer and has been used on many projects so is worth installing:

```bash
go install github.com/bflad/tfproviderlint/cmd/tfproviderlintx@latest
```

> NOTE: I suggest installing the 'extended' binary (notice `x` at the end of the name).

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
- Initialize your terraform project and then execute a plan (`terraform init && terraform plan`). 

> NOTE: every time you make a change to the terraform provider code, you'll need to rebuild the binary and then go to your consuming terraform project and reinitialize (i.e. `terraform init`) so it picks up the latest version of the `terraform-provider-mock` binary.

## Local Development

When running:

```bash
$ TF_LOG=TRACE terraform init
```

You should notice a couple of things different from what you'd normally see when initializing a new terraform project.

The first is a message highlighting the fact that a provider 'override' is in place:

```
Warning: Provider development overrides are in effect

The following provider development overrides are set in the CLI configuration:
 - integralist/mock in /Users/integralist/Code/terraform/terraform-provider-mock

The behavior may therefore not match any released version of the provider and
applying changes may cause the state to become incompatible with published
releases.
```

That is expected in this case we've followed the instructions above, which tells us how to implement an override for the sake of local testing of the provider code.

The other thing you'll notice is an error:

```
Error: Failed to query available provider packages

Could not retrieve the list of available versions for provider
integralist/mock: provider registry registry.terraform.io does not have a
provider named registry.terraform.io/integralist/mock

If you have just upgraded directly from Terraform v0.12 to Terraform v0.14
then please upgrade to Terraform v0.13 first and follow the upgrade guide for
that release, which might help you address this problem.
```

This error is expected because we've not actually published this provider to the terraform registry, so indeed it cannot be found. But the error doesn't prevent you from consuming the local provider binary still.

> NOTE: don't use Print functions from the `fmt` package in the terraform provider as depending on the execution flow terraform can treat it as input to its internal program and treat it as an error. So use Print functions from the `log` package instead.

## Terraform Execution Flow

When there is no terraform state file, then terraform won't execute any CRUD functions.

On the initial `terraform apply` you'll find CREATE is called first but what happens from there depends on how your provider works. For example, fastly and aws both call UPDATE at the end of the CREATE, where in this mock provider I call READ instead.

Once a terraform state file has been created, and you make a change to your terraform configuration file, you'll find the first operation called when running `terraform plan` is READ. This is because terraform wants to get the latest version of your infrastructure to compare against what you have defined locally in your configuration.

If you run `terraform apply` to ensure your changes are applied, then you'll find the first operation called by terraform is a READ. This is because if you don't have `terraform plan` set to save the 'execution plan' using the `-out` flag, then terraform is going to go off and get the latest data it can (you'll have noticed this as you would have had to type in "yes" manually to force the changes to be applied). After the READ, terraform calls UPDATE and what follows that is typically a READ because that's what most terraform providers do in their UPDATE function logic.

## Example Terraform Consumer Code

Below are two code files you can use to validate how to use this provider in its current form:

1. `service.tf`
2. `outputs.tf`

Here is the `service.tf` contents:

```tf
terraform {
  required_providers {
    mock = {
      source = "integralist/mock"
    }
  }
}

provider "mock" {
  foo = "example_value"
  #
  # if 'foo' wasn't set here by us, then the value would default to the value 
  # assigned to the environment variable 'MOCK_FOO' or the default value of nil
  # if the environment variable wasn't set.
}

resource "mock_example" "testing" {
  not_computed_required = "some value"

  dynamic "foo" {
    for_each = [{ number = 1 }, { number = 2 }, { number = 3 }]
    content {
      bar {
        number = foo.value.number
      }
    }
  }
  /*
   * The above is equivalent to:
   *
   * foo {
   *   bar {
   *     number = 1
   *   }
   * }
   * foo {
   *   bar {
   *     number = 2
   *   }
   * }
   * foo {
   *   bar {
   *     number = 3
   *   }
   * }
  */

  dynamic "baz" {
    // The variable inside the for_each block doesn't have to be the same as 
    // what you're assigning the value to.
    for_each = [{ something = "x" }, { something = "y" }, { something = "z" }]
    content {
      qux = baz.value.something
    }
  }
  /*
   * The above is equivalent to:
   *
   * baz {
   *   qux = "x"
   * }
   * baz {
   *   qux = "y"
   * }
   * baz {
   *   qux = "z"
   * }
  */

  some_list = ["a", "b", "c"]
}
```

Here is the `outputs.tf` contents:

```tf
output "last_updated" {
  value = mock_example.testing.last_updated
}
```

The `outputs.tf` is a terraform convention where you can specify what 'computed' values you would like to see displayed once a planned set of changes has been successfully applied.

Once you've written the above code, and you run a plan, you should see the following output:

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

If you were to run `terraform show` you would see `No state.` returned.

So let's run `terraform apply` to apply the 'planned' changes:

```bash
$ terraform apply

Warning: Provider development overrides are in effect

The following provider development overrides are set in the CLI configuration:
 - integralist/mock in /Users/integralist/Code/terraform/terraform-provider-mock

The behavior may therefore not match any released version of the provider and
applying changes may cause the state to become incompatible with published
releases.


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

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

mock_example.testing: Creating...
mock_example.testing: Creation complete after 0s [id=123]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

last_updated = "Saturday, 20-Feb-21 13:33:11 GMT"
```

Notice at the bottom of this output we see the `Outputs` section which is displaying what you had defined inside of `outputs.tf`.

If you were to now run `terraform show` you would see some state!

```bash
$ terraform show

# mock_example.testing:
resource "mock_example" "testing" {
    id                    = "123"
    last_updated          = "Saturday, 20-Feb-21 13:33:11 GMT"
    not_computed_required = "some value"
    some_list             = [
        "a",
        "b",
        "c",
    ]

    baz {
        qux = "x"
    }
    baz {
        qux = "y"
    }
    baz {
        qux = "z"
    }

    foo {
        bar {
            number  = 1
            version = "27356913-3cf2-4296-b78e-509d487f4fd0"
        }
    }
    foo {
        bar {
            number  = 2
            version = "8bd02c94-1e65-4eac-b106-f977c15ff173"
        }
    }
    foo {
        bar {
            number  = 3
            version = "b931c027-2cb0-463d-b289-f48ec2943a5e"
        }
    }
}


Outputs:

last_updated = "Saturday, 20-Feb-21 13:33:11 GMT"
```

## Reference Material

- [How Terraform Works](https://www.terraform.io/docs/extend/how-terraform-works.html): explains how providers are sourced, versioned and upgraded.
- [Schema Attributes and Types](https://www.terraform.io/docs/extend/schemas/schema-types.html): explains the various schema types you can define in your provider.
- [Writing a custom terraform provider](https://boxboat.com/2020/02/04/writing-a-custom-terraform-provider/): there actually isn't that many articles on the topic.
