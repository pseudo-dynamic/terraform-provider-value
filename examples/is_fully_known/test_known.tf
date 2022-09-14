// This example is incomplete. Please take a look at provider_meta.tf and shared.tf too!

resource "value_is_fully_known" "known" {
  value            = "test"
  unique_seed      = "known"
  proposed_unknown = value_unknown_proposer.default.value
}

output "is_known_value" {
  value = {
    is_fully_known = value_is_fully_known.known.result
  }
}