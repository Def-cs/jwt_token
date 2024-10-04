package api

var urls = make(map[string]*Handle)

func init() {
	urls["/register"] = newHandle(false, "/register", createUser)
	urls["/login"] = newHandle(false, "/login", login)
	urls["/refresh"] = newHandle(true, "/refresh", refresh)
}
