
packages {
  development = ["go@1.22.0", "delve@^1.22.0", "gotools@~0.7.0", "goreleaser@^1.24.0"]
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
  envVars      = null
}
