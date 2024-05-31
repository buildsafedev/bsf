
packages {
  development = ["go@1.22.0", "delve@^1.22.0", "gotools@~0.7.0", "goreleaser@^1.24.0"]
  runtime     = ["cacert@3.95"]
}

gomodule {
  name       = "bsf"
  src        = "./."
  doCheck = false
}

