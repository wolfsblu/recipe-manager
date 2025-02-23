package handlers

import (
	"github.com/google/wire"
	"github.com/wolfsblu/go-chef/api"
)

var Set = wire.NewSet(
	NewRecipeHandler,
	NewSecurityHandler,
	wire.Bind(new(api.SecurityHandler), new(*SecurityHandler)),
	wire.Bind(new(api.Handler), new(*RecipeHandler)),
)
