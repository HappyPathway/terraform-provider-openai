package main

import (
	"github.com/HappyPathway/terraform-provider-openai/openai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: openai.Provider,
	})
}
