// This example is complete but there are additional features implemented in terraform.tf!

resource "value_unknown_proposer" "known" {}

resource "value_is_known" "known" {
  value            = "test"
  guid_seed        = "known"
  proposed_unknown = value_unknown_proposer.known.value
}

output "is_known_value" {
  value = {
    known = value_is_known.known.result
  }
}
