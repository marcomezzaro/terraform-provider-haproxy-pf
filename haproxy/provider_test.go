package haproxy

import (
	"os"
	"strconv"
	"terraform-provider-haproxy-pf/haproxy/middleware"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Haproxy client is properly configured.
	// It is also possible to use the HAPROXY_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
provider "haproxy-pf" {
  username = "admin"
  password = "adminpwd"
	host     = "localhost:5555"
	insecure = true
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"haproxy-pf": providerserver.NewProtocol6WithError(New()),
	}
	testAccProviders map[string]*haproxyProvider
	testAccProvider  *haproxyProvider
)

func TestAccMain(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		// short circuit non-acceptance test runs
		os.Exit(1)
	}

	serverAddr := os.Getenv("HAPROXY_SERVER")
	username := os.Getenv("HAPROXY_USERNAME")
	password := os.Getenv("HAPROXY_PASSWORD")
	insecure, _ := strconv.ParseBool(os.Getenv("HAPROXY_INSECURE"))

	testClient := middleware.NewClient(username, password, serverAddr, insecure)

	err := testClient.TestApiCall()
	if err != nil {
		panic(err)
	}
}
