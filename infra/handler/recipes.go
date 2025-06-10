package handler

import (
	"context"
	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/domain"
	"time"
)

type RecipeHandler struct {
	Recipes *domain.RecipeService
}

func NewRecipeHandler(service *domain.RecipeService) *RecipeHandler {
	return &RecipeHandler{
		Recipes: service,
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
	var response []api.ReadMealPlan
	for _, item := range mealplan {
		var recipes []api.ReadRecipe
		for _, recipe := range item.Recipes {
			recipes = append(recipes, api.ReadRecipe{
				ID:   recipe.ID,
				Name: recipe.Name,
			})
		}
		response = append(response, api.ReadMealPlan{
			Date:    item.Date.Format(time.DateOnly),
			Recipes: recipes,
		})
	}
	return response, nil
}

func (h *RecipeHandler) GetRecipes(ctx context.Context) ([]api.ReadRecipe, error) {
	user := ctx.Value(ctxKeyUser).(*domain.User)
	recipes, err := h.Recipes.GetByUser(ctx, user)
	if err != nil {
		return nil, err
	}
	var response []api.ReadRecipe
	for _, recipe := range recipes {
		response = append(response, api.ReadRecipe{
			ID:   recipe.ID,
			Name: recipe.Name,
		})
	}
	return response, nil
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
