terraform {
  required_providers {
    value = {
      source  = "github.com/pseudo-dynamic/value"
      version = "0.1.0"
    }
  }
}

resource "value_unknown_proposer" "default" {}

output "proposed_unknown" {
  value = value_unknown_proposer.default.value
}
