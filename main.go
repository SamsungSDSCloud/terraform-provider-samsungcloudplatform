package main

import (
	"flag"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

//go:generate terraform fmt -recursive ./examples/
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: scp.Provider,
		Debug:        debugMode,
		ProviderAddr: "samsungsds/scp",
	}

	plugin.Serve(opts)
}
