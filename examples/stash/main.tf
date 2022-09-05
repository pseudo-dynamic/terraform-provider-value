terraform {
  required_providers {
    value = {
      source  = "github.com/pseudo-dynamic/value"
      version = "0.1.0"
    }
  }
}

resource "value_stash" "default" {
  value = [
    {
      coffee = {
        id = 3
      }

      quantity = 9
    },
    {
      coffee = {
        id = 1
      }

      quantity = 2
    }
  ]
}

output "value" {
  value = value_stash.default.value
}
