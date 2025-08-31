package permissions

type Slug string

const (
	CreateUnit Slug = "can_create_unit"
	DeleteUnit Slug = "can_delete_unit"
	UpdateUnit Slug = "can_update_unit"
	ViewUnit   Slug = "can_view_unit"
)

const (
	CreateIngredient Slug = "can_create_ingredient"
	DeleteIngredient Slug = "can_delete_ingredient"
	UpdateIngredient Slug = "can_update_ingredient"
	ViewIngredient   Slug = "can_view_ingredient"
)

const (
	CreateMealPlan Slug = "can_create_meal_plan"
	DeleteMealPlan Slug = "can_delete_meal_plan"
	UpdateMealPlan Slug = "can_update_meal_plan"
	ViewMealPlan   Slug = "can_view_meal_plan"
)

const (
	CreateRecipe Slug = "can_create_recipe"
	DeleteRecipe Slug = "can_delete_recipe"
	UpdateRecipe Slug = "can_update_recipe"
	ViewRecipe   Slug = "can_view_recipe"
)

const (
	UpdateProfile Slug = "can_update_profile"
	ViewProfile   Slug = "can_view_profile"
)

const (
	CreateShoppingList Slug = "can_create_shopping_list"
	DeleteShoppingList Slug = "can_delete_shopping_list"
	UpdateShoppingList Slug = "can_update_shopping_list"
	ViewShoppingList   Slug = "can_view_shopping_list"
)
