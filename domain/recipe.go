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

type Ingredient struct {
	ID   int64
	Name string
}

type Unit struct {
	ID   int64
	Name string
	Code *string
}

type Tag struct {
	ID   int64
	Name string
}
