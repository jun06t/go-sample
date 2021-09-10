package app.rbac

default allow = false

allow {
	user_is_admin
}

user_is_admin {
	input.roles[_] = "admin"
}

# GET /articles
allow {
	input.method = "GET"
	input.path = ["articles"]
	has_permission(input.roles, "myservice.article.list")
}

# GET /articles/:id
allow {
	some id
	input.method = "GET"
	input.path = ["articles", id]
	has_permission(input.roles, "myservice.article.get")
}

# POST /articles
allow {
	input.method = "POST"
	input.path = ["articles"]
	has_permission(input.roles, "myservice.article.create")
}

# PUT /articles/:id
allow {
	some id
	input.method = "PUT"
	input.path = ["articles", id]
	has_permission(input.roles, "myservice.article.update")
}

# DELETE /articles/:id
allow {
	some id
	input.method = "DELETE"
	input.path = ["articles", id]
	has_permission(input.roles, "myservice.article.delete")
}

has_permission(roles, p) {
	r := roles[_]
	data.role_permissions[r][_] == p
}
