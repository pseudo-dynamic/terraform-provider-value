terraform {
  required_providers {
    value = {
      source  = "github.com/pseudo-dynamic/value"
      version = "0.1.0"
    }
  }
}

data "value_temp_dir" "default" {}

output "temp_dir" {
  value = data.value_temp_dir.default.path
}
