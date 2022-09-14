// This example is incomplete. Please take a look at provider_meta.tf and shared.tf too!

resource "value_is_fully_known" "unknown" {
  value = value_promise.default.result
  unique_seed      = "unknown"
  proposed_unknown = value_unknown_proposer.default.value
}

output "is_unknown_value" {
  value = {
    is_fully_known = value_is_fully_known.unknown.result
  }
}
