package api

import (
	"fmt"

	"github.com/ogen-go/ogen/middleware"
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/infra/config"
)

var operationRoles = map[string]string{
	/*
		"addIngredient": "admin",
		"addUnit":       "admin",
		"getUnits":      "admin",
	*/
}

func Authorize() middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		user := req.Context.Value(config.CtxKeyUser).(*domain.User)
		if requiredRole, ok := operationRoles[req.OperationID]; ok {
			if requiredRole == user.Role {
				return next(req)
			} else {
				return middleware.Response{}, domain.WrapError(domain.ErrAuthorization, fmt.Errorf("operation %s requires role %s", req.OperationID, requiredRole))
			}
		}
		return next(req)
	}
}
