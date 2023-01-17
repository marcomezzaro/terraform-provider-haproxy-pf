resource "haproxy-pf_server" "s1" {
  name = "s1"
  address = "127.0.0.1"
  port = 8989
  check = "disabled"
  parent_name = "backend-name"
  
}