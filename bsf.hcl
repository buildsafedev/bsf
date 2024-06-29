
packages {
  development = ["go@1.21.6", "gotools@0.16.1", "delve@1.22.0", "go-task@~3.37.2"]
  runtime     = ["cacert@3.95"]
}

gomodule {
  name    = "bsf"
  src     = "./."
  ldFlags = null
  tags    = null
  doCheck = false
}


githubRelease "bsf" {
  owner = "buildsafedev"
  repo  = "bsf"
}