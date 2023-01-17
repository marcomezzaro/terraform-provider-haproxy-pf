resource "haproxy-pf_backend" "example" {
  name    = "b1"
  balance = "roundrobin"
  mode    = "http"
}