package haproxy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccServerResource(t *testing.T) {
	backendName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	serverName1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	serverName2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			// TODO: remove providerConfig and use global env
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "haproxy-pf_backend" "%s" {
					name = "%s"
					balance = "leastconn"
					mode    = "tcp"
				}
				resource "haproxy-pf_server" "%s" {
					name = "%s"
					address = "127.0.0.1"
					port = 9999
					check = "enabled"
					parent_name = "%s"
					depends_on = [
						haproxy-pf_backend.%s
					]
				}
				`, backendName, backendName, serverName1, serverName1, backendName, backendName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server.%s", serverName1), "name", serverName1),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server.%s", serverName1), "port", "9999"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server.%s", serverName1), "address", "127.0.0.1"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(fmt.Sprintf("haproxy-pf_server.%s", serverName1), "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      fmt.Sprintf("haproxy-pf_server.%s", serverName1),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "haproxy-pf_backend" "%s" {
					name = "%s"
					balance = "leastconn"
					mode    = "tcp"
				}
				resource "haproxy-pf_server" "%s" {
					name = "%s"
					address = "127.0.0.1"
					port = 8888
					check = "disabled"
					parent_name = "%s"
					depends_on = [
						haproxy-pf_backend.%s
					]
				}
				`, backendName, backendName, serverName1, serverName1, backendName, backendName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server.%s", serverName1), "name", serverName1),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server.%s", serverName1), "port", "8888"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server.%s", serverName1), "address", "127.0.0.1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server.%s", serverName1), "check", "disabled"),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "haproxy-pf_backend" "%s" {
					name = "%s"
					balance = "leastconn"
					mode    = "tcp"
				}
				resource "haproxy-pf_server" "%s" {
					name = "%s"
					address = "127.0.0.1"
					port = 9998
					check = "enabled"
					parent_name = "%s"
					depends_on = [
						haproxy-pf_backend.%s
					]
				}
				resource "haproxy-pf_server" "%s" {
					name = "%s"
					address = "127.0.0.1"
					port = 9999
					check = "enabled"
					parent_name = "%s"
					depends_on = [
						haproxy-pf_backend.%s
					]
				}
				`, backendName, backendName, serverName1, serverName1, backendName, backendName, serverName2, serverName2, backendName, backendName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server.%s", serverName1), "name", serverName1),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server.%s", serverName1), "port", "9998"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server.%s", serverName1), "address", "127.0.0.1"),
					resource.TestCheckResourceAttrSet(fmt.Sprintf("haproxy-pf_server.%s", serverName1), "id"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server.%s", serverName2), "name", serverName2),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server.%s", serverName2), "port", "9999"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_server.%s", serverName2), "address", "127.0.0.1"),
					resource.TestCheckResourceAttrSet(fmt.Sprintf("haproxy-pf_server.%s", serverName2), "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
