package routes

import "net/http"

type AuthRoutes struct {
	SignUp Route
	SignIn Route
}

var Auth = AuthRoutes{
	SignUp: Route{Path: "/api/auth/signup", Method: http.MethodPost},
	SignIn: Route{Path: "/api/auth/signin", Method: http.MethodPost},
}
