package domain

func NewRecipeService(notifier NotificationSender, store RecipeStore) *RecipeService {
	return &RecipeService{
		store:  store,
		sender: notifier,
	}
}

func NewUserService(notifier NotificationSender, store UserStore) *UserService {
	return &UserService{
		store:  store,
		sender: notifier,
	}
}

func NewShoppingService(store ShoppingListStore) *ShoppingService {
	return &ShoppingService{
		store: store,
	}
}
