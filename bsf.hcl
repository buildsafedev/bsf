
packages {
  development = ["go@1.21.4", "goreleaser@1.22.1", "gotools@0.7.0", "delve@1.21.2"]
  runtime     = ["cacert@3.95"]
}

gomodule {
  name       = "bsf"
  src        = "../."
  vendorHash = "sha256-f29THF+FPQ4ORx2SJ2EJVOGbvoEvF4V+V0ZkmjUY35o="
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