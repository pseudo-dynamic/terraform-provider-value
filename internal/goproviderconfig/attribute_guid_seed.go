package goproviderconfig

func GetGuidSeedAttributeDescription(resourceName string) string {
	return "Attention! The seed is being used to determine resource uniqueness prior (first plan phase) " +
	"and during apply phase (second plan phase). Very important to state is that the **seed must be fully " +
	"known during the plan phase**, otherwise, an error is thrown. Within one terraform plan & apply the " +
	"**seed of every \"" + resourceName + "\" must be unique**! I really recommend you to use the " +
	"provider configuration and/or provider_meta configuration to increase resource uniqueness. " +
	"Besides `guid_seed`, the provider configuration seed, the provider_meta configuration seed and " +
	"the resource type itself will become part of the final seed. Under certain circumstances you " +
	"may face problems if you run terraform concurrenctly. If you do so, then I recommend you to " +
	"pass-through a random value via a user (environment) variable that you then add to the seed."
}