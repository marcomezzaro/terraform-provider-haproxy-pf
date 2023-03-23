# Terraform Provider Haproxy

This provider uses Haproxy Dataplane api in order to manage Terraform objects.

Built with Terraform plugin framework.

Checklist:

- [x] Backend
- [x] Frontend
- [x] Bind
- [x] Server
- [x] Server Template

TODO:

- [ ] ACL
- [ ] More resource options

## Dev Build provider

set terraformrc file with dev_overrides

```
# ~/.terraformrc 
provider_installation {

  dev_overrides {
      "registry.terraform.io/marcomezzaro/haproxy-pf" = "<GOPATH>/bin/"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Run the following command to build the provider

```shell
$ make install
```

## build and run haproxy for testing

```bash
podman build -t haproxy-dp -f docker/Dockerfile docker/
podman run -it --rm -p 5555:5555 --name haproxy localhost/haproxy-dp:latest
```

## Run ACC Test

```shell
TF_ACC=1 HAPROXY_SERVER="localhost:5555" HAPROXY_USERNAME="admin" HAPROXY_PASSWORD="adminpwd" HAPROXY_INSECURE="true" go test -v -cover -count 1 ./haproxy/
```

## Generate doc

```shell
go generate ./...
```

## Test sample configuration

First, build and install the provider.

```shell
$ make install
```

Then, navigate to the `examples` directory. and use the example to create terraform manifests

```shell
$ cd examples
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform apply
```

## Special Thanks

https://github.com/matthisholleville