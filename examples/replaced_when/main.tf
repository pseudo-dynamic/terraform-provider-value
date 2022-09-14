terraform {
  required_providers {
    value = {
      source  = "github.com/pseudo-dynamic/value"
      version = "0.1.0"
    }
  }
}

# resource "value_promise" "bool" {
#   value = true
# }

resource "value_replaced_when" "true" {
  count     = 0
  condition = true
}

resource "value_stash" "replacement_trigger" {
  // This workaround detects not only resource creation
  // and every change of value but also the deletion of
  // resource.
  value = try(value_replaced_when.true[0].value, null)
}

resource "value_stash" "replaced" {
  lifecycle {
    // Replace me whenever value_replaced_when.true[0].value 
    // changes or value_replaced_when[0] gets deleted.
    replace_triggered_by = [
      value_stash.replacement_trigger.value
    ]
  }
}
