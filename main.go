package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"gitlab.com/alxrem/terraform-provider-hdns/hdns"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: hdns.Provider})
}
