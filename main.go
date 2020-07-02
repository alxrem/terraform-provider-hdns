package main

import (
	"github.com/alxrem/terraform-provider-hdns/hdns"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: hdns.Provider})
}
