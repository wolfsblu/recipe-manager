package sqlite

import (
	"context"
	"net/url"
	"sort"
	"time"

	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/database"
	_ "modernc.org/sqlite"
)

type recipeRelations struct {
	tags                    []database.GetTagsForRecipesRow
	images                  []database.GetImagesForRecipesRow
	steps                   []database.GetStepsForRecipesRow
	ingredients             []database.GetIngredientsForRecipesRow
	votes                   map[int64]domain.RecipeVotes
	populatedIngredientsMap map[int64]domain.Ingredient
}

func (s *Store) BrowseRecipes(ctx context.Context) (recipes []domain.Recipe, err error) {
	result, err := s.query().BrowseRecipes(ctx)
	if err != nil {
		return nil, err
	}
	for _, recipe := range result {
		recipes = append(recipes, s.mapper.ToRecipe(recipe))
	}
	return recipes, err
}

func (s *Store) CreateRecipe(ctx context.Context, recipe domain.Recipe) (domain.Recipe, error) {
	var recipeId int64
	err := s.WithTransaction(ctx, func(tx *TxStore) error {
		var err error
		recipeId, err = tx.query().CreateRecipe(ctx, tx.mapper.FromRecipe(recipe))
		if err != nil {
			return err
		}
		if err = tx.createRecipeSteps(ctx, recipeId, recipe.Steps); err != nil {
			return err
		}
		if err = tx.createRecipeImages(ctx, recipeId, recipe.Images); err != nil {
			return err
		}
		if err = tx.createRecipeTags(ctx, recipeId, recipe.Tags); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return domain.Recipe{}, err
	}
	return s.GetRecipeById(ctx, recipe.CreatedBy, recipeId)
}

func (s *Store) createRecipeSteps(ctx context.Context, recipeID int64, steps []domain.RecipeStep) error {
	for i, step := range steps {
		stepId, err := s.query().CreateRecipeStep(ctx, s.mapper.FromRecipeStep(recipeID, step, int64(i)))
		if err != nil {
			return err
		}

		if err = s.createStepIngredients(ctx, stepId, step.Ingredients); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) createStepIngredients(ctx context.Context, stepID int64, ingredients []domain.StepIngredient) error {
	for i, ingredient := range ingredients {
		_, err := s.query().CreateStepIngredient(ctx, s.mapper.FromStepIngredient(stepID, ingredient, int64(i)))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) createRecipeImages(ctx context.Context, recipeID int64, images []domain.RecipeImage) error {
	for _, image := range images {
		_, err := s.query().CreateRecipeImages(ctx, s.mapper.FromRecipeImage(recipeID, image))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) createRecipeTags(ctx context.Context, recipeID int64, tags []domain.Tag) error {
	for _, tag := range tags {
		err := s.query().CreateRecipeTag(ctx, s.mapper.FromRecipeTag(recipeID, tag))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) DeleteRecipe(ctx context.Context, id int64) error {
	return s.query().DeleteRecipe(ctx, id)
}

func (s *Store) CreateMealPlan(ctx context.Context, entry domain.MealPlanEntry) error {
	return s.query().CreateMealPlan(ctx, s.mapper.FromMealPlanEntry(entry))
}

func (s *Store) DeleteMealPlan(ctx context.Context, userID int64, recipeID int64, date time.Time) error {
	return s.query().DeleteMealPlan(ctx, database.DeleteMealPlanParams{
		UserID:   userID,
		RecipeID: recipeID,
		Date:     date.Format(time.DateOnly),
	})
}

func (s *Store) GetMealPlan(ctx context.Context, user *domain.User, from time.Time, until time.Time) ([]domain.MealPlan, error) {
	result, err := s.query().GetMealPlan(ctx, database.GetMealPlanParams{
		UserID:    user.ID,
		FromDate:  from.Format(time.DateOnly),
		UntilDate: until.Format(time.DateOnly),
	})
	if err != nil {
		return []domain.MealPlan{}, err
	}

	recipeMap := make(map[int64]domain.Recipe)
	for _, item := range result {
		recipe := s.mapper.ToRecipe(item.Recipe)
		recipeMap[recipe.ID] = recipe
	}

	recipes := make([]domain.Recipe, 0, len(recipeMap))
	for _, recipe := range recipeMap {
		recipes = append(recipes, recipe)
	}

	populatedRecipes, err := s.populateRecipeRelations(ctx, user, recipes)
	if err != nil {
		return []domain.MealPlan{}, err
	}

	populatedRecipeMap := make(map[int64]domain.Recipe)
	for _, recipe := range populatedRecipes {
		populatedRecipeMap[recipe.ID] = recipe
	}

	grouped := make(map[string][]domain.Recipe)
	for _, item := range result {
		recipe := populatedRecipeMap[item.Recipe.ID]
		grouped[item.MealPlan.Date] = append(grouped[item.MealPlan.Date], recipe)
	}

	i := 0
	mealPlan := make([]domain.MealPlan, len(grouped))
	for key, recipes := range grouped {
		date, err := time.Parse(time.DateOnly, key)
		if err != nil {
			return []domain.MealPlan{}, err
		}
		mealPlan[i] = domain.MealPlan{
			Date:    date,
			Recipes: recipes,
		}
		i++
	}

	sort.Slice(mealPlan, func(i, j int) bool {
		return mealPlan[i].Date.Before(mealPlan[j].Date)
	})
	return mealPlan, nil
}

func (s *Store) GetIngredients(ctx context.Context) ([]domain.Ingredient, error) {
	result, err := s.query().GetIngredients(ctx)
	if err != nil {
		return nil, err
	}

	ingredients := make([]domain.Ingredient, len(result))
	for i, ingredient := range result {
		ingredients[i] = s.mapper.ToIngredient(ingredient)
	}

	return s.populateIngredientNutrients(ctx, ingredients)
}

func (s *Store) populateIngredientNutrients(ctx context.Context, ingredients []domain.Ingredient) ([]domain.Ingredient, error) {
	if len(ingredients) == 0 {
		return ingredients, nil
	}

	ingredientIDs := make([]int64, len(ingredients))
	for i, ing := range ingredients {
		ingredientIDs[i] = ing.ID
	}

	nutrients, err := s.query().GetNutrientsForIngredients(ctx, ingredientIDs)
	if err != nil {
		return nil, err
	}

	// Group nutrients by ingredient ID
	nutrientsByIngredient := make(map[int64][]domain.IngredientNutrient)
	for _, nutrient := range nutrients {
		nutrientsByIngredient[nutrient.IngredientID] = append(
			nutrientsByIngredient[nutrient.IngredientID],
			domain.IngredientNutrient{
				Nutrient: s.mapper.ToNutrient(nutrient.Nutrient),
				Amount:   nutrient.Amount,
			},
		)
	}

	// Assign nutrients to ingredients
	populatedIngredients := make([]domain.Ingredient, len(ingredients))
	for i, ingredient := range ingredients {
		ingredient.Nutrients = nutrientsByIngredient[ingredient.ID]
		if ingredient.Nutrients == nil {
			ingredient.Nutrients = []domain.IngredientNutrient{}
		}
		populatedIngredients[i] = ingredient
	}

	return populatedIngredients, nil
}

func (s *Store) GetUnits(ctx context.Context) ([]domain.Unit, error) {
	result, err := s.query().GetUnits(ctx)
	if err != nil {
		return nil, err
	}

	units := make([]domain.Unit, len(result))
	for i, unit := range result {
		units[i] = s.mapper.ToUnit(unit)
	}

	return units, nil
}

func (s *Store) GetTags(ctx context.Context) ([]domain.Tag, error) {
	result, err := s.query().GetTags(ctx)
	if err != nil {
		return nil, err
	}

	tags := make([]domain.Tag, len(result))
	for i, tag := range result {
		tags[i] = s.mapper.ToTag(tag)
	}

	return tags, nil
}

func (s *Store) GetRecipeById(ctx context.Context, user *domain.User, id int64) (recipe domain.Recipe, _ error) {
	result, err := s.query().GetRecipe(ctx, id)
	if err != nil {
		return recipe, err
	}

	recipe = s.mapper.ToRecipe(result)
	populatedRecipes, err := s.populateRecipeRelations(ctx, user, []domain.Recipe{recipe})
	if err != nil {
		return recipe, err
	}

	return populatedRecipes[0], nil
}

func (s *Store) GetRecipesByUser(ctx context.Context, user *domain.User) (recipes []domain.Recipe, _ error) {
	result, err := s.query().ListRecipes(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	recipes = make([]domain.Recipe, len(result))
	for i, recipe := range result {
		recipes[i] = s.mapper.ToRecipe(recipe)
	}

	return s.populateRecipeRelations(ctx, user, recipes)
}

func (s *Store) populateRecipeRelations(ctx context.Context, user *domain.User, recipes []domain.Recipe) ([]domain.Recipe, error) {
	if len(recipes) == 0 {
		return recipes, nil
	}

	recipeIds := make([]int64, len(recipes))
	for i, recipe := range recipes {
		recipeIds[i] = recipe.ID
	}

	relations, err := s.getRecipeRelations(ctx, user, recipeIds)
	if err != nil {
		return nil, err
	}

	tagsByRecipe := s.groupTagsByRecipe(relations.tags)
	ingredientsByStep := s.groupIngredientsByStep(relations.ingredients, relations.populatedIngredientsMap)
	stepsByRecipe := s.groupStepsByRecipe(relations.steps, ingredientsByStep)
	imagesByRecipe, err := s.groupImagesByRecipe(relations.images)
	if err != nil {
		return nil, err
	}

	populatedRecipes := make([]domain.Recipe, len(recipes))
	for i, recipe := range recipes {
		recipe.Tags = tagsByRecipe[recipe.ID]
		recipe.Images = imagesByRecipe[recipe.ID]
		recipe.Votes = relations.votes[recipe.ID]
		recipe.Steps = stepsByRecipe[recipe.ID]
		populatedRecipes[i] = recipe
	}
	return populatedRecipes, nil
}

func (s *Store) getRecipeRelations(ctx context.Context, user *domain.User, recipeIds []int64) (*recipeRelations, error) {
	tags, err := s.query().GetTagsForRecipes(ctx, recipeIds)
	if err != nil {
		return nil, err
	}

	images, err := s.query().GetImagesForRecipes(ctx, recipeIds)
	if err != nil {
		return nil, err
	}

	votes, err := s.getVotes(ctx, user, recipeIds)
	if err != nil {
		return nil, err
	}

	steps, err := s.query().GetStepsForRecipes(ctx, recipeIds)
	if err != nil {
		return nil, err
	}

	ingredients, err := s.query().GetIngredientsForRecipes(ctx, recipeIds)
	if err != nil {
		return nil, err
	}

	// Extract unique ingredients and populate with nutrients
	uniqueIngredients := make([]domain.Ingredient, 0)
	seenIngredients := make(map[int64]bool)
	for _, ing := range ingredients {
		if !seenIngredients[ing.IngredientID] {
			uniqueIngredients = append(uniqueIngredients, s.mapper.ToIngredientFromRecipeRow(ing))
			seenIngredients[ing.IngredientID] = true
		}
	}

	populatedIngredients, err := s.populateIngredientNutrients(ctx, uniqueIngredients)
	if err != nil {
		return nil, err
	}

	// Store populated ingredients in a map for quick lookup
	populatedIngredientsMap := make(map[int64]domain.Ingredient)
	for _, ingredient := range populatedIngredients {
		populatedIngredientsMap[ingredient.ID] = ingredient
	}

	return &recipeRelations{
		tags:                    tags,
		images:                  images,
		steps:                   steps,
		ingredients:             ingredients,
		votes:                   votes,
		populatedIngredientsMap: populatedIngredientsMap,
	}, nil
}

func (s *Store) getVotes(ctx context.Context, user *domain.User, recipeIds []int64) (map[int64]domain.RecipeVotes, error) {
	voteStats, err := s.query().GetVotesForRecipes(ctx, recipeIds)
	if err != nil {
		return nil, err
	}
	userVotes, err := s.query().GetUserVotesForRecipes(ctx, s.mapper.FromUserVotesParams(recipeIds, user.ID))
	if err != nil {
		return nil, err
	}

	votesByRecipe := make(map[int64]int64)
	for _, vote := range voteStats {
		if vote.TotalScore != nil {
			votesByRecipe[vote.RecipeID] = vote.TotalScore.(int64)
		}
	}

	userVotesByRecipe := make(map[int64]int64)
	for _, vote := range userVotes {
		userVotesByRecipe[vote.RecipeID] = vote.Vote
	}

	result := make(map[int64]domain.RecipeVotes)
	for _, recipeId := range recipeIds {
		result[recipeId] = s.mapper.ToRecipeVotes(votesByRecipe[recipeId], userVotesByRecipe[recipeId])
	}

	return result, nil
}

func (s *Store) groupTagsByRecipe(tags []database.GetTagsForRecipesRow) map[int64][]domain.Tag {
	tagsByRecipe := make(map[int64][]domain.Tag)
	for _, tag := range tags {
		tagsByRecipe[tag.RecipeID] = append(tagsByRecipe[tag.RecipeID], s.mapper.ToTag(tag.Tag))
	}
	return tagsByRecipe
}

func (s *Store) groupImagesByRecipe(images []database.GetImagesForRecipesRow) (map[int64][]domain.RecipeImage, error) {
	imagesByRecipe := make(map[int64][]domain.RecipeImage)
	for _, image := range images {
		imageUrl, err := url.ParseRequestURI(image.Path)
		if err != nil {
			return nil, err
		}
		imagesByRecipe[image.RecipeID] = append(imagesByRecipe[image.RecipeID], domain.RecipeImage{
			ID:  image.ID,
			URL: imageUrl,
		})
	}
	return imagesByRecipe, nil
}

func (s *Store) groupIngredientsByStep(ingredients []database.GetIngredientsForRecipesRow, populatedIngredientsMap map[int64]domain.Ingredient) map[int64][]domain.StepIngredient {
	ingredientsByStep := make(map[int64][]domain.StepIngredient)
	for _, ing := range ingredients {
		stepIngredient := s.mapper.ToStepIngredient(ing)
		stepIngredient.Ingredient = populatedIngredientsMap[ing.IngredientID]
		ingredientsByStep[ing.StepID] = append(ingredientsByStep[ing.StepID], stepIngredient)
	}
	return ingredientsByStep
}

func (s *Store) groupStepsByRecipe(steps []database.GetStepsForRecipesRow, ingredientsByStep map[int64][]domain.StepIngredient) map[int64][]domain.RecipeStep {
	stepsByRecipe := make(map[int64][]database.GetStepsForRecipesRow)
	for _, step := range steps {
		stepsByRecipe[step.RecipeID] = append(stepsByRecipe[step.RecipeID], step)
	}

	result := make(map[int64][]domain.RecipeStep)
	for recipeID, recipeSteps := range stepsByRecipe {
		stepsForRecipe := make([]domain.RecipeStep, len(recipeSteps))
		for i, step := range recipeSteps {
			recipeStep := s.mapper.ToRecipeStep(step)
			recipeStep.Ingredients = ingredientsByStep[step.ID]
			stepsForRecipe[i] = recipeStep
		}
		result[recipeID] = stepsForRecipe
	}
	return result
}

func (s *Store) UpdateRecipe(ctx context.Context, recipe domain.Recipe) (domain.Recipe, error) {
	err := s.WithTransaction(ctx, func(tx *TxStore) error {
		err := tx.query().UpdateRecipe(ctx, tx.mapper.FromRecipeForUpdate(recipe))
		if err != nil {
			return err
		}

		if err = tx.query().DeleteRecipeIngredients(ctx, recipe.ID); err != nil {
			return err
		}
		if err = tx.query().DeleteRecipeSteps(ctx, recipe.ID); err != nil {
			return err
		}
		if err = tx.query().DeleteRecipeImages(ctx, recipe.ID); err != nil {
			return err
		}
		if err = tx.query().DeleteRecipeTags(ctx, recipe.ID); err != nil {
			return err
		}

		if err = tx.createRecipeSteps(ctx, recipe.ID, recipe.Steps); err != nil {
			return err
		}
		if err = tx.createRecipeImages(ctx, recipe.ID, recipe.Images); err != nil {
			return err
		}
		if err = tx.createRecipeTags(ctx, recipe.ID, recipe.Tags); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return domain.Recipe{}, err
	}
	return s.GetRecipeById(ctx, recipe.CreatedBy, recipe.ID)
}

func (s *Store) AddVote(ctx context.Context, recipeID int64, userID int64, vote int64) error {
	return s.query().AddVote(ctx, database.AddVoteParams{
		RecipeID: recipeID,
		UserID:   userID,
		Vote:     vote,
	})
}

func (s *Store) RemoveVote(ctx context.Context, recipeID int64, userID int64) error {
	return s.query().RemoveVote(ctx, database.RemoveVoteParams{
		RecipeID: recipeID,
		UserID:   userID,
	})
}

func (s *Store) GetRecipeVotes(ctx context.Context, recipeID int64) (int64, error) {
	result, err := s.query().GetRecipeVotes(ctx, recipeID)
	if err != nil {
		return 0, err
	}
	if result == nil {
		return 0, nil
	}
	return result.(int64), nil
}

func (s *Store) GetUserVote(ctx context.Context, recipeID int64, userID int64) (int64, error) {
	return s.query().GetUserVote(ctx, database.GetUserVoteParams{
		RecipeID: recipeID,
		UserID:   userID,
	})
}

func (s *Store) CreateIngredient(ctx context.Context, ingredient domain.Ingredient) (domain.Ingredient, error) {
	var id int64
	err := s.WithTransaction(ctx, func(tx *TxStore) error {
		var err error
		id, err = tx.query().CreateIngredient(ctx, ingredient.Name)
		if err != nil {
			return err
		}

		for _, nutrient := range ingredient.Nutrients {
			err = tx.query().AddIngredientNutrient(ctx, tx.mapper.FromIngredientNutrient(id, nutrient))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return domain.Ingredient{}, err
	}
	ingredient.ID = id
	populated, err := s.populateIngredientNutrients(ctx, []domain.Ingredient{ingredient})
	if err != nil {
		return domain.Ingredient{}, err
	}
	return populated[0], nil
}

func (s *Store) UpdateIngredient(ctx context.Context, ingredient domain.Ingredient) (domain.Ingredient, error) {
	err := s.WithTransaction(ctx, func(tx *TxStore) error {
		err := tx.query().UpdateIngredient(ctx, database.UpdateIngredientParams{
			Name: ingredient.Name,
			ID:   ingredient.ID,
		})
		if err != nil {
			return err
		}

		// Delete existing nutrients and add new ones
		err = tx.query().DeleteIngredientNutrients(ctx, ingredient.ID)
		if err != nil {
			return err
		}

		for _, nutrient := range ingredient.Nutrients {
			err = tx.query().AddIngredientNutrient(ctx, tx.mapper.FromIngredientNutrient(ingredient.ID, nutrient))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return domain.Ingredient{}, err
	}
	populated, err := s.populateIngredientNutrients(ctx, []domain.Ingredient{ingredient})
	if err != nil {
		return domain.Ingredient{}, err
	}
	return populated[0], nil
}

func (s *Store) DeleteIngredient(ctx context.Context, id int64) error {
	return s.query().DeleteIngredient(ctx, id)
}

func (s *Store) CreateUnit(ctx context.Context, unit domain.Unit) (domain.Unit, error) {
	id, err := s.query().CreateUnit(ctx, database.CreateUnitParams{
		Name:   unit.Name,
		Symbol: unit.Symbol,
	})
	if err != nil {
		return domain.Unit{}, err
	}
	return domain.Unit{
		ID:     id,
		Name:   unit.Name,
		Symbol: unit.Symbol,
	}, nil
}

func (s *Store) UpdateUnit(ctx context.Context, unit domain.Unit) error {
	return s.query().UpdateUnit(ctx, database.UpdateUnitParams{
		Name:   unit.Name,
		Symbol: unit.Symbol,
		ID:     unit.ID,
	})
}

func (s *Store) DeleteUnit(ctx context.Context, id int64) error {
	return s.query().DeleteUnit(ctx, id)
}
