// This example is complete but there are additional features implemented in terraform.tf!

resource "value_unknown_proposer" "known_with_nested_unknown" {}

resource "value_promise" "known_with_nested_unknown" {
  value = "test"
}

resource "value_is_fully_known" "known_with_nested_unknown" {
  value = {
    nested = value_promise.known_with_nested_unknown.result
  }

  guid_seed        = "nested_unknown"
  proposed_unknown = value_unknown_proposer.known_with_nested_unknown.value
}

output "is_known_with_nested_unknown_value" {
  value = {
    fully_known = value_is_fully_known.known_with_nested_unknown.result
  }
}
