
packages {
  development = ["go@1.21.4", "gotools@0.7.0", "delve@1.21.2"]
  runtime     = ["cacert"]
}

gomodule {
  name       = "github.com/buildsafedev/bsf"
  src        = "../."
  vendorHash = ""
}
