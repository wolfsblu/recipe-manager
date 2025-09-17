package sqlite

import (
	"context"
	"net/url"
	"time"

	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/database"
	_ "modernc.org/sqlite"
)

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
	if err := s.Begin(ctx); err != nil {
		return domain.Recipe{}, err
	}
	defer s.Rollback()

	recipeId, err := s.query().CreateRecipe(ctx, s.mapper.FromRecipe(recipe))
	if err != nil {
		return domain.Recipe{}, err
	}
	if err = s.createRecipeSteps(ctx, recipeId, recipe.Steps); err != nil {
		return domain.Recipe{}, err
	}
	if err = s.createRecipeImages(ctx, recipeId, recipe.Images); err != nil {
		return domain.Recipe{}, err
	}
	if err = s.createRecipeTags(ctx, recipeId, recipe.Tags); err != nil {
		return domain.Recipe{}, err
	}
	if err = s.Commit(); err != nil {
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

func (s *Store) GetMealPlan(ctx context.Context, user *domain.User, from time.Time, until time.Time) ([]domain.MealPlan, error) {
	result, err := s.query().GetMealPlan(ctx, database.GetMealPlanParams{
		UserID:    user.ID,
		FromDate:  from.Format(time.DateOnly),
		UntilDate: until.Format(time.DateOnly),
	})
	if err != nil {
		return []domain.MealPlan{}, err
	}

	grouped := make(map[string][]domain.Recipe)
	for _, item := range result {
		grouped[item.MealPlan.Date] = append(grouped[item.MealPlan.Date], s.mapper.ToRecipe(item.Recipe))
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

	return ingredients, nil
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
	recipeIds := []int64{id}
	steps, err := s.query().GetStepsForRecipes(ctx, recipeIds)
	if err != nil {
		return recipe, err
	}
	ingredients, err := s.query().GetIngredientsForRecipes(ctx, recipeIds)
	if err != nil {
		return recipe, err
	}
	databaseTags, err := s.query().GetTagsForRecipes(ctx, recipeIds)
	if err != nil {
		return recipe, err
	}
	images, err := s.query().GetImagesForRecipes(ctx, recipeIds)
	if err != nil {
		return recipe, err
	}

	ingredientsByStep := make(map[int64][]domain.StepIngredient)
	for _, ing := range ingredients {
		ingredientsByStep[ing.StepID] = append(ingredientsByStep[ing.StepID], s.mapper.ToStepIngredient(ing))
	}

	recipeSteps := make([]domain.RecipeStep, len(steps))
	for i, step := range steps {
		recipeStep := s.mapper.ToRecipeStep(step)
		recipeStep.Ingredients = ingredientsByStep[step.ID]
		recipeSteps[i] = recipeStep
	}

	recipeTags := make([]domain.Tag, len(databaseTags))
	for i, databaseTag := range databaseTags {
		recipeTags[i] = s.mapper.ToTag(databaseTag.Tag)
	}

	recipeImages := make([]domain.RecipeImage, len(images))
	for i, image := range images {
		recipeImages[i] = s.mapper.ToRecipeImage(image)
		imageUrl, err := url.ParseRequestURI(image.Path)
		if err != nil {
			return recipe, err
		}
		recipeImages[i].URL = imageUrl
	}

	votes, err := s.getVotes(ctx, user, []int64{id})
	if err != nil {
		return recipe, err
	}

	recipe.Steps = recipeSteps
	recipe.Tags = recipeTags
	recipe.Images = recipeImages
	recipe.Votes = votes[id]
	return recipe, nil
}

func (s *Store) GetRecipesByUser(ctx context.Context, user *domain.User) (recipes []domain.Recipe, _ error) {
	result, err := s.query().ListRecipes(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, len(result))
	for i, item := range result {
		ids[i] = item.ID
		recipes = append(recipes, s.mapper.ToRecipe(item))
	}

	tags, err := s.query().GetTagsForRecipes(ctx, ids)
	if err != nil {
		return nil, err
	}

	images, err := s.query().GetImagesForRecipes(ctx, ids)
	if err != nil {
		return nil, err
	}

	tagsByRecipe := make(map[int64][]domain.Tag)
	for _, tag := range tags {
		tagsByRecipe[tag.RecipeID] = append(tagsByRecipe[tag.RecipeID], s.mapper.ToTag(tag.Tag))
	}

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

	votes, err := s.getVotes(ctx, user, ids)
	if err != nil {
		return nil, err
	}

	recipes = make([]domain.Recipe, len(result))
	for i, recipe := range result {
		r := s.mapper.ToRecipe(recipe)
		r.Images = imagesByRecipe[r.ID]
		r.Tags = tagsByRecipe[r.ID]
		r.Votes = votes[r.ID]
		recipes[i] = r
	}
	return recipes, nil
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

func (s *Store) UpdateRecipe(ctx context.Context, recipe domain.Recipe) (domain.Recipe, error) {
	if err := s.Begin(ctx); err != nil {
		return domain.Recipe{}, err
	}
	defer s.Rollback()

	err := s.query().UpdateRecipe(ctx, s.mapper.FromRecipeForUpdate(recipe))
	if err != nil {
		return domain.Recipe{}, err
	}

	if err = s.query().DeleteRecipeIngredients(ctx, recipe.ID); err != nil {
		return domain.Recipe{}, err
	}
	if err = s.query().DeleteRecipeSteps(ctx, recipe.ID); err != nil {
		return domain.Recipe{}, err
	}
	if err = s.query().DeleteRecipeImages(ctx, recipe.ID); err != nil {
		return domain.Recipe{}, err
	}
	if err = s.query().DeleteRecipeTags(ctx, recipe.ID); err != nil {
		return domain.Recipe{}, err
	}

	if err = s.createRecipeSteps(ctx, recipe.ID, recipe.Steps); err != nil {
		return domain.Recipe{}, err
	}
	if err = s.createRecipeImages(ctx, recipe.ID, recipe.Images); err != nil {
		return domain.Recipe{}, err
	}
	if err = s.createRecipeTags(ctx, recipe.ID, recipe.Tags); err != nil {
		return domain.Recipe{}, err
	}

	if err = s.Commit(); err != nil {
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
