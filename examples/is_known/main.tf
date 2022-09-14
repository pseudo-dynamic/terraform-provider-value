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

resource "value_is_known" "unknown" {
  value = value_promise.default.result
}

resource "value_is_known" "nested_unknown" {
  value = {
    nested = value_promise.default.result
  }
}

resource "value_is_known" "known" {
  value = "test"
}

output "is_unknown_value" {
  value = {
    known = value_is_known.unknown.result
  }
}

output "is_nested_unknown" {
  value = {
    known = value_is_known.nested_unknown.result
  }
}

output "is_known_value" {
  value = {
    known = value_is_known.known.result
  }
}
