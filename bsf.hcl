
packages {
  development = ["gotools@0.7.0", "delve@1.21.2", "goreleaser@1.21.1", "go@1.21.4"]
  runtime     = ["cacert@3.95"]
}

gomodule {
  name       = "bsf"
  src        = "../."
  vendorHash = "sha256-x/9O7qBzA1PaHBmZEzd8Kt0XuSw6feuJnd9FVUbdTt4="
  doCheck    = false
}

export "dev" {
  artifactType = "oci"
  name         = "ttl.sh/bsfdev/bsf:dev"
  cmd          = ["/result/bin/bsf \n"]
  entrypoint   = null
  platform     = "linux/arm64"
}
