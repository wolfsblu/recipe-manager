package operations

type ID string

const (
	// Recipes
	GetRecipes    ID = "getRecipes"
	AddRecipe     ID = "addRecipe"
	GetRecipeById ID = "getRecipeById"
	UpdateRecipe  ID = "updateRecipe"
	DeleteRecipe  ID = "deleteRecipe"

	// User
	Login          ID = "login"
	Logout         ID = "logout"
	Register       ID = "register"
	ConfirmUser    ID = "confirmUser"
	UpdatePassword ID = "updatePassword"
	ResetPassword  ID = "resetPassword"
	GetUserProfile ID = "getUserProfile"

	// Meal Plan
	GetMealPlan ID = "getMealPlan"

	// Ingredients
	GetIngredients ID = "getIngredients"

	// Units
	GetUnits ID = "getUnits"
)
