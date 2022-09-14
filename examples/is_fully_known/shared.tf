// An "(known after apply)" value producer
resource "value_unknown_proposer" "default" {}

resource "value_promise" "default" {
  value = "test"
}
