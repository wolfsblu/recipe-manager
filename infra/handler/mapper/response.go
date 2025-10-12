package mapper

import (
	"net/url"
	"time"

	"github.com/wolfsblu/recipe-manager/api"
	"github.com/wolfsblu/recipe-manager/domain"
)

func (m *APIMapper) ToIngredientNutrient(nutrient domain.IngredientNutrient) api.IngredientNutrient {
	return api.IngredientNutrient{
		ID:     nutrient.Nutrient.ID,
		Name:   nutrient.Nutrient.Name,
		Unit:   nutrient.Nutrient.Unit,
		Amount: nutrient.Amount,
	}
}

func (m *APIMapper) ToIngredient(ingredient domain.Ingredient) *api.Ingredient {
	nutrients := make([]api.IngredientNutrient, len(ingredient.Nutrients))
	for i, nutrient := range ingredient.Nutrients {
		nutrients[i] = m.ToIngredientNutrient(nutrient)
	}
	return &api.Ingredient{
		ID:        ingredient.ID,
		Name:      ingredient.Name,
		Nutrients: nutrients,
	}
}

func (m *APIMapper) ToReadStepIngredient(ingredient domain.StepIngredient) api.ReadStepIngredient {
	apiIngredient := m.ToIngredient(ingredient.Ingredient)

	return api.ReadStepIngredient{
		ID:        apiIngredient.ID,
		Name:      apiIngredient.Name,
		Nutrients: apiIngredient.Nutrients,
		Unit:      *m.ToUnit(ingredient.Unit),
		Amount:    ingredient.Amount,
	}
}

func (m *APIMapper) ToReadRecipeStep(step domain.RecipeStep) (api.ReadRecipeStep, error) {
	ingredients := make([]api.ReadStepIngredient, len(step.Ingredients))
	for i, ingredient := range step.Ingredients {
		ingredients[i] = m.ToReadStepIngredient(ingredient)
	}

	return api.ReadRecipeStep{
		ID:           step.ID,
		Ingredients:  ingredients,
		Instructions: step.Instructions,
	}, nil
}

func (m *APIMapper) ToIngredients(ingredients []domain.Ingredient) ([]api.Ingredient, error) {
	result := make([]api.Ingredient, len(ingredients))
	for i, ingredient := range ingredients {
		result[i] = *m.ToIngredient(ingredient)
	}
	return result, nil
}

func (m *APIMapper) ToMealPlan(mealPlan domain.MealPlan) (api.ReadMealPlan, error) {
	recipes := make([]api.ReadRecipe, len(mealPlan.Recipes))
	for i, recipe := range mealPlan.Recipes {
		response, err := m.ToReadRecipe(recipe)
		if err != nil {
			return api.ReadMealPlan{}, err
		}
		recipes[i] = *response
	}

	return api.ReadMealPlan{
		Date:    mealPlan.Date.Format(time.DateOnly),
		Recipes: recipes,
	}, nil
}

func (m *APIMapper) ToMealPlans(mealPlans []domain.MealPlan) ([]api.ReadMealPlan, error) {
	result := make([]api.ReadMealPlan, len(mealPlans))
	for i, mealPlan := range mealPlans {
		response, err := m.ToMealPlan(mealPlan)
		if err != nil {
			return nil, err
		}
		result[i] = response
	}
	return result, nil
}

func (m *APIMapper) ToReadRecipe(recipe domain.Recipe) (*api.ReadRecipe, error) {
	images, err := m.ToRecipeImageURLs(recipe.Images)
	if err != nil {
		return nil, err
	}

	tags := m.ToTags(recipe.Tags)
	steps := make([]api.ReadRecipeStep, len(recipe.Steps))
	for i, step := range recipe.Steps {
		readStep, err := m.ToReadRecipeStep(step)
		if err != nil {
			return nil, err
		}
		steps[i] = readStep
	}

	return &api.ReadRecipe{
		ID:          recipe.ID,
		Name:        recipe.Name,
		Description: recipe.Description,
		Servings:    recipe.Servings,
		Minutes:     recipe.Minutes,
		Images:      images,
		Tags:        tags,
		Steps:       steps,
		Votes:       *m.ToRecipeVotes(recipe.Votes),
	}, nil
}

func (m *APIMapper) ToRecipes(recipes []domain.Recipe) ([]api.ReadRecipe, error) {
	result := make([]api.ReadRecipe, len(recipes))
	for i, recipe := range recipes {
		mapped, err := m.ToReadRecipe(recipe)
		if err != nil {
			return nil, err
		}
		result[i] = *mapped
	}
	return result, nil
}

func (m *APIMapper) ToUnit(unit domain.Unit) *api.ReadUnit {
	return &api.ReadUnit{
		ID:     unit.ID,
		Name:   unit.Name,
		Symbol: ToNilString(unit.Symbol),
	}
}

func (m *APIMapper) ToUnits(units []domain.Unit) []api.ReadUnit {
	result := make([]api.ReadUnit, len(units))
	for i, unit := range units {
		result[i] = *m.ToUnit(unit)
	}
	return result
}

func (m *APIMapper) ToTag(tag domain.Tag) api.ReadTag {
	return api.ReadTag{
		ID:   tag.ID,
		Name: tag.Name,
	}
}

func (m *APIMapper) ToTags(tags []domain.Tag) []api.ReadTag {
	result := make([]api.ReadTag, len(tags))
	for i, tag := range tags {
		result[i] = m.ToTag(tag)
	}
	return result
}

func (m *APIMapper) ToRecipeVotes(votes domain.RecipeVotes) *api.RecipeVotes {
	return &api.RecipeVotes{
		Total: votes.Total,
		User:  api.RecipeVotesUser(votes.User),
	}
}

func (m *APIMapper) ToRecipeImageURLs(images []domain.RecipeImage) ([]url.URL, error) {
	urls := make([]url.URL, len(images))
	for i, image := range images {
		urls[i] = *image.URL
	}
	return urls, nil
}

func ToNilString(s *string) api.NilString {
	if s != nil {
		return api.NewNilString(*s)
	}
	return api.NilString{
		Value: "",
		Null:  true,
	}
}

func (m *APIMapper) ToShoppingListItem(item domain.ShoppingListItem) (*api.ReadShoppingListItem, error) {
	return &api.ReadShoppingListItem{
		ID:         item.ID,
		Ingredient: item.Ingredient,
		Quantity:   ToNilString(item.Quantity),
		Unit:       ToNilString(item.Unit),
		Done:       item.Done,
		SortOrder:  item.SortOrder,
	}, nil
}

func (m *APIMapper) ToShoppingListItems(items []domain.ShoppingListItem) ([]api.ReadShoppingListItem, error) {
	result := make([]api.ReadShoppingListItem, len(items))
	for i, item := range items {
		mapped, err := m.ToShoppingListItem(item)
		if err != nil {
			return nil, err
		}
		result[i] = *mapped
	}
	return result, nil
}

func (m *APIMapper) ToShoppingList(list domain.ShoppingList) (*api.ReadShoppingList, error) {
	items, err := m.ToShoppingListItems(list.Items)
	if err != nil {
		return nil, err
	}
	return &api.ReadShoppingList{
		ID:    list.ID,
		Name:  list.Name,
		Items: items,
	}, nil
}

func (m *APIMapper) ToShoppingLists(lists []domain.ShoppingList) ([]api.ReadShoppingList, error) {
	result := make([]api.ReadShoppingList, len(lists))
	for i, list := range lists {
		mapped, err := m.ToShoppingList(list)
		if err != nil {
			return nil, err
		}
		result[i] = *mapped
	}
	return result, nil
}
