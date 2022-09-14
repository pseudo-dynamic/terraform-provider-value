terraform {
  required_providers {
    value = {
      source  = "github.com/pseudo-dynamic/value"
      version = "0.1.0"
    }
  }
}

resource "value_promise" "default" {
  value = "test"
}

resource "value_is_fully_known" "unknown" {
  value = value_promise.default.result
}

resource "value_is_fully_known" "nested_unknown" {
  value = {
    nested = value_promise.default.result
  }
}

resource "value_is_fully_known" "known" {
  value = "test"
}

output "is_unknown_value" {
  value = {
    fully_known = value_is_fully_known.unknown.result
  }
}

output "is_nested_unknown" {
  value = {
    fully_known = value_is_fully_known.nested_unknown.result
  }
}

output "is_known_value" {
  value = {
    fully_known = value_is_fully_known.known.result
  }
}
