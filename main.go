package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/integralist/terraform-provider-mock/mock"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: mock.Provider,
	})
}
