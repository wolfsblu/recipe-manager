package handler

import (
	"context"
	"time"

	"github.com/ogen-go/ogen/json"
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
	user := ctx.Value(config.CtxKeyUser).(*domain.User)
	recipe, err := h.Recipes.GetById(ctx, user, params.RecipeId)
	if err != nil {
		return nil, domain.ErrRecipeNotFound
	}
	return h.mapper.ToReadRecipe(recipe)
}

func (h *RecipeHandler) UpdateRecipe(ctx context.Context, req *api.WriteRecipe, params api.UpdateRecipeParams) (*api.ReadRecipe, error) {
	user := ctx.Value(config.CtxKeyUser).(*domain.User)
	recipe := h.mapper.FromWriteRecipe(req)
	recipe.ID = params.RecipeId
	recipe.CreatedBy = user

	result, err := h.Recipes.UpdateRecipe(ctx, recipe)
	if err != nil {
		return nil, err
	}

	return h.mapper.ToReadRecipe(result)
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
	return h.mapper.ToUnits(units), nil
}

func (h *RecipeHandler) GetTags(ctx context.Context) ([]api.ReadTag, error) {
	tags, err := h.Recipes.GetTags(ctx)
	if err != nil {
		return nil, err
	}
	return h.mapper.ToTags(tags), nil
}

func (h *RecipeHandler) AddVote(ctx context.Context, req *api.Vote, params api.AddVoteParams) (*api.RecipeVotes, error) {
	user := ctx.Value(config.CtxKeyUser).(*domain.User)

	votes, err := h.Recipes.AddVote(ctx, user, params.RecipeId, int64(req.Vote))
	if err != nil {
		return nil, err
	}

	return h.mapper.ToRecipeVotes(votes), nil
}

func (h *RecipeHandler) RemoveVote(ctx context.Context, params api.RemoveVoteParams) (*api.RecipeVotes, error) {
	user := ctx.Value(config.CtxKeyUser).(*domain.User)

	votes, err := h.Recipes.RemoveVote(ctx, user, params.RecipeId)
	if err != nil {
		return nil, err
	}

	return h.mapper.ToRecipeVotes(votes), nil
}

func (h *RecipeHandler) AddIngredient(ctx context.Context, req *api.WriteIngredient) (*api.Ingredient, error) {
	ingredient := h.mapper.FromWriteIngredient(req)

	result, err := h.Recipes.AddIngredient(ctx, ingredient)
	if err != nil {
		return nil, err
	}

	return h.mapper.ToIngredient(result), nil
}

func (h *RecipeHandler) UpdateIngredient(ctx context.Context, req *api.WriteIngredient, params api.UpdateIngredientParams) (*api.Ingredient, error) {
	ingredient := h.mapper.FromWriteIngredient(req)
	ingredient.ID = params.IngredientId

	result, err := h.Recipes.UpdateIngredient(ctx, ingredient)
	if err != nil {
		return nil, err
	}

	return h.mapper.ToIngredient(result), nil
}

func (h *RecipeHandler) DeleteIngredient(ctx context.Context, params api.DeleteIngredientParams) error {
	return h.Recipes.DeleteIngredient(ctx, params.IngredientId)
}

func (h *RecipeHandler) AddUnit(ctx context.Context, req *api.WriteUnit) (*api.ReadUnit, error) {
	unit := h.mapper.FromWriteUnit(req)

	result, err := h.Recipes.AddUnit(ctx, unit)
	if err != nil {
		return nil, err
	}

	return h.mapper.ToUnit(result), nil
}

func (h *RecipeHandler) UpdateUnit(ctx context.Context, req *api.WriteUnit, params api.UpdateUnitParams) (*api.ReadUnit, error) {
	unit := h.mapper.FromWriteUnit(req)
	unit.ID = params.UnitId

	result, err := h.Recipes.UpdateUnit(ctx, unit)
	if err != nil {
		return nil, err
	}

	return h.mapper.ToUnit(result), nil
}

func (h *RecipeHandler) DeleteUnit(ctx context.Context, params api.DeleteUnitParams) error {
	return h.Recipes.DeleteUnit(ctx, params.UnitId)
}

func (h *RecipeHandler) PatchRecipe(ctx context.Context, req *api.RecipeWriteFields, params api.PatchRecipeParams) (*api.ReadRecipe, error) {
	user := ctx.Value(config.CtxKeyUser).(*domain.User)

	existingRecipe, err := h.Recipes.GetById(ctx, user, params.RecipeId)
	if err != nil {
		return nil, err
	}

	writeRecipe, err := h.mapper.ToWriteRecipe(existingRecipe)
	if err != nil {
		return nil, err
	}

	var mergedWriteRecipe api.WriteRecipe
	patchedRecipe, err := h.mergePatch(writeRecipe, req)
	if err := json.Unmarshal(patchedRecipe, &mergedWriteRecipe); err != nil {
		return nil, err
	}

	recipe := h.mapper.FromWriteRecipe(&mergedWriteRecipe)
	recipe.ID = params.RecipeId
	recipe.CreatedBy = user

	result, err := h.Recipes.UpdateRecipe(ctx, recipe)
	if err != nil {
		return nil, err
	}

	return h.mapper.ToReadRecipe(result)
}
