package handler

import (
	"context"
	"time"

	"github.com/wolfsblu/recipe-manager/api"
	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/config"
	"github.com/wolfsblu/recipe-manager/infra/env"
	"github.com/wolfsblu/recipe-manager/infra/handler/mapper"
)

type RecipeHandler struct {
	mapper  *mapper.APIMapper
	Recipes *domain.RecipeService
}

func NewRecipeHandler(service *domain.RecipeService) *RecipeHandler {
	return &RecipeHandler{
		mapper:  mapper.NewAPIMapper(env.MustGet("BASE_URL")),
		Recipes: service,
	}
}

func (h *RecipeHandler) AddRecipe(ctx context.Context, req *api.WriteRecipe) (*api.ReadRecipe, error) {
	user := ctx.Value(config.CtxKeyUser).(*domain.User)
	recipe := h.mapper.FromWriteRecipe(req)
	recipe.CreatedBy = user

	result, err := h.Recipes.Add(ctx, recipe)
	if err != nil {
		return nil, err
	}

	return h.mapper.ToReadRecipe(result)
}

func (h *RecipeHandler) BrowseRecipes(ctx context.Context) ([]api.ReadRecipe, error) {
	recipes, err := h.Recipes.Browse(ctx)
	if err != nil {
		return nil, err
	}
	return h.mapper.ToRecipes(recipes)
}

func (h *RecipeHandler) DeleteRecipe(ctx context.Context, params api.DeleteRecipeParams) error {
	user := ctx.Value(config.CtxKeyUser).(*domain.User)
	err := h.Recipes.Delete(ctx, user, params.RecipeId)
	if err != nil {
		return err
	}
	return nil
}

func (h *RecipeHandler) GetMealPlan(ctx context.Context, params api.GetMealPlanParams) ([]api.ReadMealPlan, error) {
	user := ctx.Value(config.CtxKeyUser).(*domain.User)
	from := params.From.Or(time.Now())
	until := params.Until.Or(from.Add(7 * 24 * time.Hour))
	mealplan, err := h.Recipes.GetMealPlan(ctx, user, from, until)
	if err != nil {
		return nil, err
	}
	return h.mapper.ToMealPlans(mealplan)
}

func (h *RecipeHandler) GetRecipes(ctx context.Context) ([]api.ReadRecipe, error) {
	user := ctx.Value(config.CtxKeyUser).(*domain.User)
	recipes, err := h.Recipes.GetByUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return h.mapper.ToRecipes(recipes)
}

func (h *RecipeHandler) GetRecipeById(ctx context.Context, params api.GetRecipeByIdParams) (*api.ReadRecipe, error) {
	recipe, err := h.Recipes.GetById(ctx, params.RecipeId)
	if err != nil {
		return nil, domain.ErrRecipeNotFound
	}
	return h.mapper.ToReadRecipe(recipe)
}

func (h *RecipeHandler) UpdateRecipe(_ context.Context, _ *api.WriteRecipe, _ api.UpdateRecipeParams) (*api.ReadRecipe, error) {
	// TODO: Implement
	return &api.ReadRecipe{}, nil
}

func (h *RecipeHandler) GetIngredients(ctx context.Context) ([]api.Ingredient, error) {
	ingredients, err := h.Recipes.GetIngredients(ctx)
	if err != nil {
		return nil, err
	}
	return h.mapper.ToIngredients(ingredients)
}

func (h *RecipeHandler) GetUnits(ctx context.Context) ([]api.ReadUnit, error) {
	units, err := h.Recipes.GetUnits(ctx)
	if err != nil {
		return nil, err
	}
	return h.mapper.ToUnits(units)
}
