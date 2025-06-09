package sqlite

import (
	"context"
	"github.com/wolfsblu/go-chef/domain"
	_ "modernc.org/sqlite"
	"time"
)

func (s *Store) BrowseRecipes(ctx context.Context) (recipes []domain.Recipe, err error) {
	result, err := s.query().BrowseRecipes(ctx)
	if err != nil {
		return nil, err
	}
	for _, recipe := range result {
		recipes = append(recipes, recipe.AsDomainModel())
	}
	return recipes, err
}

func (s *Store) CreateRecipe(ctx context.Context, r domain.RecipeDetails) (recipe domain.Recipe, _ error) {
	payload := CreateRecipeParams{
		Name:        r.Name,
		Servings:    0,
		Minutes:     0,
		Description: "",
		CreatedBy:   r.CreatedBy.ID,
	}

	result, err := s.query().CreateRecipe(ctx, payload)
	if err != nil {
		return recipe, err
	}

	return result.AsDomainModel(), nil
}

func (s *Store) DeleteRecipe(ctx context.Context, id int64) error {
	return s.query().DeleteRecipe(ctx, id)
}

func (s *Store) GetMealPlan(ctx context.Context, user *domain.User, from time.Time, until time.Time) ([]domain.MealPlan, error) {
	result, err := s.query().GetMealPlan(ctx, GetMealPlanParams{
		UserID:    user.ID,
		FromDate:  from.Format(time.DateOnly),
		UntilDate: until.Format(time.DateOnly),
	})
	if err != nil {
		return []domain.MealPlan{}, err
	}

	grouped := make(map[string][]domain.Recipe)
	for _, item := range result {
		grouped[item.MealPlan.Date] = append(grouped[item.MealPlan.Date], item.Recipe.AsDomainModel())
	}

	mealplan := make([]domain.MealPlan, len(grouped))
	for key, recipes := range grouped {
		date, err := time.Parse(time.DateOnly, key)
		if err != nil {
			return []domain.MealPlan{}, err
		}
		mealplan = append(mealplan, domain.MealPlan{
			Date:    date,
			Recipes: recipes,
		})
	}
	return mealplan, nil
}

func (s *Store) GetRecipeById(ctx context.Context, id int64) (recipe domain.Recipe, _ error) {
	result, err := s.query().GetRecipe(ctx, id)
	if err != nil {
		return recipe, err
	}
	return result.AsDomainModel(), nil
}

func (s *Store) GetRecipesByUser(ctx context.Context, user *domain.User) (recipes []domain.Recipe, _ error) {
	result, err := s.query().ListRecipes(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	for _, item := range result {
		recipes = append(recipes, item.AsDomainModel())
	}
	return recipes, nil
}
