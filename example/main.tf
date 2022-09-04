terraform {
  required_providers {
    value = {
      source  = "github.com/teneko/value"
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

  lifecycle {
    ignore_changes = all
  }
}

output "edu_order" {
  value = value_stash.default.value
}
