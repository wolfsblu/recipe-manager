import type { PageLoad } from './$types';
import { getRecipe, getIngredients, getUnits, getTags } from "$lib/api/recipes/recipes.svelte.js";
import { superValidate } from "sveltekit-superforms";
import { zod4 } from "sveltekit-superforms/adapters";
import { formSchema } from "../../add/schema";

export const prerender = false

export const load: PageLoad = async ({ params }) => {
    const recipeId = Number(params.id)

    const [recipe, ingredientsResponse, unitsResponse, tagsResponse] = await Promise.all([
        getRecipe(recipeId),
        getIngredients({ limit: 100 }),
        getUnits({ limit: 100 }),
        getTags({ limit: 100 })
    ])

    const ingredients = ingredientsResponse.data
    const units = unitsResponse.data
    const tags = tagsResponse.data

    const breadcrumbs = [
        { link: "/", name: "Home" },
        { link: "/recipes", name: "Recipes" },
        { link: `/recipes/${recipeId}`, name: recipe.name },
        { link: `/recipes/${recipeId}/edit`, name: "Edit" },
    ]

    const form = await superValidate({
        name: recipe.name,
        description: recipe.description,
        servings: recipe.servings,
        minutes: recipe.minutes,
        images: recipe.images?.map(image => image || '') || [],
        tags: recipe.tags?.map(tag => tag.id) || [],
        steps: recipe.steps?.map(step => ({
            instructions: step.instructions,
            ingredients: step.ingredients?.map(ingredient => ({
                amount: ingredient.amount,
                ingredientId: ingredient.id,
                unitId: ingredient.unit.id,
            })) || []
        })) || []
    }, zod4(formSchema));

    return {
        breadcrumbs,
        recipe,
        ingredients,
        units,
        tags,
        form,
    };
};