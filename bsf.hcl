
packages {
  development = ["gotools@0.7.0", "delve@1.21.2", "go@1.21.4"]
  runtime     = ["cacert@3.95"]
}

gomodule {
  name       = "github.com/buildsafedev/bsf"
  src        = "../."
  vendorHash = ""
}
