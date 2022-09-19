// This example is complete but there are additional features implemented in terraform.tf!

resource "value_unknown_proposer" "unknown" {}

resource "value_promise" "unknown" {
  value = "test"
}

resource "value_is_known" "unknown" {
  value            = value_promise.unknown.result
  guid_seed        = "unknown"
  proposed_unknown = value_unknown_proposer.unknown.value
}

output "is_unknown_value" {
  value = {
    known = value_is_known.unknown.result
  }
}
