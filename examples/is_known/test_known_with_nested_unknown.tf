// This example is incomplete. Please take a look at provider_meta.tf and shared.tf too!

resource "value_is_known" "known_with_nested_unknown" {
  value = {
    nested = value_promise.default.result
  }

  unique_seed      = "nested_known"
  proposed_unknown = value_unknown_proposer.default.value
}

output "is_known_with_nested_unknown_value" {
  value = {
    known = value_is_known.known_with_nested_unknown.result
  }
}
