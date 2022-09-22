package goproviderconfig

func GetProposedUnknownAttributeDescription() string {
	return "It is very crucial that this field is **not** filled by any " +
		"custom value except the one produced by `value_unknown_proposer` (resource). " +
		"This has the reason as its `value` is **always** unknown during the plan phase. " +
		"On this behaviour this resource must rely and it cannot check if you do not so! "
}
