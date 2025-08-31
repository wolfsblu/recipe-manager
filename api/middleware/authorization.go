package middleware

import (
	"fmt"

	"github.com/ogen-go/ogen/middleware"
	"github.com/wolfsblu/go-chef/api/operations"
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/domain/permissions"
	"github.com/wolfsblu/go-chef/infra/config"
)

var operationPermissions = map[operations.ID]permissions.Slug{
	operations.AddRecipe:     permissions.CreateRecipe,
	operations.DeleteRecipe:  permissions.DeleteRecipe,
	operations.UpdateRecipe:  permissions.UpdateRecipe,
	operations.GetRecipes:    permissions.ViewRecipe,
	operations.GetRecipeById: permissions.ViewRecipe,
}

func Authorize() middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		user, ok := req.Context.Value(config.CtxKeyUser).(*domain.User)
		if !ok {
			return next(req)
		}

		if requiredPermission, ok := operationPermissions[operations.ID(req.OperationID)]; ok {
			for _, permission := range user.Role.Permissions {
				if permission.Slug == requiredPermission {
					return next(req)
				}
			}
			return middleware.Response{}, domain.WrapError(domain.ErrAuthorization, fmt.Errorf("operation %s requires role %s", req.OperationID, requiredPermission))
		}

		return next(req)
	}
}
