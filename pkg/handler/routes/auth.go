package routes

import (
	"net/http"
	"tic-tac-toe/pkg/common"
)

type AuthRoutes struct {
	Register common.ApiRoute
	Login    common.ApiRoute
}

var Auth = AuthRoutes{
	Register: common.ApiRoute{
		Path:   "/api/auth/register",
		Method: http.MethodPost,
	},
	Login: common.ApiRoute{
		Path:   "/api/auth/login",
		Method: http.MethodPost,
	},
}
