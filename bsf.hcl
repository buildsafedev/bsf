
packages {
  development = ["go@1.21.4", "goreleaser@1.22.1", "gotools@0.7.0", "delve@1.21.2"]
  runtime     = ["cacert@3.95"]
}

gomodule {
  name       = "bsf"
  src        = "../."
  vendorHash = "sha256-x/9O7qBzA1PaHBmZEzd8Kt0XuSw6feuJnd9FVUbdTt4="
}


export "dev"{
  name = "ttl.sh/bsfdev/bsf:dev"
  artifactType = "oci"
  platform = "linux/arm64"
  cmd = ["/result/bin/bsf"]
}

export "prod"{
  name = "ttl.sh/bsfdev/bsf:prod"
  artifactType = "oci"
  platform = "linux/arm64"
  cmd = ["/result/bin/bsf"]
  publish = true
}