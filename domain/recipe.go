package domain

import (
	"net/url"
	"time"
)

type RecipeDetails struct {
	Name        string
	Description string
	CreatedBy   *User
	Servings    int32
	Minutes     int32
}

type StepIngredient struct {
	Unit       Unit
	Amount     float64
	Ingredient Ingredient
}

type RecipeStep struct {
	ID           int32
	Instructions string
	Ingredients  []StepIngredient
}

type Recipe struct {
	ID        int32
	CreatedAt time.Time
	Tags      []Tag
	Images    []RecipeImage
	Steps     []RecipeStep
	RecipeDetails
}

type RecipeImage struct {
	ID  int32
	URL *url.URL
}

type MealPlan struct {
	Date    time.Time
	Entries []MealPlanRecipe
}

type MealPlanRecipe struct {
	ID     int32
	Recipe Recipe
}

type MealPlanEntry struct {
	UserID    int32
	RecipeID  int32
	Date      time.Time
	SortOrder int32
}

type Ingredient struct {
	ID        int32
	Name      string
	Nutrients []IngredientNutrient
}

type Nutrient struct {
	ID   int32
	Name string
	Unit string
}

type IngredientNutrient struct {
	Nutrient Nutrient
	Amount   float32
}

type Unit struct {
	ID     int32
	Name   string
	Symbol *string
}

type Tag struct {
	ID   int32
	Name string
}
