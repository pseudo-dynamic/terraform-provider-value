terraform {
  required_providers {
    value = {
      source  = "github.com/pseudo-dynamic/value"
      version = "0.1.0"
    }
  }
}

resource "value_promise" "default" {
  value = {
    timestamp = timestamp()
  }
}

output "result" {
  value = value_promise.default.result
}
