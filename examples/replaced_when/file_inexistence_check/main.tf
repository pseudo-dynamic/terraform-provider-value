terraform {
  required_providers {
    value = {
      source  = "github.com/pseudo-dynamic/value"
      version = "0.1.0"
    }
  }
}

resource "value_unknown_proposer" "default" {}

locals {
  files = {
    "findme" : {
      fullname = "${path.module}/findme"
    }
  }
}

resource "value_path_exists" "findme" {
  path             = local.files["findme"].fullname
  guid_seed        = "findme"
  proposed_unknown = value_unknown_proposer.default.value
}

resource "value_replaced_when" "findme_inexistence" {
  condition = !value_path_exists.findme.exists
}

resource "local_file" "findme" {
  count    = !value_path_exists.findme.exists ? 1 : 0
  content  = ""
  filename = local.files["findme"].fullname
}

output "is_findme_inexistent" {
  value = !value_path_exists.findme.exists
}

output "findme_inexistence_caused_new_value" {
  value = value_replaced_when.findme_inexistence.value
}
