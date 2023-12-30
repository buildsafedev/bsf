
packages {
  development = ["go@1.21.4", "goreleaser@1.22.1", "gotools@0.7.0", "delve@1.21.2"]
  runtime     = ["cacert@3.95"]
}

gomodule {
  name       = "bsf"
  src        = "../."
  vendorHash = "sha256-f29THF+FPQ4ORx2SJ2EJVOGbvoEvF4V+V0ZkmjUY35o="
}
