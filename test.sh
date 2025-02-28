make build || exit 1
cp /Users/darnold/git/terraform-provider-openai/terraform-provider-openai ~/.terraform.d/plugins/registry.terraform.io/happypathway/openai/0.1.0/darwin_amd64/
rm -rf .terraform*
terraform init -upgrade && terraform apply -auto-approve || exit 1

