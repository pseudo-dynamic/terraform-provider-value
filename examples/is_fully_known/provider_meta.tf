terraform {
  required_providers {
    value = {
      source  = "github.com/pseudo-dynamic/value"
      version = "0.1.0"
    }
  }

  provider_meta "value" {
    // {workdir} -> a placeholder (see docs)
    seed_prefix = "{workdir}"
  }
}