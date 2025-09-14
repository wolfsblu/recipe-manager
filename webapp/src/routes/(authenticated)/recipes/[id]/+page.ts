import type { PageLoad } from './$types';
import {getRecipe} from "$lib/api/recipes/recipes.svelte";

export const prerender = false

export const load: PageLoad = async ({ params }) => {
    const recipeId = Number(params.id)
    const recipe = await getRecipe(recipeId)

    const breadcrumbs = [
        { link: "/", name: "Home" },
        { link: "/recipes", name: "Recipes" },
        { link: "", name: recipe.name },
    ]

    return {
        breadcrumbs,
        recipe,
    };
};