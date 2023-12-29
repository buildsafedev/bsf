
packages {
  development = ["gotools@0.7.0", "delve@1.21.2", "go@1.21.4"]
  runtime     = ["cacert@3.95"]
}

gomodule {
  name       = "github.com/buildsafedev/bsf"
  src        = "../."
  vendorHash = "sha256-f29THF+FPQ4ORx2SJ2EJVOGbvoEvF4V+V0ZkmjUY35o="
}
