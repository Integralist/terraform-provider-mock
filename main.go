package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/integralist/terraform-provider-mock/mock"
)

// Article:
// https://boxboat.com/2020/02/04/writing-a-custom-terraform-provider/

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return mock.Provider()
		},
	})
}
