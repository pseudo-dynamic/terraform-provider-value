terraform {
  required_providers {
    value = {
      source  = "github.com/pseudo-dynamic/value"
      version = "0.1.0"
    }
  }

  provider_meta "value" {
    // Module-scoped seed addition.
    // {workdir} -> a placeholder (see docs)
    guid_seed_addition = "module(is-fully-known)"
  }
}

provider "value" {
  // Project-wide seed addition.
  // Won't overwrite module-scoped seed addition,
  // instead both serve are now considered as seed addition.
  guid_seed_addition = "{workdir}"
}
