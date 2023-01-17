resource "haproxy-pf_bind" "bind1" {
  name = "bind1"
  address = "127.0.0.1"
  port = 8888
  parent_name = "frontend-name"
}