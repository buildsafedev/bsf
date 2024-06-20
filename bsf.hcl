
packages {
  development = ["delve@1.22.0", "gotools@0.16.1", "go-task@~3.37.2", "go@1.21.6"]
  runtime     = ["cacert@3.95"]
}

gomodule {
  name    = "bsf"
  src     = "./."
  ldFlags = null
  tags    = null
  doCheck = false
}
