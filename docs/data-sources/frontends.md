---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "haproxy-pf_frontends Data Source - haproxy-pf"
subcategory: ""
description: |-
  
---

# haproxy-pf_frontends (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `frontends` (Attributes List) (see [below for nested schema](#nestedatt--frontends))

<a id="nestedatt--frontends"></a>
### Nested Schema for `frontends`

Optional:

- `default_backend` (String)
- `http_connection_mode` (String) possible values: httpclose,http-server-close,http-keep-alive
- `maxconn` (Number)
- `mode` (String)

Read-Only:

- `id` (String)
- `name` (String)


