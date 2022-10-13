package authz

default allow = false

# allow all dapr requests
allow {
	contains(input.method, "/dapr.proto.runtime.v1.AppCallback/")
}

allow {
	method_perms := input.methodPerms[_]
	user_perms := input.userPerms[_]

	# check if any user permission is in the set of method permissions
	user_perms == method_perms
}