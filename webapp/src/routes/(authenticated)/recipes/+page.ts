import type { PageLoad } from './$types';
import { getRecipes } from "$lib/api/recipes/recipes.svelte";

export const load: PageLoad = async () => {
    const breadcrumbs = [
        { link: "/", name: "Home" },
        { link: "/recipes", name: "Recipes" },
    ]

    const recipes = await getRecipes()
    return {
        breadcrumbs,
        recipes,
    };
};