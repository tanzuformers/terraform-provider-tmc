package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/tanzuformers/terraform-provider-tmc/tmc"
)

// Run "go generate" to format example terraform files and generate the docs for
// the registry/website

// If you do not have terraform installed, you can remove the formatting command
// but its suggested to ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on
// how it works and how docs can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	// remove date and time stamp from log output as the plugin SDK already adds its own
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	var debugMode bool

	flag.BoolVar(&debugMode, "debuggable", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return tmc.Provider()
		},
	}

	if debugMode {
		err := plugin.Debug(context.Background(),
			"registry.terraform.io/tanzuformers/terraform-provider-tmc",
			opts)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		plugin.Serve(opts)
	}
}
