package handler

import (
	"context"
	"time"

	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/infra/urlbuilder"
)

type RecipeHandler struct {
	mapper     *responseMapper
	Recipes    *domain.RecipeService
	URLBuilder *urlbuilder.Builder
}

func NewRecipeHandler(service *domain.RecipeService, urlBuilder *urlbuilder.Builder) *RecipeHandler {
	return &RecipeHandler{
		mapper:     newResponseMapper(urlBuilder),
		Recipes:    service,
		URLBuilder: urlBuilder,
	}
}

func (h *RecipeHandler) AddRecipe(ctx context.Context, req *api.WriteRecipeMultipart) (*api.ReadRecipe, error) {
	user := ctx.Value(ctxKeyUser).(*domain.User)
	recipe, err := h.Recipes.Add(ctx, domain.RecipeDetails{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   user,
		Servings:    req.Servings,
		Minutes:     req.Minutes,
	})
	if err != nil {
		return nil, err
	}
	return &api.ReadRecipe{
		ID:   recipe.ID,
		Name: recipe.Name,
	}, nil
}

func (h *RecipeHandler) BrowseRecipes(ctx context.Context) ([]api.ReadRecipe, error) {
	recipes, err := h.Recipes.Browse(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]api.ReadRecipe, len(recipes))
	for i, r := range recipes {
		result[i] = api.ReadRecipe{
			ID:   r.ID,
			Name: r.Name,
		}
	}
	return result, nil
}

func (h *RecipeHandler) DeleteRecipe(ctx context.Context, params api.DeleteRecipeParams) error {
	err := h.Recipes.Delete(ctx, params.RecipeId)
	if err != nil {
		return err
	}
	return nil
}

func (h *RecipeHandler) GetMealPlan(ctx context.Context, params api.GetMealPlanParams) ([]api.ReadMealPlan, error) {
	user := ctx.Value(ctxKeyUser).(*domain.User)
	from := params.From.Or(time.Now())
	until := params.Until.Or(from.Add(7 * 24 * time.Hour))
	mealplan, err := h.Recipes.GetMealPlan(ctx, user, from, until)
	if err != nil {
		return nil, err
	}
	return h.mapper.toReadMealPlanList(mealplan)
}

func (h *RecipeHandler) GetRecipes(ctx context.Context) ([]api.ReadRecipe, error) {
	user := ctx.Value(ctxKeyUser).(*domain.User)
	recipes, err := h.Recipes.GetByUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return h.mapper.toRecipeListResponse(recipes)
}

func (h *RecipeHandler) GetRecipeById(ctx context.Context, params api.GetRecipeByIdParams) (*api.ReadRecipe, error) {
	recipe, err := h.Recipes.GetById(ctx, params.RecipeId)
	if err != nil {
		return nil, domain.ErrRecipeNotFound
	}
	return &api.ReadRecipe{
		ID:   recipe.ID,
		Name: recipe.Name,
	}, nil
}

func (h *RecipeHandler) UpdateRecipe(_ context.Context, _ *api.WriteRecipeMultipart, _ api.UpdateRecipeParams) (*api.ReadRecipe, error) {
	// TODO: Implement
	return &api.ReadRecipe{}, nil
}
