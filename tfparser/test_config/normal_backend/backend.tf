terraform {
  required_version = ">= 0.13.0, < 0.13.2"
  token            = "EmN5pXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  backend "remote" {
    hostname     = "app.terraform.io"
    organization = "test-org"

    workspace {
      name = "test-workspace"
    }
  }
}
