terraform {
  required_providers {
    value = {
      source  = "github.com/pseudo-dynamic/value"
      version = "0.1.0"
    }
  }
}

resource "value_replaced_when" "findme_inexistence" {
  for_each = {
    "findme" : {
      fullname = "${path.module}/findme"
    }
  }

  condition = !fileexists(each.value.fullname)
}

output "findme_inexistence_check" {
  // On craetion this should contain a non-null value.value
  value = value_replaced_when.findme_inexistence
}
