package haproxy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBindResource(t *testing.T) {
	frontendName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	bindName1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	bindName2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "haproxy-pf_frontend" "%s" {
					name = "%s"
					maxconn = 2000
					mode = "http"
					http_connection_mode = ""
				}
				resource "haproxy-pf_bind" "%s" {
					name = "%s"
					address = "127.0.0.1"
					port = 9999
					parent_name = "%s"
					depends_on = [
						haproxy-pf_frontend.%s
					]
				}
				`, frontendName, frontendName, bindName1, bindName1, frontendName, frontendName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_bind.%s", bindName1), "name", bindName1),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_bind.%s", bindName1), "port", "9999"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_bind.%s", bindName1), "address", "127.0.0.1"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(fmt.Sprintf("haproxy-pf_bind.%s", bindName1), "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      fmt.Sprintf("haproxy-pf_bind.%s", bindName1),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "haproxy-pf_frontend" "%s" {
					name = "%s"
					maxconn = 2000
					mode = "http"
					http_connection_mode = ""
				}
				resource "haproxy-pf_bind" "%s" {
					name = "%s"
					address = "127.0.0.1"
					port = 8888
					parent_name = "%s"
					depends_on = [
						haproxy-pf_frontend.%s
					]
				}
				`, frontendName, frontendName, bindName1, bindName1, frontendName, frontendName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_bind.%s", bindName1), "name", bindName1),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_bind.%s", bindName1), "port", "8888"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_bind.%s", bindName1), "address", "127.0.0.1"),
				),
			},
			// add two resources at the same time
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "haproxy-pf_frontend" "%s" {
					name = "%s"
					maxconn = 2000
					mode = "http"
					http_connection_mode = ""
				}
				resource "haproxy-pf_bind" "%s" {
					name = "%s"
					address = "127.0.0.1"
					port = 9998
					parent_name = "%s"
					depends_on = [
						haproxy-pf_frontend.%s
					]
				}
				resource "haproxy-pf_bind" "%s" {
					name = "%s"
					address = "127.0.0.1"
					port = 9999
					parent_name = "%s"
					depends_on = [
						haproxy-pf_frontend.%s
					]
				}
				`, frontendName, frontendName, bindName1, bindName1, frontendName, frontendName, bindName2, bindName2, frontendName, frontendName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_bind.%s", bindName1), "name", bindName1),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_bind.%s", bindName1), "port", "9998"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_bind.%s", bindName1), "address", "127.0.0.1"),
					resource.TestCheckResourceAttrSet(fmt.Sprintf("haproxy-pf_bind.%s", bindName1), "id"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_bind.%s", bindName2), "name", bindName2),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_bind.%s", bindName2), "port", "9999"),
					resource.TestCheckResourceAttr(fmt.Sprintf("haproxy-pf_bind.%s", bindName2), "address", "127.0.0.1"),
					resource.TestCheckResourceAttrSet(fmt.Sprintf("haproxy-pf_bind.%s", bindName2), "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
