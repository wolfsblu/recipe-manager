package domain

import (
	"net/url"
	"time"
)

type RecipeDetails struct {
	Name        string
	Description string
	CreatedBy   *User
	Servings    int64
	Minutes     int64
}

type StepIngredient struct {
	Unit       Unit
	Amount     float64
	Ingredient Ingredient
}

type RecipeStep struct {
	ID           int64
	Instructions string
	Ingredients  []StepIngredient
}

type Recipe struct {
	ID     int64
	Tags   []Tag
	Images []RecipeImage
	Steps  []RecipeStep
	Votes  RecipeVotes
	RecipeDetails
}

type RecipeImage struct {
	ID  int64
	URL *url.URL
}

type MealPlan struct {
	Date    time.Time
	Recipes []Recipe
}

type MealPlanEntry struct {
	UserID    int64
	RecipeID  int64
	Date      time.Time
	SortOrder int64
}

type Ingredient struct {
	ID        int64
	Name      string
	Nutrients []IngredientNutrient
}

type Nutrient struct {
	ID   int64
	Name string
	Unit string
}

type IngredientNutrient struct {
	Nutrient Nutrient
	Amount   float64
}

type Unit struct {
	ID     int64
	Name   string
	Symbol *string
}

type Tag struct {
	ID   int64
	Name string
}

type RecipeVotes struct {
	Total int64
	User  int64
}
