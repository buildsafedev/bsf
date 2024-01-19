
packages {
  development = ["athens@0.13.0", "go@1.21.4", "gotools@0.7.0", "delve@1.21.2", "goreleaser@1.23.0"]
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
  cmd          = ["/result/bin/bsf"]
  entrypoint   = null
  platform     = "linux/arm64"
}
