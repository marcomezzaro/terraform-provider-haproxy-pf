package main

import (
	"context"
	"flag"
	"log"
	"terraform-provider-haproxy-pf/haproxy"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name haproxy-pf

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/marcomezzaro/haproxy-pf",
		Debug:   debug,
	}
	err := providerserver.Serve(context.Background(), haproxy.New, opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
