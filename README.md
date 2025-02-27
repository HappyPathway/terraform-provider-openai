# Terraform Provider OpenAI

This Terraform provider enables interaction with OpenAI APIs to manage OpenAI resources.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
```sh
git clone git@github.com:HappyPathway/terraform-provider-openai
```

2. Enter the repository directory
```sh
cd terraform-provider-openai
```

3. Build the provider
```sh
make build
```

## Using the provider

You can use the provider by adding it to your Terraform configuration:

```hcl
terraform {
  required_providers {
    openai = {
      source = "HappyPathway/openai"
    }
  }
}

provider "openai" {
  api_key = var.openai_api_key
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.19+ is *required*).

To compile the provider:

```sh
make build
```

To run the tests:

```sh
make test
```

### Running Tests

The provider has both unit tests and acceptance tests. You can run them using:

```sh
make test
make testacc
```

*Note:* Acceptance tests create real resources and often cost money to run.

## License

This provider is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.