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

## Releasing the Provider

To release a new version of the provider:

1. Generate a GPG key for signing releases if you haven't already:
   ```sh
   gpg --full-generate-key
   ```
   - Choose RSA and RSA (default)
   - Choose 4096 bits
   - Enter your name and email
   - Set a secure passphrase

2. Export your GPG private key and note down your key fingerprint:
   ```sh
   # Get your key fingerprint
   gpg --list-secret-keys --keyid-format=long
   
   # Export the private key (replace [YOUR_KEY_ID] with your key ID)
   gpg --armor --export-secret-key [YOUR_KEY_ID]
   ```

3. Configure GitHub repository secrets:
   - Go to your repository Settings > Secrets and variables > Actions
   - Add the following secrets:
     - `GPG_PRIVATE_KEY`: Your exported GPG private key (the entire armored output)
     - `PASSPHRASE`: The passphrase you set for your GPG key

4. To create a new release:
   ```sh
   # Tag a new version (replace v0.1.0 with your version)
   git tag v0.1.0
   git push origin v0.1.0
   ```

   The GitHub Actions workflow will automatically:
   - Build the provider for all supported platforms
   - Sign the release with your GPG key
   - Create a GitHub release with the built artifacts
   - Generate release notes

### Release Requirements

For the automated release process to work:
- All tests must pass
- The version tag must start with 'v' (e.g., v0.1.0)
- Required environment variables must be set in GitHub Actions secrets
- The repository must have write permissions for GitHub Actions (already configured in workflow)

## License

This provider is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.