package domain

import "log"

func NewRecipeService(notifier NotificationSender, store RecipeStore) *RecipeService {
	log.Println("creating new recipe service")
	return &RecipeService{
		store:  store,
		sender: notifier,
	}
}
