package handler

import (
	"context"
	"time"

	"github.com/ogen-go/ogen/json"
	"github.com/wolfsblu/recipe-manager/api"
	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/domain/pagination"
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
	user, ok := ctx.Value(config.CtxKeyUser).(*domain.User)
	if !ok || user == nil {
		return nil, domain.ErrAuthentication
	}
	recipe := h.mapper.FromWriteRecipe(req)
	recipe.CreatedBy = user

	result, err := h.Recipes.Add(ctx, recipe)
	if err != nil {
		return nil, err
	}

	return h.mapper.ToReadRecipe(result)
}

func (h *RecipeHandler) DeleteRecipe(ctx context.Context, params api.DeleteRecipeParams) error {
	user, ok := ctx.Value(config.CtxKeyUser).(*domain.User)
	if !ok || user == nil {
		return domain.ErrAuthentication
	}
	err := h.Recipes.Delete(ctx, user, params.RecipeId)
	if err != nil {
		return err
	}
	return nil
}

func (h *RecipeHandler) GetMealPlan(ctx context.Context, params api.GetMealPlanParams) (*api.PaginatedMealPlan, error) {
	user, ok := ctx.Value(config.CtxKeyUser).(*domain.User)
	if !ok || user == nil {
		return nil, domain.ErrAuthentication
	}

	from := params.From.Or(time.Now())
	until := params.Until.Or(from.Add(7 * 24 * time.Hour))

	paginationReq, err := pagination.ValidatePage(params.Cursor.Value, int(params.Limit.Or(pagination.DefaultLimit)))
	if err != nil {
		return nil, err
	}

	result, err := h.Recipes.GetMealPlan(ctx, user, from, until, paginationReq)
	if err != nil {
		return nil, err
	}

	mealplans, err := h.mapper.ToMealPlans(result.Data)
	if err != nil {
		return nil, err
	}

	var nextCursor api.OptNilString
	if result.NextCursor != nil {
		nextCursor = api.NewOptNilString(*result.NextCursor)
	}

	return &api.PaginatedMealPlan{
		Data:       mealplans,
		NextCursor: nextCursor,
		HasMore:    result.HasMore,
	}, nil
}

func (h *RecipeHandler) CreateMealPlan(ctx context.Context, req *api.WriteMealPlan) error {
	user, ok := ctx.Value(config.CtxKeyUser).(*domain.User)
	if !ok || user == nil {
		return domain.ErrAuthentication
	}
	return h.Recipes.CreateMealPlan(ctx, user, req.RecipeId, req.Date)
}

func (h *RecipeHandler) DeleteMealPlan(ctx context.Context, params api.DeleteMealPlanParams) error {
	user, ok := ctx.Value(config.CtxKeyUser).(*domain.User)
	if !ok || user == nil {
		return domain.ErrAuthentication
	}
	return h.Recipes.DeleteMealPlan(ctx, user, params.RecipeId, params.Date)
}

func (h *RecipeHandler) GetRecipes(ctx context.Context, params api.GetRecipesParams) (*api.PaginatedRecipes, error) {
	user, ok := ctx.Value(config.CtxKeyUser).(*domain.User)
	if !ok || user == nil {
		return nil, domain.ErrAuthentication
	}

	paginationReq, err := pagination.ValidatePage(params.Cursor.Value, int(params.Limit.Or(pagination.DefaultLimit)))
	if err != nil {
		return nil, err
	}

	result, err := h.Recipes.GetByUser(ctx, user, paginationReq)
	if err != nil {
		return nil, err
	}

	recipes, err := h.mapper.ToRecipes(result.Data)
	if err != nil {
		return nil, err
	}

	var nextCursor api.OptNilString
	if result.NextCursor != nil {
		nextCursor = api.NewOptNilString(*result.NextCursor)
	}

	return &api.PaginatedRecipes{
		Data:       recipes,
		NextCursor: nextCursor,
		HasMore:    result.HasMore,
	}, nil
}

func (h *RecipeHandler) GetRecipeById(ctx context.Context, params api.GetRecipeByIdParams) (*api.ReadRecipe, error) {
	user, ok := ctx.Value(config.CtxKeyUser).(*domain.User)
	if !ok || user == nil {
		return nil, domain.ErrAuthentication
	}
	recipe, err := h.Recipes.GetById(ctx, user, params.RecipeId)
	if err != nil {
		return nil, domain.ErrRecipeNotFound
	}
	return h.mapper.ToReadRecipe(recipe)
}

func (h *RecipeHandler) UpdateRecipe(ctx context.Context, req *api.WriteRecipe, params api.UpdateRecipeParams) (*api.ReadRecipe, error) {
	user, ok := ctx.Value(config.CtxKeyUser).(*domain.User)
	if !ok || user == nil {
		return nil, domain.ErrAuthentication
	}
	recipe := h.mapper.FromWriteRecipe(req)
	recipe.ID = params.RecipeId
	recipe.CreatedBy = user

	result, err := h.Recipes.UpdateRecipe(ctx, recipe)
	if err != nil {
		return nil, err
	}

	return h.mapper.ToReadRecipe(result)
}

func (h *RecipeHandler) GetIngredients(ctx context.Context, params api.GetIngredientsParams) (*api.PaginatedIngredients, error) {
	paginationReq, err := pagination.ValidatePage(params.Cursor.Value, int(params.Limit.Or(pagination.DefaultLimit)))
	if err != nil {
		return nil, err
	}

	result, err := h.Recipes.GetIngredients(ctx, paginationReq)
	if err != nil {
		return nil, err
	}

	ingredients, err := h.mapper.ToIngredients(result.Data)
	if err != nil {
		return nil, err
	}

	var nextCursor api.OptNilString
	if result.NextCursor != nil {
		nextCursor = api.NewOptNilString(*result.NextCursor)
	}

	return &api.PaginatedIngredients{
		Data:       ingredients,
		NextCursor: nextCursor,
		HasMore:    result.HasMore,
	}, nil
}

func (h *RecipeHandler) GetUnits(ctx context.Context, params api.GetUnitsParams) (*api.PaginatedUnits, error) {
	paginationReq, err := pagination.ValidatePage(params.Cursor.Value, int(params.Limit.Or(pagination.DefaultLimit)))
	if err != nil {
		return nil, err
	}

	result, err := h.Recipes.GetUnits(ctx, paginationReq)
	if err != nil {
		return nil, err
	}

	var nextCursor api.OptNilString
	if result.NextCursor != nil {
		nextCursor = api.NewOptNilString(*result.NextCursor)
	}

	return &api.PaginatedUnits{
		Data:       h.mapper.ToUnits(result.Data),
		NextCursor: nextCursor,
		HasMore:    result.HasMore,
	}, nil
}

func (h *RecipeHandler) GetTags(ctx context.Context, params api.GetTagsParams) (*api.PaginatedTags, error) {
	paginationReq, err := pagination.ValidatePage(params.Cursor.Value, int(params.Limit.Or(pagination.DefaultLimit)))
	if err != nil {
		return nil, err
	}

	result, err := h.Recipes.GetTags(ctx, paginationReq)
	if err != nil {
		return nil, err
	}

	var nextCursor api.OptNilString
	if result.NextCursor != nil {
		nextCursor = api.NewOptNilString(*result.NextCursor)
	}

	return &api.PaginatedTags{
		Data:       h.mapper.ToTags(result.Data),
		NextCursor: nextCursor,
		HasMore:    result.HasMore,
	}, nil
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
	user, ok := ctx.Value(config.CtxKeyUser).(*domain.User)
	if !ok || user == nil {
		return nil, domain.ErrAuthentication
	}

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
