package permissions

type Slug string

// Units
const (
	CreateUnit Slug = "can_create_unit"
	DeleteUnit Slug = "can_delete_unit"
	ListUnits  Slug = "can_list_units"
	UpdateUnit Slug = "can_update_unit"
	ViewUnit   Slug = "can_view_unit"
)

// Ingredients
const (
	CreateIngredient Slug = "can_create_ingredient"
	DeleteIngredient Slug = "can_delete_ingredient"
	ListIngredients  Slug = "can_list_ingredients"
	UpdateIngredient Slug = "can_update_ingredient"
	ViewIngredient   Slug = "can_view_ingredient"
)

// Meal Plans
const (
	CreateMealPlan Slug = "can_create_meal_plan"
	DeleteMealPlan Slug = "can_delete_meal_plan"
	ListMealPlans  Slug = "can_list_meal_plans"
	UpdateMealPlan Slug = "can_update_meal_plan"
	ViewMealPlan   Slug = "can_view_meal_plan"
)

// Recipes
const (
	CreateRecipe Slug = "can_create_recipe"
	DeleteRecipe Slug = "can_delete_recipe"
	ListRecipes  Slug = "can_list_recipes"
	UpdateRecipe Slug = "can_update_recipe"
	ViewRecipe   Slug = "can_view_recipe"
)

// Profile
// Note: 'ListProfiles' would typically be an admin-only permission to see all users.
const (
	UpdateProfile Slug = "can_update_profile"
	ViewProfile   Slug = "can_view_profile"
)

// Shopping Lists
const (
	CreateShoppingList Slug = "can_create_shopping_list"
	DeleteShoppingList Slug = "can_delete_shopping_list"
	ListShoppingLists  Slug = "can_list_shopping_lists"
	UpdateShoppingList Slug = "can_update_shopping_list"
	ViewShoppingList   Slug = "can_view_shopping_list"
)
